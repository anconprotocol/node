// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v0/code": {
            "post": {
                "description": "Execute library smartcontracts.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "Upload hybrid smartcontracts",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v0/dagcbor": {
            "post": {
                "description": "Writes a dag-cbor block which syncs with IPFS. Returns a CID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dag-cbor"
                ],
                "summary": "Stores CBOR as dag-json",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v0/dagcbor/{cid}/{path}": {
            "get": {
                "description": "Returns CBOR",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dag-cbor"
                ],
                "summary": "Reads CBOR from a dag-cbor block",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/v0/dagjson": {
            "post": {
                "description": "Writes a dag-json block which syncs with IPFS. Returns a CID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dag-json"
                ],
                "summary": "Stores JSON as dag-json",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v0/dagjson/{cid}/{path}": {
            "get": {
                "description": "Returns JSON",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dag-json"
                ],
                "summary": "Reads JSON from a dag-json block",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/v0/file": {
            "post": {
                "description": "Writes a raw block which syncs with IPFS. Returns a CID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "Stores files",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v0/file/{cid}/{path}": {
            "get": {
                "description": "Returns JSON",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "Reads JSON from a dag-json block",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/v0/proofs": {
            "post": {
                "description": "Writes an ics23 proof",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "proofs"
                ],
                "summary": "Create",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v0/proofs/verify": {
            "post": {
                "description": "Verifies an ics23 proof",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "proofs"
                ],
                "summary": "Verifies an ics23 proofs",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v0/proofs/{cid}/{path}": {
            "get": {
                "description": "Returns JSON",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "proofs"
                ],
                "summary": "Reads an existing proof",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/v0/tip": {
            "post": {
                "description": "Writes a dag-json block which syncs with IPFS. Returns a CID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "dag-json"
                ],
                "summary": "Stores JSON as dag-json",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.4.0",
	Host:        "api.ancon.did.pa",
	BasePath:    "/v0",
	Schemes:     []string{},
	Title:       "Ancon Protocol Sync API v0.4.0",
	Description: "API",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
