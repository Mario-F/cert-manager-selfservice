/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

/**
 * This is a single certificate for consuming on clients
 */
export type Certificate = {
    /**
     * The domain for which the certificate is issued
     */
    domain: string;
    /**
     * The certificate authority for the certificate
     */
    authority: string;
    /**
     * The private key for the certificate
     */
    key: string;
    /**
     * The certificate for the certificate
     */
    crt: string;
};
