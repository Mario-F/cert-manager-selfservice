# cert-manager-selfservice

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
![GitHub issues](https://img.shields.io/github/issues/Mario-F/cert-manager-selfservice)

This project aims to utilize a working cert-manager installation to provide certificates outside kubernetes as easy as possible.

## Usage

### Testing

The most simplest usage (for testing) would to run cert-manager-selfservice with your local kubeconfig, this can be done by:

```shell
./cert-manager-selfservice server --issuer-name your-issuer-to-use
```

Then you can request a certificate by calling: `http://localhost:8030/cert/your.domain.tld/pem`

If the certificate not exists a certificate ressource will automatically be created, until there is no valid secret (issued certificate) a HTTP 202 will be returned.

There a also other endpoints like `crt`, `key`, `ca` availiable.

### Deploy

* [Helm Chart](https://github.com/Mario-F/helm-charts/tree/main/charts/cert-manager-selfservice)
