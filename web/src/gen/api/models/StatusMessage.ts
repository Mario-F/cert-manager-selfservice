/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type StatusMessage = {
    time: string;
    severity: StatusMessage.severity;
    message: string;
};

export namespace StatusMessage {

    export enum severity {
        INFO = 'info',
        WARNING = 'warning',
        ERROR = 'error',
    }


}
