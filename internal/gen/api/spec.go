// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RWwW7kNgz9FUHt0RvPZvfk2yIFihwWKJreFjkwNu3RxpZcinZqBPPvBSV77Bl7Mslh",
	"bxlTenx6fCTzqnPXtM6iZa+zV+3zPTYQ/rxDYlOaHBjlZ4E+J9OycVZn+p+98cp4BcobW9Wo8vm0Kh2p",
	"3FnfNcZWylmV1yYkSHRLrpWTGFJAx3tHhoetBKeYx6MBnU+jOtE8tKgz7ZmMrfQh0TnxddR3YhWuAWO3",
	"4WIsIL3sTb4/x1NBKN9hsYX8jBfe3pLp5fYzvvPFh0QT/tsZwkJnPybKyULjmC0q83i8755+Ys7C5d6W",
	"Tsic1qhH8oHWFssxqFx5JPipAQsV0iePdemRepNfZzsl2aL1wMCdXxNbqBF+/05Y6kz/ls6OTkc7pxHj",
	"bnnjkOgGvYcq3t6wNxAq7NGymg6qllzR5Viop+Hagw1j805e3yO8UBpfD0QwrFQ68k1OH39ZtbszjS41",
	"seuFN75IIaGul07zyscCnDdvCaYWXq9TcmMZKyTJ3qItpM6bQUIohu0QO4Z6O9TZZ+te7FbwTKYIMpOY",
	"MiYT5RnssnJTTVa2a+bAqps99jiNM7RdI2yMdFWiX4BsJINEjhaZ5+tsmoBbOmqAdaYLYPwUvl5roPHQ",
	"RG5BZf1EuWvGXmfDtcTEKOp7NLJ6mI2svv11r5N5COjPN7ubnZB1LVpojc70l5vdzRfRG3gfJEoX9klf",
	"4yQ6SKDCMJNFThAP3hc6038iL2z6xzS3WiBokJG8zn68MXjZqQp5NXVLR9KEclh46URbEHXnwTjLx9Rh",
	"Mq6+jboeHuWwb5310QK3u10YQM4y2vAiaNta8hpn058+TssZ760JsNyyoTKnD/2mBE7F0inCltDLRDpf",
	"B4dE3+5uz1jtualTxv/kxxr6AalHUkwGvagIQY6VkIPrlN+7ri4U06CgEtFr4NiVX3dfP5R0WsAqB2sd",
	"qyepVGcLZVhB7Z18j59zQmAsgtN91zRAg87038gdWQW2iKO5g7oexrObm70yPdrRKwErnax/yYz3sWF/",
	"WcUD/gdLbeKdtRI+hGRcyBqGJ9fx24spSOCPK/WSCA/TzP9lMowZPiiEn25tSRGDV/8ZCSl9cH+cLR3V",
	"OtMptCbtP+vD4+H/AAAA///HggM7FgsAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

