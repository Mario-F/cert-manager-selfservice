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

	"H4sIAAAAAAAC/7RWwW7jNhD9FYLtUYm92T3ptkjRIocFiqa3RQ4TaSRzI5Ha4dCpEOjfiyElS7HlODns",
	"zfKQbx7fvBnyRReu7ZxFy17nL9oXO2wh/rxFYlOZAhjls0RfkOnYOKtz/e/OeGW8AuWNrRtUxbxaVY5U",
	"4awPrbG1clYVjYkJMt2R62QlxhQQeOfIcL+W4DXmYWlE59dRnWnuO9S59kzG1nrIdEF8GfWdWKVrwdh1",
	"uBSLSM87U+yO8VQUygcs15Cf8MzZOzJ72f2E7zzxkGnCn8EQljr/PlHOFhqnbEmZh8N+9/gDCxYud7Zy",
	"QuZ1jfZIPtJaYzkGlasOBK9asFAjXXlsKo+0N8VltlOSNVr3DBz8KbGFGvH7d8JK5/q3zezozWjnTcK4",
	"Xe4YMt2i91Cn3Sv2BkKFe7SspoWqI1eGAkv12F86sGFs38nrW4IXSuPpgQj6E5UOfLPXhz+v2u2RRuea",
	"2O2FNz5LIaFplk7zyqcCHDdvBaYRXi9TcmMZayTJ3qEtpc6rQUIo+/UQO4ZmPRTsk3XPdi14JFMCmUlM",
	"GbOJ8gx2XrmpJie2a+fASTd73OM0ztCGVtgY6apMPwPZRAaJHC0yz9vZtBG3ctQC61yXwHgV/73UQOOi",
	"idyCyukRZa8Ze50NNxITo6hvycjqfjay+vr3nc7mIaA/XW+vt0LWdWihMzrXn6+3159Fb+BdlGizsM/m",
	"JU2iQQI1xpkscoJ48K7Uuf4LeWHTP6a51QFBi4zkdf79jcHLTtXIJ1O3ciRNKIuFl860BVF3HoyzfEwB",
	"s/HqW6nrkB2n/zMWaJp6hD40nKkSKwgNe6H0wzs75f8ZkPqZwFjdZcJxp871uG8yz/gpE/swvSFau11x",
	"0PAgp/Kdsz559Wa7jZPSWUYbpYeua0Qg4+wmgh/u/EujavkcEAftuG02jP8JrHy/VuhrVEAlyynCjtDL",
	"JD2+xoZM32xvjki+DX2PtEdSTAaj1BDLeGKA3gXldy40pWLqFdRilgY4TZMv2y8fSjo9HFQB1jpWj+Kw",
	"YEtlWEHj3eL/ghAYyyiSD20L1Otc/4McyCqwZbpTAjRNP65dfZLUZo92NHnE2kw9e66L7tKk+WUOiPjD",
	"x2pt0p5TJXwMSSfI+wEeXeC3b9QogT+8Bc6JcD9dVr9MhjHDB4Xw0641KVLw4isqpvTR/mkoBmp0rjfQ",
	"mc3+kx4ehv8DAAD//0chB3jPCwAA",
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

