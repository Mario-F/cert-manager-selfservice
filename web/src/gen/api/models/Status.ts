/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { StatusCertificates } from './StatusCertificates';
import type { StatusMessage } from './StatusMessage';

export type Status = {
    /**
     * This are event messages produced by the cert-manager-selfservice
     */
    messages: Array<StatusMessage>;
    certificates: StatusCertificates;
};
