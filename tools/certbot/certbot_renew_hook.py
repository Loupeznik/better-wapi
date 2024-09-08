#!/usr/bin/env python3

import requests
import os
import re
from dotenv import load_dotenv
import sys
import json
from time import sleep


def get_tld(validation_domain):
    tld_pattern = r'\.([^.]+\.[^.]+)$'
    match = re.search(tld_pattern, validation_domain)
    if match:
        tld = match.group(1)
        return tld
    else:
        return None


def get_validation_subdomain(validation_domain, domain):
    return validation_domain.replace('.' + domain, '')


def authorize():
    payload = {
        'login': os.environ.get('BW_USER_LOGIN'),
        'secret': os.environ.get('BW_USER_SECRET')
    }

    response = requests.post(os.environ.get(
        'BW_BASE_URL') + '/api/auth/token', json=payload)

    if response.status_code == 200:
        return response.json()['token']

    print('Error authorizing. Response code: ' + str(response.status_code))
    raise Exception('Authorization error')


def perform_dns_challenge(validation_domain, validation_token):
    domain = get_tld(validation_domain)

    headers = {
        'Authorization': 'Bearer ' + authorize(),
        'Content-Type': 'application/json'
    }

    payload = json.dumps({
        'autocommit': True,
        'data': validation_token,
        'subdomain': '_acme-challenge.' + get_validation_subdomain(validation_domain, domain),
        'ttl': 300,
        'type': 'TXT'
    })

    response = requests.post(os.environ.get(
        'BW_BASE_URL') + '/api/v1/domain/' + domain + '/record', data=payload, headers=headers)

    if response.status_code == 201:
        sleep(600)
        print('DNS challenge completed successfully.')
    else:
        print('Error performing DNS challenge.')
        raise Exception('DNS challenge failed.')


def cleanup_dns_challenge(validation_domain, validation_token):
    domain = get_tld(validation_domain)

    headers = {
        'Authorization': 'Bearer ' + authorize(),
        'Content-Type': 'application/json'
    }

    payload = json.dumps({
        'autocommit': True,
        'data': validation_token,
        'subdomain': '_acme-challenge.' + get_validation_subdomain(validation_domain, domain),
        'ttl': 300,
        'type': 'TXT'
    })

    response = requests.delete(os.environ.get('BW_BASE_URL') + '/api/v1/domain/' + domain
                               + '/record/', data=payload, headers=headers)

    if response.status_code == 200:
        print('DNS challenge cleanup completed successfully.')
    else:
        print('Error performing DNS challenge cleanup.')
        print('Response: ' + response.text)
        raise Exception('DNS challenge cleanup failed.')


if __name__ == '__main__':
    load_dotenv()

    hook_action = sys.argv[1]
    domain = os.environ.get('CERTBOT_DOMAIN')
    token = os.environ.get('CERTBOT_VALIDATION')

    if hook_action == 'deploy_challenge':
        perform_dns_challenge(domain, token)
    elif hook_action == 'clean_challenge':
        cleanup_dns_challenge(domain, token)
