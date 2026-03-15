#!/usr/bin/env python3

import requests
import os
import re
import socket
import struct
from dotenv import load_dotenv
import sys
import json
from time import sleep

CHALLENGE_DOMAINS_FILE = "/tmp/certbot_challenge_domains.json"
DNS_POLL_INTERVAL = 15
DNS_POLL_TIMEOUT = 600
DNS_INITIAL_WAIT = 60
NAMESERVER = "ns.wedos.cz"


def get_tld(validation_domain):
    parts = validation_domain.split(".")
    if len(parts) <= 2:
        return validation_domain
    tld_pattern = r"\.([^.]+\.[^.]+)$"
    match = re.search(tld_pattern, validation_domain)
    if match:
        return match.group(1)
    return None


def get_validation_subdomain(validation_domain, domain):
    if validation_domain == domain:
        return ""
    return validation_domain.replace("." + domain, "")


def authorize():
    payload = {
        "login": os.environ.get("BW_USER_LOGIN"),
        "secret": os.environ.get("BW_USER_SECRET"),
    }

    response = requests.post(
        os.environ.get("BW_BASE_URL") + "/api/auth/token", json=payload
    )

    if response.status_code == 200:
        return response.json()["token"]

    print("Error authorizing. Response code: " + str(response.status_code))
    raise Exception("Authorization error")


def save_challenge_domain(fqdn, token):
    domains = {}
    if os.path.exists(CHALLENGE_DOMAINS_FILE):
        with open(CHALLENGE_DOMAINS_FILE, "r") as f:
            domains = json.load(f)
    domains[fqdn] = token
    with open(CHALLENGE_DOMAINS_FILE, "w") as f:
        json.dump(domains, f)


def load_challenge_domains():
    if os.path.exists(CHALLENGE_DOMAINS_FILE):
        with open(CHALLENGE_DOMAINS_FILE, "r") as f:
            return json.load(f)
    return {}


def clear_challenge_domains():
    if os.path.exists(CHALLENGE_DOMAINS_FILE):
        os.remove(CHALLENGE_DOMAINS_FILE)


def resolve_nameserver(hostname):
    try:
        return socket.getaddrinfo(hostname, 53, socket.AF_INET)[0][4][0]
    except socket.gaierror:
        return "185.8.238.1"


def build_dns_query(fqdn):
    transaction_id = os.urandom(2)
    flags = struct.pack("!H", 0x0100)
    counts = struct.pack("!HHHH", 1, 0, 0, 0)
    question = b""
    for label in fqdn.rstrip(".").split("."):
        question += struct.pack("B", len(label)) + label.encode()
    question += b"\x00"
    question += struct.pack("!HH", 16, 1)
    return transaction_id + flags + counts + question


def parse_dns_response(data):
    txt_records = []
    if len(data) < 12:
        return txt_records
    rcode = data[3] & 0x0F
    if rcode != 0:
        return txt_records
    qdcount = struct.unpack("!H", data[4:6])[0]
    ancount = struct.unpack("!H", data[6:8])[0]
    offset = 12
    for _ in range(qdcount):
        while offset < len(data):
            length = data[offset]
            if length == 0:
                offset += 1
                break
            if length >= 192:
                offset += 2
                break
            offset += 1 + length
        offset += 4
    for _ in range(ancount):
        if offset >= len(data):
            break
        if data[offset] >= 192:
            offset += 2
        else:
            while offset < len(data):
                length = data[offset]
                if length == 0:
                    offset += 1
                    break
                offset += 1 + length
        if offset + 10 > len(data):
            break
        rtype, rclass, rttl, rdlength = struct.unpack(
            "!HHIH", data[offset : offset + 10]
        )
        offset += 10
        if rtype == 16 and offset + rdlength <= len(data):
            rdata = data[offset : offset + rdlength]
            pos = 0
            txt_value = ""
            while pos < len(rdata):
                txt_len = rdata[pos]
                pos += 1
                txt_value += rdata[pos : pos + txt_len].decode(
                    "utf-8", errors="replace"
                )
                pos += txt_len
            txt_records.append(txt_value)
        offset += rdlength
    return txt_records


def check_dns_txt(fqdn, expected_token):
    try:
        ns_ip = resolve_nameserver(NAMESERVER)
        query = build_dns_query(fqdn)
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.settimeout(5)
        sock.sendto(query, (ns_ip, 53))
        data, _ = sock.recvfrom(4096)
        sock.close()
        for record in parse_dns_response(data):
            if expected_token in record:
                return True
    except Exception:
        pass
    return False


def wait_for_dns_propagation(challenge_domains):
    print(f"Waiting {DNS_INITIAL_WAIT}s for initial DNS propagation before polling...")
    sleep(DNS_INITIAL_WAIT)
    print(
        f"Polling DNS for {len(challenge_domains)} challenge record(s) "
        f"(interval={DNS_POLL_INTERVAL}s, timeout={DNS_POLL_TIMEOUT}s)..."
    )
    elapsed = 0
    while elapsed < DNS_POLL_TIMEOUT:
        all_resolved = True
        for fqdn, token in challenge_domains.items():
            if not check_dns_txt(fqdn, token):
                all_resolved = False
                break
        if all_resolved:
            print(f"All DNS TXT records resolved after {elapsed}s.")
            return True
        sleep(DNS_POLL_INTERVAL)
        elapsed += DNS_POLL_INTERVAL
        print(f"  Waiting... ({elapsed}s/{DNS_POLL_TIMEOUT}s)")
    print(f"DNS propagation timeout after {DNS_POLL_TIMEOUT}s.")
    return False


def perform_dns_challenge(validation_domain, validation_token):
    domain = get_tld(validation_domain)

    headers = {
        "Authorization": "Bearer " + authorize(),
        "Content-Type": "application/json",
    }

    validation_sub = get_validation_subdomain(validation_domain, domain)
    subdomain = "_acme-challenge" + ("." + validation_sub if validation_sub else "")

    payload = json.dumps(
        {
            "autocommit": True,
            "data": validation_token,
            "subdomain": subdomain,
            "ttl": 300,
            "type": "TXT",
        }
    )

    response = requests.post(
        os.environ.get("BW_BASE_URL") + "/api/v1/domain/" + domain + "/record",
        data=payload,
        headers=headers,
    )

    if response.status_code == 201:
        print("DNS TXT record created for " + validation_domain)
        challenge_fqdn = subdomain + "." + domain
        save_challenge_domain(challenge_fqdn, validation_token)
        remaining = int(os.environ.get("CERTBOT_REMAINING_CHALLENGES", "0"))
        if remaining == 0:
            challenge_domains = load_challenge_domains()
            if not wait_for_dns_propagation(challenge_domains):
                raise Exception("DNS propagation timed out.")
            clear_challenge_domains()
        print("DNS challenge completed successfully.")
    else:
        print("Error performing DNS challenge.")
        raise Exception("DNS challenge failed.")


def cleanup_dns_challenge(validation_domain, validation_token):
    domain = get_tld(validation_domain)

    headers = {
        "Authorization": "Bearer " + authorize(),
        "Content-Type": "application/json",
    }

    validation_sub = get_validation_subdomain(validation_domain, domain)
    subdomain = "_acme-challenge" + ("." + validation_sub if validation_sub else "")

    payload = json.dumps(
        {
            "autocommit": True,
            "data": validation_token,
            "subdomain": subdomain,
            "ttl": 300,
            "type": "TXT",
        }
    )

    response = requests.delete(
        os.environ.get("BW_BASE_URL") + "/api/v1/domain/" + domain + "/record/",
        data=payload,
        headers=headers,
    )

    if response.status_code == 200:
        print("DNS challenge cleanup completed successfully.")
    else:
        print("Error performing DNS challenge cleanup.")
        print("Response: " + response.text)
        raise Exception("DNS challenge cleanup failed.")


if __name__ == "__main__":
    load_dotenv()

    hook_action = sys.argv[1]
    domain = os.environ.get("CERTBOT_DOMAIN")
    token = os.environ.get("CERTBOT_VALIDATION")

    if hook_action == "deploy_challenge":
        perform_dns_challenge(domain, token)
    elif hook_action == "clean_challenge":
        cleanup_dns_challenge(domain, token)
