# cert-manager-selfservice

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
![GitHub issues](https://img.shields.io/github/issues/Mario-F/cert-manager-selfservice)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate?hosted_button_id=34NHCDNHRRV6G)

This project aims to utilize a working cert-manager installation to provide certificates outside kubernetes as easy as possible.

## Overview

What does cert-manager-selfservice (CMS) offer?

* Just make an http call to get your certificate, example: <http://localhost:8030/api/v1/certificate/your.domain.tld/pem>
* CMS creates certificates ressources automatically
* CMS keep track when certificates are accessed
* CMS cleanup certificates not requested in a time

## Install

You need a working [cert-manager](https://cert-manager.io/) installation and issuer for the domains you want to get certificates.

Simply use the [Helm Chart](https://github.com/Mario-F/helm-charts/tree/main/charts/cert-manager-selfservice) to get started.

## Usage

Expose your selfservice by ingress for example, the following examples assume that the selfservice http is reachable under <http://selfservice.example.com>

Login to the target system that want to use a certificate, create an directory and simply execute the commands:

```shell
mkdir /etc/ssl/selfservice
wget -O /etc/ssl/selfservice/service.test.example.com.pem http://selfservice.example.com/api/v1/certificate/service.test.example.com/pem
```

This will request a certificte for domain `service.test.example.com` from selfservice, at the very first request for this domain the file under `/etc/ssl/selfservice/service.test.example.com.pem` will created empty.

This is because cert-manager creating certificates asynchronously the commonly used lets-encrypt certificates will normally take more than one minute to populate.

Selfservice will return HTTP Code 202 until the certificate is ready to use and normal Code 200 when its ready, this means you should check your request for HTTP Code 200.

This [example script](./examples/get-certificate.sh) can be used to get certificates only when ready, this simple call will put the final certificate under `/etc/ssl/selfservice/service.test.example.com.pem` when ready:

```shell
get-certificate.sh http://selfservice.example.com service.test.example.com
```

If you run it in a cronjob the certificate will automatically renewed regullary.

## Development

At the moment you need a working kubernetes cluster with cert-manager to get started in development.

### OpenAPI

The API Server and Client are generated using [OpenAPI](./openapi.yaml), generation can be done with `make generate`.

The following tools are needed to generate the OpenAPI files:

* GO: [oapi-codegen](https://github.com/deepmap/oapi-codegen)
* Typescript: [openapi](https://github.com/openapi/openapi)

### Web

The [web project](./web/README.md) will bundled on build time into the executable, so you need to make sure the web project is build before.

### Testing

The most simplest usage (for testing) would to run cert-manager-selfservice with your local kubeconfig, this can be done by:

```shell
./cert-manager-selfservice server --issuer-name your-issuer-to-use
```

Then you can request a certificate by calling: `http://localhost:8030/api/v1/certificate/your.domain.tld/pem`

If the certificate not exists a certificate ressource will automatically be created, until there is no valid secret (issued certificate) a HTTP 202 will be returned.

There a also other endpoints like `crt`, `key`, `ca`, `json` availiable.

### Debugger

There a serveral ways to easy start development and using live debugging provided by [delve](https://github.com/go-delve/delve)

#### VSCode integrated Console

The provided launch.json has a debug task `Launch File` predefined, just hit start and it should run with the args provided in launch.json.

#### VSCode external Terminal

A more advanced way to test in an external Terminal is provided by the `External Debugging` launch config and `./debug` script:

1. Execute debug script with arguments as normal: `./debug server --issuer-name your-issuer-to-use`
2. Start the `External Debugging` session in vscode

Unfortunately the order is importend because vscode does not try to automatically connect after start.
