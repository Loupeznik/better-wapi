# Certbot certificate create/renew hook

A hook for certbot to create/renew TLS certificates using DNS challenge via Better WAPI.

## Usage

```bash
pip3 install -r requirements.txt
```

```bash
sudo certbot certonly --manual --preferred-challenges dns \
--manual-auth-hook "./certbot_renew_hook.py deploy_challenge" \ 
--manual-cleanup-hook "./certbot_renew_hook.py clean_challenge" \
-d '*.example.com'
```

Note that the Better WAPI `.env` file needs to be present in the same directory as the hook script.
