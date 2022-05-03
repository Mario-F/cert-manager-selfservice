/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
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

}