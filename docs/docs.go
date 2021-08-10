// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/config/get": {
            "get": {
                "description": "Get gateway's system config",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Get system config",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/api/config/modify": {
            "post": {
                "description": "Modify gateway's system config, include: \"CQHTTP_API_ADDRESS\", \"CQHTTP_SECRET\", \"ADMIN_QQ\", \"CMD PREFIX\"",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "config"
                ],
                "summary": "Modify config",
                "parameters": [
                    {
                        "description": "config struct",
                        "name": "config",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.SystemConfig"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/api/node/add": {
            "post": {
                "description": "Add a proxy remote node into the gateway's selector",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Selector"
                ],
                "summary": "Add Proxy Node",
                "parameters": [
                    {
                        "description": "Proxy node's remote address, eg: 127.0.0.1:8081",
                        "name": "node",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/selector.Node"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/api/node/getAll": {
            "get": {
                "description": "Get all proxy nodes and current selector's function",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Selector"
                ],
                "summary": "Get all nodes",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/api/node/modifySelector": {
            "post": {
                "description": "Change gateway's selector method, supported: \"random\", \"round_robin\", \"hash\", \"weight\"",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Selector"
                ],
                "summary": "Modify Selector",
                "parameters": [
                    {
                        "description": "Delete node",
                        "name": "selector",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.SelectorFuncName"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/api/node/remove": {
            "post": {
                "description": "Delete a proxy node",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Selector"
                ],
                "summary": "Delete Node",
                "parameters": [
                    {
                        "description": "Node's remote address",
                        "name": "node",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/selector.Node"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.SelectorFuncName": {
            "type": "object",
            "required": [
                "func_name"
            ],
            "properties": {
                "func_name": {
                    "type": "string"
                }
            }
        },
        "main.SystemConfig": {
            "type": "object",
            "required": [
                "admin_qq",
                "cqhttp_address",
                "prefix",
                "secret"
            ],
            "properties": {
                "admin_qq": {
                    "type": "string"
                },
                "cqhttp_address": {
                    "type": "string"
                },
                "prefix": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "secret": {
                    "type": "string"
                }
            }
        },
        "response.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "selector.Node": {
            "type": "object",
            "required": [
                "remote_addr"
            ],
            "properties": {
                "remote_addr": {
                    "type": "string"
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
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
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
	swag.Register(swag.Name, &s{})
}
