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

# follow these steps

1. clone my github `ecommerce-gin-clean-arch` repository to your system
2. use these bash commands

``` bash commands
## these bash comands are set-up on makefilie ##

# Navigate into project directory
cd ./ecommerce-gin-clean-arch

# Install needed dependencies
make deps

# !create `.env` in these work_dir with your database details

#generate wire_gen.go for dependency injection

```
