// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/authors": {
            "get": {
                "description": "there can be multiple search parameters, id or translit, or both\nexample: ?id=1\u0026id=2\u0026translit=author-21\u0026id=3\u0026translit=author-756342",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authors"
                ],
                "summary": "Find authors",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Author Id, must be a primitive id hex",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Translit author name",
                        "name": "translit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Authors",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/authors_pb.GetAuthorsResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Field was not in a correct format",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Authors not found",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Some internal error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "503": {
                        "description": "Service does not responding (maybe crush)",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/genres/all": {
            "get": {
                "description": "Fetches all categories (with genres included)",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "genres"
                ],
                "summary": "Get all genres",
                "responses": {
                    "200": {
                        "description": "Genres fetched successfully",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/genres_pb.CategoryModel"
                            }
                        }
                    },
                    "404": {
                        "description": "There's no genres in database",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal error occured",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "503": {
                        "description": "Service does not responding (maybe crush)",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "params goes in specific order: id -\u003e login -\u003e email\nfirst found user will be returned. If no user found, there'll be an error with details",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Find user by credentials",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User Id",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "User login",
                        "name": "login",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "email",
                        "description": "User email",
                        "name": "email",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User DTO model",
                        "schema": {
                            "$ref": "#/definitions/users_pb.UserModel"
                        }
                    },
                    "400": {
                        "description": "Request's field was not in a correct format",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "503": {
                        "description": "Service does not responding (maybe crush)",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/users/auth": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "check if current user has legit token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Authenticates user",
                "responses": {
                    "200": {
                        "description": "User successfully authorized",
                        "schema": {
                            "$ref": "#/definitions/user.UserAuthenticate.UserResponse"
                        }
                    },
                    "401": {
                        "description": "User does not authorized",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "User does not exists in database",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "503": {
                        "description": "Service does not responding (maybe crush)",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/users/login": {
            "post": {
                "description": "log in user with provided login and password",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Authorizes user",
                "parameters": [
                    {
                        "description": "Login field can be presented as login and email as well",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users_pb.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User successfully authorized",
                        "schema": {
                            "$ref": "#/definitions/user.UserLogin.UserResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request data",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "503": {
                        "description": "Service does not responding (maybe crush)",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/users/signin": {
            "post": {
                "description": "creates new user and authorizes it",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Registration",
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users_pb.RegistrationRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User registered and authorized",
                        "schema": {
                            "$ref": "#/definitions/user.UserRegister.UserResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request data",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Some internal error occured",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "503": {
                        "description": "Service does not responding (maybe crush)",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/middleware.CustomError"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "details": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/shared_pb.ErrorDetail"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "authors_pb.AuthorModel": {
            "type": "object",
            "properties": {
                "about": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "profilepicture": {
                    "type": "string"
                },
                "rating": {
                    "type": "number"
                },
                "translitname": {
                    "type": "string"
                }
            }
        },
        "authors_pb.GetAuthorsResponse": {
            "type": "object",
            "properties": {
                "authors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/authors_pb.AuthorModel"
                    }
                }
            }
        },
        "genres_pb.CategoryModel": {
            "type": "object",
            "properties": {
                "genres": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/genres_pb.GenreModel"
                    }
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "translitName": {
                    "type": "string"
                }
            }
        },
        "genres_pb.GenreModel": {
            "type": "object",
            "properties": {
                "bookCount": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "translitName": {
                    "type": "string"
                }
            }
        },
        "middleware.CustomError": {
            "description": "General error object. This structure always returns when error occured",
            "type": "object",
            "properties": {
                "code": {
                    "description": "Internal gRPC error code (e.g. 3)",
                    "type": "integer",
                    "example": 3
                },
                "details": {
                    "description": "Error details. Check 'ErrorDetail' structure for more information",
                    "type": "array",
                    "items": {}
                },
                "message": {
                    "description": "Error message. Can be shown to users",
                    "type": "string",
                    "example": "Bad token provided"
                },
                "type": {
                    "description": "Error code in string (e.g. InvalidArgument)",
                    "type": "string",
                    "example": "InvalidArgument"
                }
            }
        },
        "shared_pb.ErrorDetail": {
            "description": "Error detail contains information about error",
            "type": "object",
            "properties": {
                "actualvalue": {
                    "description": "Actual value of field that causes the error. Note: 'password' field will be hidden",
                    "type": "string",
                    "example": "token"
                },
                "description": {
                    "description": "Error description. Only development purposes, do not show users",
                    "type": "string",
                    "example": "Field must be a jwt token"
                },
                "field": {
                    "description": "Field that error occured on",
                    "type": "string",
                    "example": "Token"
                },
                "struct": {
                    "description": "Structure that contains field",
                    "type": "string",
                    "example": "users_pb.TokenRequest"
                },
                "tag": {
                    "description": "Failed validation tag",
                    "type": "string",
                    "example": "jwt"
                }
            }
        },
        "user.UserAuthenticate.UserResponse": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string",
                    "example": "admin"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "user"
                    ]
                }
            }
        },
        "user.UserLogin.UserResponse": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string",
                    "example": "admin"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "user"
                    ]
                }
            }
        },
        "user.UserRegister.UserResponse": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string",
                    "example": "admin"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "user"
                    ]
                }
            }
        },
        "users_pb.LoginRequest": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "description": "Can be presented as login or email",
                    "type": "string",
                    "example": "admin"
                },
                "password": {
                    "type": "string",
                    "example": "admin"
                }
            }
        },
        "users_pb.RegistrationRequest": {
            "type": "object",
            "required": [
                "email",
                "login",
                "password",
                "password_repeat"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "login": {
                    "type": "string",
                    "maxLength": 16,
                    "minLength": 4
                },
                "password": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 8
                },
                "password_repeat": {
                    "type": "string"
                }
            }
        },
        "users_pb.UserModel": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "Cookie"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9000",
	BasePath:         "/api/v1/",
	Schemes:          []string{},
	Title:            "API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
