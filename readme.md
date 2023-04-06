# `Ecommerce` Project using `Go programing` with `gin` framework
# Here I'm follwing `clean-code architecture` for my project

# Used Packages 
1. [GIN](github.com/gin-goin/gin) is a web framework written in Go (Golang). It features a martini-like API with performance that is up to 40 times faster thanks to httprouter. If you need performance and good productivity, you will love Gin.
2. [JWT](github.com/golang-jwt/jwt) A go (or 'golang' for search engine friendliness) implementation of JSON Web Tokens.
3. [GORM](https://gorm.io/index.html) with [PostgresSQL](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL)The fantastic ORM library for Golang aims to be developer friendly.
4. [Wire](https://github.com/google/wire) is a code generation tool that automates connecting components using dependency injection.
5. [Viper](https://github.com/spf13/viper) is a complete configuration solution for Go applications including 12-Factor apps. It is designed to work within an application, and can handle all types of configuration needs and formats.
6. [swag](https://github.com/swaggo/swag) converts Go annotations to Swagger Documentation 2.0 with [gin-swagger](https://github.com/swaggo/gin-swagger) and [swaggo files](github.com/swaggo/files)

# To Use and check my `ecommerce-gin-clean-arch` Project

# Follow these steps

1. clone my github `ecommerce-gin-clean-arch` repository to your system
2. use these bash commands

``` bash commands
## these bash comands are set-up on makefilie ##

# Step 1 :  Navigate into project directory
cd ./ecommerce-gin-clean-arch

# Step 2 : Install needed dependencies
make deps 
#or
go mod tidy

# Step 3 : Setup Env file directions are given below

# Step 4 : If you want to use swagger api documetation

make swag

# Step 5 : Run the Server
make run

# If you want to use in post man : then check localhost:8000
 
# If you want to use the swagger 
#  Open any browser and visit localhost:8000/swagger/index.html (!you should generate swagger files use `make swag`)

# !Direction for creating creating Env file
# Crate create a new " .env " file in your root directory

# database details
DB_HOST=(Your Datbase host name (localhost))
DB_NAME="Your Database name"
DB_USER="Your Database owner name"
DB_PASSWORD="Your database owner password"
DB_PORT="Your database host port number"

#JWT
JWT_CODE= "secret code that you wan't to use for jwt signature"

# Twilio
# Register an account on twilio and create a messaging service
AUTH_TOKEN="Your Twilio Authenticatino token (authe token)"
ACCOUNT_SID="Your Twilio Account SID (account sid)"
SERVICE_SID="Your Twilio SID (messaging service id)"

#Razorpay
# Create an razorpay account
RAZOR_PAY_KEY="Your Razorpay test key"
RAZOR_PAY_SECRET="your Razropay secret key"

```
