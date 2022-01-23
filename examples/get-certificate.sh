#!/bin/sh
#
# This is an example script using https://github.com/Mario-F/cert-manager-selfservice
#
SCRIPTNAME=$0
API=$1
DOMAIN=$2
TYPE=${3:-pem}

usage () {
    echo "usage: $SCRIPTNAME apiurl domain [type]"
    echo "type can be pem,crt,key,ca,json (default is pem)"
    exit 1
}

if [ -z "${API}" ]; then usage; fi
if [ -z "${DOMAIN}" ]; then usage; fi

REQUEST_URL="${API}/cert/${DOMAIN}/${TYPE}"
TMP_FILE="/tmp/${DOMAIN}.${TYPE}"
TARGET_FOLDER="/etc/ssl/selfservice"
TARGET_FILE="${TARGET_FOLDER}/${DOMAIN}.${TYPE}"

# Download cert to /tmp/ and store http response code
HTTP_RES_CODE=$(wget -O ${TMP_FILE} --server-response --quiet "${REQUEST_URL}" 2>&1 | awk 'NR==1{print $2}')

if [ "$HTTP_RES_CODE" = "200" ]; then
  mkdir -p ${TARGET_FOLDER}
  cp ${TMP_FILE} ${TARGET_FILE}
fi
