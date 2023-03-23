// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
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
        "/": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "after user login user will seen this page with user informations",
                "tags": [
                    "Home"
                ],
                "summary": "api for showing home page of user",
                "operationId": "Home",
                "responses": {
                    "200": {
                        "description": "Welcome Home"
                    },
                    "400": {
                        "description": "Faild to load user home page"
                    }
                }
            }
        },
        "/carts": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can see all productItem that stored in cart",
                "tags": [
                    "Carts"
                ],
                "summary": "api for get all cart item of user",
                "operationId": "UserCart",
                "responses": {
                    "200": {
                        "description": "successfully got user cart items",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "500": {
                        "description": "faild to get cart items",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can inrement or drement count of a productItem in cart (min=1)",
                "tags": [
                    "Carts"
                ],
                "summary": "api for updte productItem count",
                "operationId": "UpdateCart",
                "parameters": [
                    {
                        "description": "Input Field",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.ReqCartCount"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully productItem count change on cart"
                    },
                    "400": {
                        "description": "invalid input"
                    },
                    "500": {
                        "description": "can't update the count of product item on cart"
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can add a stock in product to user cart",
                "tags": [
                    "Carts"
                ],
                "summary": "api for add productItem to user cart",
                "operationId": "AddToCart",
                "parameters": [
                    {
                        "description": "Input Field",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.ReqCart"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully productItem added to cart"
                    },
                    "400": {
                        "description": "can't add the product item into cart"
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can remove a signle productItem full quantity from cart",
                "tags": [
                    "Carts"
                ],
                "summary": "api for remove a product from cart",
                "operationId": "RemoveFromCart",
                "parameters": [
                    {
                        "description": "Input Field",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.ReqCart"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully productItem removed from cart",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "400": {
                        "description": "invalid input",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "500": {
                        "description": "can't remove product item from cart",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/carts/checkout": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can checkout user cart items",
                "tags": [
                    "Carts"
                ],
                "summary": "api for cart checkout",
                "operationId": "CheckOutCart",
                "responses": {
                    "200": {
                        "description": "successfully got checkout data",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "401": {
                        "description": "cart is empty so user can't call this api",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "500": {
                        "description": "faild to get checkout items",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/carts/place-order/:address_id": {
            "post": {
                "description": "user can place after checkout",
                "tags": [
                    "Carts"
                ],
                "summary": "api for place order of all items in user cart",
                "operationId": "PlaceOrderByCart",
                "responses": {
                    "200": {
                        "description": "successfully placed your order for COD",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "400": {
                        "description": "faild to place to order",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/login": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Enter this fields on login page post",
                "tags": [
                    "Login"
                ],
                "summary": "to get the json format for login",
                "operationId": "LoginGet",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Enter user_name/phone/email with password",
                "tags": [
                    "Login"
                ],
                "summary": "api for user login",
                "operationId": "LoginPost",
                "parameters": [
                    {
                        "description": "Input Field",
                        "name": "inputs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.LoginStruct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successfully logged in",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "400": {
                        "description": "faild to login",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "500": {
                        "description": "faild to generat JWT",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/login/otp-send": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can enter email/user_name/phone will send an otp to user phone",
                "tags": [
                    "Login"
                ],
                "summary": "api for user login with otp",
                "operationId": "LoginOtpSend",
                "parameters": [
                    {
                        "description": "Input Field",
                        "name": "inputs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.OTPLoginStruct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully Otp Send to registered number",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "400": {
                        "description": "Enter input properly",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "500": {
                        "description": "Faild to send otp",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/login/otp-verify": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "enter your otp that send to your registered number",
                "tags": [
                    "Login"
                ],
                "summary": "varify user login otp",
                "operationId": "LoginOtpVerify",
                "parameters": [
                    {
                        "description": "Input Field",
                        "name": "inputs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.OTPVerifyStruct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successfully logged in uing otp"
                    },
                    "400": {
                        "description": "invalid input otp"
                    },
                    "500": {
                        "description": "Faild to generate JWT"
                    }
                }
            }
        },
        "/logout": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can logout",
                "tags": [
                    "Logout"
                ],
                "summary": "api for user to lgout",
                "operationId": "Logout",
                "responses": {
                    "200": {
                        "description": "successfully logged out"
                    }
                }
            }
        },
        "/orders": {
            "get": {
                "description": "user can see all user order history",
                "tags": [
                    "Orders"
                ],
                "summary": "api for showing user order list",
                "operationId": "GetOrdersOfUser",
                "responses": {
                    "200": {
                        "description": "successfully got shop order list of user",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "500": {
                        "description": "faild to get user shop order list",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/orders/": {
            "put": {
                "description": "admin can change user order status",
                "tags": [
                    "Orders"
                ],
                "summary": "api for admin to change the status of order",
                "operationId": "GetAllShopOrders",
                "responses": {
                    "200": {
                        "description": "Successfully order cancelled",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "400": {
                        "description": "invalid input",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/orders/items": {
            "get": {
                "description": "user can place after checkout",
                "tags": [
                    "Orders"
                ],
                "summary": "api for show order items of a specific order",
                "operationId": "GetOrderItemsForUser",
                "responses": {
                    "200": {
                        "description": "successfully got order items",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "500": {
                        "description": "faild to get order list of user",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/profile": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can edit user details",
                "tags": [
                    "Profile"
                ],
                "summary": "api for edit user details",
                "operationId": "EditAccount",
                "parameters": [
                    {
                        "description": "input field",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.ReqUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successfully edited user details",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "400": {
                        "description": "invalid input",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/profile/": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can see user details",
                "tags": [
                    "Profile"
                ],
                "summary": "api for showing user details",
                "operationId": "Account",
                "responses": {
                    "200": {
                        "description": "Successfully user account details found"
                    },
                    "500": {
                        "description": "faild to show user details",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/profile/address": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can show all adderss",
                "tags": [
                    "Address"
                ],
                "summary": "api for get all address of user",
                "operationId": "GetAddresses",
                "responses": {
                    "200": {
                        "description": "successfully got user addresses",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "500": {
                        "description": "faild to show user addresses",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can change existing address",
                "tags": [
                    "Address"
                ],
                "summary": "api for edit user address",
                "operationId": "EditAddress",
                "parameters": [
                    {
                        "description": "Input Field",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.ReqEditAddress"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successfully addresses updated",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "400": {
                        "description": "can't update the address",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get a new address from user to store the the database",
                "tags": [
                    "Address"
                ],
                "summary": "api for adding a new address for user",
                "operationId": "AddAddress",
                "parameters": [
                    {
                        "description": "Input Field",
                        "name": "inputs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/req.ReqAddress"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully address added",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "400": {
                        "description": "inavlid input",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        },
        "/signup": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can see what are the fields to enter to create a new account",
                "tags": [
                    "Signup"
                ],
                "summary": "api for user to signup page",
                "operationId": "SignUpGet",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "user can send user details and validate and create new account",
                "tags": [
                    "Signup"
                ],
                "summary": "api for user to post the user details",
                "operationId": "SignUpPost",
                "responses": {
                    "200": {
                        "description": "Successfully account created"
                    },
                    "400": {
                        "description": "Faild to create account"
                    }
                }
            }
        },
        "/wishlist": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Wishlist"
                ],
                "summary": "api get all wish list items of user",
                "operationId": "GetWishListI",
                "responses": {
                    "200": {
                        "description": "Wish list is empty"
                    },
                    "400": {
                        "description": "faild to get user wish list items"
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Wishlist"
                ],
                "summary": "api to add a productItem to wish list",
                "operationId": "AddToWishList",
                "parameters": [
                    {
                        "description": "product_id",
                        "name": "product_id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successfully added product item to wishlist",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "400": {
                        "description": "invalid input",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Wishlist"
                ],
                "summary": "api to remove a productItem from wish list",
                "operationId": "RemoveFromWishList",
                "responses": {
                    "200": {
                        "description": "successfully removed product item from wishlist",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    },
                    "400": {
                        "description": "invalid input",
                        "schema": {
                            "$ref": "#/definitions/res.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "req.LoginStruct": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 30,
                    "minLength": 5
                },
                "phone": {
                    "type": "string",
                    "maxLength": 10,
                    "minLength": 10
                },
                "user_name": {
                    "type": "string",
                    "maxLength": 15,
                    "minLength": 3
                }
            }
        },
        "req.OTPLoginStruct": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "phone": {
                    "type": "string",
                    "maxLength": 10,
                    "minLength": 10
                },
                "user_name": {
                    "type": "string",
                    "maxLength": 16,
                    "minLength": 3
                }
            }
        },
        "req.OTPVerifyStruct": {
            "type": "object",
            "required": [
                "id",
                "otp"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                },
                "otp": {
                    "type": "string",
                    "maxLength": 8,
                    "minLength": 4
                }
            }
        },
        "req.ReqAddress": {
            "type": "object",
            "required": [
                "country_id",
                "house",
                "land_mark",
                "name",
                "phone_number",
                "pincode"
            ],
            "properties": {
                "area": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "country_id": {
                    "type": "integer"
                },
                "house": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_default": {
                    "type": "boolean"
                },
                "land_mark": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 2
                },
                "phone_number": {
                    "type": "string",
                    "maxLength": 10,
                    "minLength": 10
                },
                "pincode": {
                    "type": "integer"
                }
            }
        },
        "req.ReqCart": {
            "type": "object",
            "required": [
                "product_item_id"
            ],
            "properties": {
                "product_item_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "req.ReqCartCount": {
            "type": "object",
            "required": [
                "increment",
                "product_item_id"
            ],
            "properties": {
                "count": {
                    "type": "integer",
                    "minimum": 1
                },
                "increment": {
                    "type": "boolean"
                },
                "product_item_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "req.ReqEditAddress": {
            "type": "object",
            "required": [
                "country_id",
                "house",
                "id",
                "land_mark",
                "name",
                "phone_number",
                "pincode"
            ],
            "properties": {
                "area": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "country_id": {
                    "type": "integer"
                },
                "house": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_default": {
                    "type": "boolean"
                },
                "land_mark": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 2
                },
                "phone_number": {
                    "type": "string",
                    "maxLength": 10,
                    "minLength": 10
                },
                "pincode": {
                    "type": "integer"
                }
            }
        },
        "req.ReqUpdateOrder": {
            "type": "object",
            "required": [
                "shop_order_id"
            ],
            "properties": {
                "order_status_id": {
                    "type": "integer"
                },
                "shop_order_id": {
                    "type": "integer"
                }
            }
        },
        "req.ReqUser": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 2
                },
                "last_name": {
                    "type": "string",
                    "maxLength": 50,
                    "minLength": 1
                },
                "phone": {
                    "type": "string",
                    "maxLength": 10,
                    "minLength": 10
                },
                "user_name": {
                    "type": "string",
                    "maxLength": 15,
                    "minLength": 3
                }
            }
        },
        "res.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "errors": {},
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
