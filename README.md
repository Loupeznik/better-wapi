# Better WAPI

This project is inteded to serve as a more standardized wrapper around the Wedos API (WAPI).

![GitHub](https://img.shields.io/github/license/loupeznik/better-wapi?style=for-the-badge)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/loupeznik/better-wapi?style=for-the-badge)

It currently offers following functionality:

- Add DNS record for a specific domain
- Update DNS record for specific domain
- Remove DNS record for specific domain
- List all DNS records for specific domain
- List a particular DNS record for specific domain

The *BetterWAPI* project uses RESTful API style and tries to take a more standardized approach,
something that the original WAPI is missing entirely. The core functionality remains similar.

For creating and updating records, only operations on **A** record types are currently supported.

This project is still in active development and far from final release. All features are subject to change
(while trying to preserve backward compatibility of course).

## Installation

For the API to work, it is first important to whitelist the IP address of the host machine in the
*[WEDOS management dashboard](https://client.wedos.com/client/wapi.html)*. It is best to have a server with a static
IP address assigned and have this address whitelisted (for production environments).

To run the project locally:

```bash
git clone https://github.com/Loupeznik/better-wapi.git
cd better-wapi
go get .
cp .env.example .env
```

Fill the .env file with your credentials.

- The BW_WAPI_ variables are your WAPI credentials from the WEDOS management dashboard
- The BW_USER_ variables are credentials to use within your API

Alternatively, it is possible to use environment variables without using the .env file.

Example in Powershell:

```powershell
$Env:BW_USER_LOGIN = "admin"
$Env:BW_USER_SECRET = "admin"
$Env:BW_WAPI_LOGIN = "admin@example.com"
$Env:BW_WAPI_PASSWORD = "yourpassword"
```

Finally, to run the API.

```bash
go run .
```

For production workloads, a web server like NGINX is needed, the .env file also needs to be present.

## Running in Docker

An option to run the API in Docker is available as well.

Building the image:

```bash
docker build -t better-wapi:latest .
docker run -d -p 8083:8000 --env-file .\.env better-wapi:latest
```

Alternatively, get the image from Dockerhub

```bash
docker pull docker pull loupeznik/better-wapi
docker run -d -p 8083:8000 --env-file .\.env loupeznik/better-wapi:latest
```

## Documentation

- WAPI documentation - <https://kb.wedos.com/en/kategorie/wapi-api-interface/>
- Better WAPI documentation - Pending

## Example usage

### List all subdomains

```bash
curl --location --request GET 'http://127.0.0.1:8000/api/domain/dzarsky.eu/info' \
--header 'Authorization: Basic aGVsb3U6eWVz'
```

### Update a record

```bash
curl --location --request PUT 'http://127.0.0.1:8000/api/domain/dzarsky.eu/record' \
--header 'Authorization: Basic aGVsb3U6eWVz' \
--header 'Content-Type: application/json' \
--data-raw '{
    "subdomain": "*",
    "ip": "123.123.123.123"
}'
```

The API uses basic auth with the BW_USER credentials set in the .env file.

## License

This project is [GPL-3.0 licensed](https://github.com/Loupeznik/better-wapi/blob/master/LICENSE).

Created by [Dominik Zarsky](https://github.com/Loupeznik).
