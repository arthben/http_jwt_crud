// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Yohanes Catur",
            "url": "www.linkedin.com/in/yohanescatur",
            "email": "yohanescatur@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1.0/access-token": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Token"
                ],
                "summary": "Get Access Token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "application/json",
                        "name": "Content-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client key provided by server",
                        "name": "X-Client-Key",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "format: 2006-01-02T15:04:05+07:00",
                        "name": "X-Timestamp",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Generated Signature",
                        "name": "X-Signature",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.TokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Message"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.Message"
                        }
                    }
                }
            }
        },
        "/v1.0/todo": {
            "get": {
                "description": "Get list of todo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todo"
                ],
                "summary": "Todo List",
                "parameters": [
                    {
                        "type": "string",
                        "description": "application/json",
                        "name": "Content-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Key provided by server",
                        "name": "X-Client-Key",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "format: 2006-01-02T15:04:05+07:00",
                        "name": "X-Timestamp",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "123425234",
                        "name": "X-Signature",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of todo",
                        "name": "ID",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/database.TableTodos"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.Message"
                        }
                    }
                }
            },
            "post": {
                "description": "## Description \nAdd new task to todo list.\n\n## Response Code\n| HTTP  | Service | Code | Description                  |\n| ----- | ------- | ---- | -----------------------------|\n|  200  |    24   |  -   | Success                      |\n|  400  |    24   |  00  | Bad Request / Unauthorized   |\n|  400  |    24   |  01  | Invalid Field Format         |\n|  400  |    24   |  02  | Missing Mandatory Field      |\n|  500  |    24   |  00  | Internal Server Error        |",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todo"
                ],
                "summary": "Add Todo item",
                "parameters": [
                    {
                        "type": "string",
                        "description": "application/json",
                        "name": "Content-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Client Key provided by server",
                        "name": "X-Client-Key",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "format: 2006-01-02T15:04:05+07:00",
                        "name": "X-Timestamp",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "123425234",
                        "name": "X-Signature",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/todo.AddTodosRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/database.TableTodos"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Message"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.TokenResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "expiresIn": {
                    "type": "string"
                },
                "responseCode": {
                    "type": "string"
                },
                "responseMessage": {
                    "type": "string"
                },
                "tokenType": {
                    "type": "string"
                }
            }
        },
        "database.TableTodos": {
            "type": "object",
            "properties": {
                "completedDate": {
                    "type": "string"
                },
                "createdDate": {
                    "type": "string"
                },
                "detail": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "statusCompleted": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updatedDate": {
                    "type": "string"
                }
            }
        },
        "response.Message": {
            "type": "object",
            "properties": {
                "responseCode": {
                    "type": "string"
                },
                "responseMessage": {
                    "type": "string"
                }
            }
        },
        "todo.AddTodosRequest": {
            "type": "object",
            "required": [
                "title"
            ],
            "properties": {
                "detail_todo": {
                    "type": "string",
                    "maxLength": 1024
                },
                "title": {
                    "type": "string",
                    "maxLength": 256,
                    "minLength": 1
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "HTTP JWT CRUD",
	Description:      "Demonstrate HTTP with Middleware, JWT, SQLX and slog package\nResponse Code Format : HTTP Status - Service Code - Response Code",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
