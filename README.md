# cert-manager-selfservice

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
![GitHub issues](https://img.shields.io/github/issues/Mario-F/cert-manager-selfservice)

This project aims to utilize a working cert-manager installation to provide certificates outside kubernetes as easy as possible.

## Overview

What does cert-manager-selfservice (CMS) offer?

* Just make an http call to get your certificate, example: <http://localhost:8030/cert/your.domain.tld/pem>
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
wget -O /etc/ssl/selfservice/service.test.example.com.pem http://selfservice.example.com/cert/service.test.example.com/pem
```

This will request a certificte for domain `service.test.example.com` from selfservice, at the very first request for this domain the file under `/etc/ssl/selfservice/service.test.example.com.pem` will created empty.

This is because cert-manager creating certificates asynchronously the commonly used lets-encrypt certificates will normally take more than one minute to populate.

Selfservice will return HTTP Code 204 until the certificate is ready to use and normal Code 200 when its ready, this means you should check your request for HTTP Code 200.

## Development

### Testing

The most simplest usage (for testing) would to run cert-manager-selfservice with your local kubeconfig, this can be done by:

```shell
./cert-manager-selfservice server --issuer-name your-issuer-to-use
```

Then you can request a certificate by calling: `http://localhost:8030/cert/your.domain.tld/pem`

If the certificate not exists a certificate ressource will automatically be created, until there is no valid secret (issued certificate) a HTTP 202 will be returned.

There a also other endpoints like `crt`, `key`, `ca`, `json` availiable.
