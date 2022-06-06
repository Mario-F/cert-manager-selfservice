/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Certificate } from '../models/Certificate';
import type { Info } from '../models/Info';
import type { Status } from '../models/Status';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class DefaultService {

    /**
     * Returns information about the cert-manager-selfservice
     * @returns Info A json object represent the info
     * @throws ApiError
     */
    public static getInfo(): CancelablePromise<Info> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/info',
        });
    }

    /**
     * Returns status of the cert-manager-selfservice
     * @returns Status A json object represent the status
     * @throws ApiError
     */
    public static getStatus(): CancelablePromise<Status> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/status',
        });
    }

    /**
     * Return and eventually create certificate for the given domain
     * @param domain The domain to get the certificate for
     * @param format Format of the result, defaults to json
     * @returns Certificate A json object represent the certificate
     * @returns any Server tries to aquire the certificate you should try again later
     * @throws ApiError
     */
    public static getCertificate(
        domain: string,
        format: 'json' | 'crt' | 'key' | 'ca' | 'pem' = 'json',
    ): CancelablePromise<Certificate | any> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/certificate/{domain}',
            path: {
                'domain': domain,
            },
            query: {
                'format': format,
            },
            errors: {
                404: `The cert cannot be found it also canot be created`,
            },
        });
    }

}