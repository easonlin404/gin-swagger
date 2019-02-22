# swag

<img align="right" width="180px" src="https://raw.githubusercontent.com/swaggo/swag/master/assets/swaggo.png">

[![Travis Status](https://img.shields.io/travis/swaggo/swag/master.svg)](https://travis-ci.org/swaggo/swag)
[![Coverage Status](https://img.shields.io/codecov/c/github/swaggo/swag/master.svg)](https://codecov.io/gh/swaggo/swag)
 [![Go Report Card](https://goreportcard.com/badge/github.com/swaggo/swag)](https://goreportcard.com/badge/github.com/swaggo/swag)
[![codebeat badge](https://codebeat.co/badges/71e2f5e5-9e6b-405d-baf9-7cc8b5037330)](https://codebeat.co/projects/github-com-swaggo-swag-master)
[![Go Doc](https://godoc.org/github.com/swaggo/swagg?status.svg)](https://godoc.org/github.com/swaggo/swag)

Swag converts Go annotations to Swagger Documentation 2.0. We've created a variety of plugins for popular [Go web frameworks](#supported-web-frameworks). This allows you to quickly integrate with an existing Go project (using Swagger UI).

## Contents
 - [Getting started](#getting-started)
 - [Supported Web Frameworks](#supported-web-frameworks)
 - [How to use it with Gin](#how-to-use-it-with-gin)
 - [Implementation Status](#implementation-status)
 - [swag cli](#swag-cli)
 - [Declarative Comments Format](#declarative-comments-format)
	- [General API Info](##general-api-info)
	- [API Operation](#api-operation)
	- [Security](#security)
 - [Examples](#tips)
	- [Descriptions over multiple lines](#descriptions-over-multiple-lines)
	- [User defined structure with an array type](#user-defined-structure-with-an-array-type)
	- [Add a headers in response](#add-a-headers-in-response) 
	- [Use multiple path params](#use-multiple-path-params)
	- [Example value of struct](#example-value-of-struct)
	- [Description of struct](#description-of-struct)
	- [Override swagger type of a struct field](#Override-swagger-type-of-a-struct-field)
	- [Add extension info to struct field](#add-extension-info-to-struct-field)
	- [How to using security annotations](#how-to-using-security-annotations)
- [About the Project](#about-the-project)

## Getting started

1. Add comments to your API source code, See [Declarative Comments Format](#declarative-comments-format).

2. Download swag by using:
```sh
$ go get -u github.com/swaggo/swag/cmd/swag
```

3. Run `swag init` in the project's root folder which contains the `main.go` file. This will parse your comments and generate the required files (`docs` folder and `docs/docs.go`).
```sh
$ swag init
```

  Make sure to import the generated `docs/docs.go` so that your specific configuration gets `init`'ed. If your General API annotations do not live in `main.go`, you can let swag know with `-g` flag.
  ```sh
  swag init -g http/api.go
  ```

## Supported Web Frameworks

- [gin](http://github.com/swaggo/gin-swagger)
- [echo](http://github.com/swaggo/echo-swagger)
- [buffalo](https://github.com/swaggo/buffalo-swagger)
- [net/http](https://github.com/swaggo/http-swagger)

## How to use it with Gin

Find the example source code [here](https://github.com/swaggo/swag/tree/master/example/celler).

1. After using `swag init` to generate Swagger 2.0 docs, import the following packages:
```go
import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files
```

2. Add [General API](#general-api-info) annotations in `main.go` code:

```go
// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

func main() {
	r := gin.Default()

	c := controller.NewController()

	v1 := r.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.GET(":id", c.ShowAccount)
			accounts.GET("", c.ListAccounts)
			accounts.POST("", c.AddAccount)
			accounts.DELETE(":id", c.DeleteAccount)
			accounts.PATCH(":id", c.UpdateAccount)
			accounts.POST(":id/images", c.UploadAccountImage)
		}
    //...
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
//...
```

Additionally some general API info can be set dynamically. The generated code package `docs` exports `SwaggerInfo` variable which we can use to set the title, description, version, host and base path programatically. Example using Gin:

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"./docs" // docs is generated by Swag CLI, you have to import it.
)

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://swagger.io/terms/

func main() {

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"

	r := gin.New()

	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
```

3. Add [API Operation](#api-operation) annotations in `controller` code

``` go
package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"github.com/swaggo/swag/example/celler/model"
)

// ShowAccount godoc
// @Summary Show a account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} model.Account
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /accounts/{id} [get]
func (c *Controller) ShowAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	account, err := model.AccountOne(aid)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, account)
}

// ListAccounts godoc
// @Summary List accounts
// @Description get accounts
// @Accept  json
// @Produce  json
// @Param q query string false "name search by q"
// @Success 200 {array} model.Account
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /accounts [get]
func (c *Controller) ListAccounts(ctx *gin.Context) {
	q := ctx.Request.URL.Query().Get("q")
	accounts, err := model.AccountsAll(q)
	if err != nil {
		httputil.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

//...
```

```console
$ swag init
```

4.Run your app, and browse to http://localhost:8080/swagger/index.html. You will see Swagger 2.0 Api documents as shown below:

![swagger_index.html](https://raw.githubusercontent.com/swaggo/swag/master/assets/swagger-image.png)

## Implementation Status

[Swagger 2.0 document](https://swagger.io/docs/specification/2-0/basic-structure/)

- [x] Basic Structure
- [x] API Host and Base Path
- [x] Paths and Operations
- [x] Describing Parameters
- [x] Describing Request Body
- [x] Describing Responses
- [x] MIME Types
- [x] Authentication
  - [x] Basic Authentication
  - [x] API Keys
- [x] Adding Examples
- [x] File Upload
- [x] Enums
- [x] Grouping Operations With Tags
- [ ] Swagger Extensions

# swag cli

```sh
$ swag init -h
NAME:
   swag init - Create docs.go

USAGE:
   swag init [command options] [arguments...]

OPTIONS:
   --generalInfo value, -g value       Go file path in which 'swagger general API Info' is written (default: "main.go")
   --dir value, -d value               Directory you want to parse (default: "./")
   --swagger value, -s value           Output the swagger conf for json and yaml (default: "./docs/swagger")
   --propertyStrategy value, -p value  Property Naming Strategy like snakecase,camelcase,pascalcase (default: "camelcase")
```

# Declarative Comments Format

## General API Info

**Example**
[celler/main.go](https://github.com/swaggo/swag/blob/master/example/celler/main.go)

| annotation  | description                                | example                         |
|-------------|--------------------------------------------|---------------------------------|
| title       | **Required.** The title of the application.| // @title Swagger Example API   |
| version     | **Required.** Provides the version of the application API.| // @version 1.0  |
| description | A short description of the application.    |// @description This is a sample server celler server.         																 |
| tag.name    | Name of a tag.| // @tag.name This is the name of the tag                     |
| tag.description   | Description of the tag  | // @tag.description Cool Description         |
| tag.docs.url      | Url of the external Documentation of the tag | // @tag.docs.url https://example.com|
| tag.docs.descripiton  | Description of the external Documentation of the tag| // @tag.docs.descirption Best example documentation |
| termsOfService | The Terms of Service for the API.| // @termsOfService http://swagger.io/terms/                     |
| contact.name | The contact information for the exposed API.| // @contact.name API Support  |
| contact.url  | The URL pointing to the contact information. MUST be in the format of a URL.  | // @contact.url http://www.swagger.io/support|
| contact.email| The email address of the contact person/organization. MUST be in the format of an email address.| // @contact.email support@swagger.io                                   |
| license.name | **Required.** The license name used for the API.|// @license.name Apache 2.0|
| license.url  | A URL to the license used for the API. MUST be in the format of a URL.                       | // @license.url http://www.apache.org/licenses/LICENSE-2.0.html |
| host        | The host (name or ip) serving the API.     | // @host localhost:8080         |
| BasePath    | The base path on which the API is served. | // @BasePath /api/v1             |

## API Operation

**Example**
[celler/controller](https://github.com/swaggo/swag/tree/master/example/celler/controller)


| annotation         | description                                                                                                                |
|--------------------|----------------------------------------------------------------------------------------------------------------------------|
| description        | A verbose explanation of the operation behavior.                                                                           |
| id                 | A unique string used to identify the operation. Must be unique among all API operations.                                   |
| tags               | A list of tags to each API operation that separated by commas.                                                             |
| summary            | A short summary of what the operation does.                                                                                |
| accept             | A list of MIME types the APIs can consume. Value MUST be as described under [Mime Types](#mime-types).                     |
| produce            | A list of MIME types the APIs can produce. Value MUST be as described under [Mime Types](#mime-types).                     |
| param              | Parameters that separated by spaces. `param name`,`param type`,`data type`,`is mandatory?`,`comment` `attribute(optional)` |
| security           | [Security](#security) to each API operation.                                                                               |
| success            | Success response that separated by spaces. `return code`,`{param type}`,`data type`,`comment`                              |
| failure            | Failure response that separated by spaces. `return code`,`{param type}`,`data type`,`comment`                              |
| header             | Header in response that separated by spaces. `return code`,`{param type}`,`data type`,`comment`                              |
| router             | Path definition that separated by spaces. `path`,`[httpMethod]`                                                           |

## Mime Types

`swag` accepts all MIME Types which are in the correct format, that is, match `*/*`.
Besides that, `swag` also accepts aliases for some MIME Types as follows:

| Alias                 | MIME Type                         |
|-----------------------|-----------------------------------|
| json                  | application/json                  |
| xml                   | text/xml                          |
| plain                 | text/plain                        |
| html                  | text/html                         |
| mpfd                  | multipart/form-data               |
| x-www-form-urlencoded | application/x-www-form-urlencoded |
| json-api              | application/vnd.api+json          |
| json-stream           | application/x-json-stream         |
| octet-stream          | application/octet-stream          |
| png                   | image/png                         |
| jpeg                  | image/jpeg                        |
| gif                   | image/gif                         |



## Param Type

- object (struct)
- string (string)
- integer (int, uint, uint32, uint64)
- number (float32)
- boolean (bool)
- array

## Data Type

- string (string)
- integer (int, uint, uint32, uint64)
- number (float32)
- boolean (bool)
- user defined struct

## Security
| annotation | description | parameters | example |
|------------|-------------|------------|---------|
| securitydefinitions.basic  | [Basic](https://swagger.io/docs/specification/2-0/authentication/basic-authentication/) auth.  |                                   | // @securityDefinitions.basic BasicAuth                      |
| securitydefinitions.apikey | [API key](https://swagger.io/docs/specification/2-0/authentication/api-keys/) auth.            | in, name                          | // @securityDefinitions.apikey ApiKeyAuth                    |
| securitydefinitions.oauth2.application  | [OAuth2 application](https://swagger.io/docs/specification/authentication/oauth2/) auth.       | tokenUrl, scope                   | // @securitydefinitions.oauth2.application OAuth2Application |
| securitydefinitions.oauth2.implicit     | [OAuth2 implicit](https://swagger.io/docs/specification/authentication/oauth2/) auth.          | authorizationUrl, scope           | // @securitydefinitions.oauth2.implicit OAuth2Implicit       |
| securitydefinitions.oauth2.password     | [OAuth2 password](https://swagger.io/docs/specification/authentication/oauth2/) auth.          | tokenUrl, scope                   | // @securitydefinitions.oauth2.password OAuth2Password       |
| securitydefinitions.oauth2.accessCode   | [OAuth2 access code](https://swagger.io/docs/specification/authentication/oauth2/) auth.       | tokenUrl, authorizationUrl, scope | // @securitydefinitions.oauth2.accessCode OAuth2AccessCode   |


| parameters annotation | example                                                  |
|-----------------------|----------------------------------------------------------|
| in                    | // @in header                                            |
| name                  | // @name Authorization                                   |
| tokenUrl              | // @tokenUrl https://example.com/oauth/token             |
| authorizationurl      | // @authorizationurl https://example.com/oauth/authorize |
| scope.hoge            | // @scope.write Grants write access                      |


## Attribute

```go
// @Param enumstring query string false "string enums" Enums(A, B, C)
// @Param enumint query int false "int enums" Enums(1, 2, 3)
// @Param enumnumber query number false "int enums" Enums(1.1, 1.2, 1.3)
// @Param string query string false "string valid" minlength(5) maxlength(10)
// @Param int query int false "int valid" mininum(1) maxinum(10)
// @Param default query string false "string default" default(A)
```

It also works for the struct fields:

```go
type Foo struct {
    Bar string `minLength:"4" maxLength:"16"`
    Baz int `minimum:"10" maximum:"20" default:"15"`
    Qux []string `enums:"foo,bar,baz"`
}
```

### Available

Field Name | Type | Description
---|:---:|---
<a name="parameterDefault"></a>default | * | Declares the value of the parameter that the server will use if none is provided, for example a "count" to control the number of results per page might default to 100 if not supplied by the client in the request. (Note: "default" has no meaning for required parameters.)  See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-6.2. Unlike JSON Schema this value MUST conform to the defined [`type`](#parameterType) for this parameter.
<a name="parameterMaximum"></a>maximum | `number` | See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.1.2.
<a name="parameterMinimum"></a>minimum | `number` | See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.1.3.
<a name="parameterMaxLength"></a>maxLength | `integer` | See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.2.1.
<a name="parameterMinLength"></a>minLength | `integer` | See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.2.2.
<a name="parameterEnums"></a>enums | [\*] | See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.5.1.
<a name="parameterFormat"></a>format | `string` | The extending format for the previously mentioned [`type`](#parameterType). See [Data Type Formats](#dataTypeFormat) for further details.

### Future

Field Name | Type | Description
---|:---:|---
<a name="parameterMultipleOf"></a>multipleOf | `number` | See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.1.1.
<a name="parameterPattern"></a>pattern | `string` | See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.2.3.
<a name="parameterMaxItems"></a>maxItems | `integer` | See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.3.2.
<a name="parameterMinItems"></a>minItems | `integer` | See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.3.3.
<a name="parameterUniqueItems"></a>uniqueItems | `boolean` | See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.3.4.
<a name="parameterCollectionFormat"></a>collectionFormat | `string` | Determines the format of the array if type array is used. Possible values are: <ul><li>`csv` - comma separated values `foo,bar`. <li>`ssv` - space separated values `foo bar`. <li>`tsv` - tab separated values `foo\tbar`. <li>`pipes` - pipe separated values <code>foo&#124;bar</code>. <li>`multi` - corresponds to multiple parameter instances instead of multiple values for a single instance `foo=bar&foo=baz`. This is valid only for parameters [`in`](#parameterIn) "query" or "formData". </ul> Default value is `csv`.

## Examples

### Descriptions over multiple lines

You can add descriptions spanning multiple lines in either the general api description or routes definitions like so: 

```go
// @description This is the first line
// @description This is the second line
// @description And so forth.
```

### User defined structure with an array type

```go
// @Success 200 {array} model.Account <-- This is a user defined struct.
```

```go
package model

type Account struct {
    ID   int    `json:"id" example:"1"`
    Name string `json:"name" example:"account name"`
}
```
### Add a headers in response

```go
// @Success 200 {string} string	"ok"
// @Header 200 {string} Location "/entity/1"
// @Header 200 {string} Token "qwerty"
```

### Use multiple path params

```go
/// ...
// @Param group_id path int true "Group ID"
// @Param account_id path int true "Account ID"
// ...
// @Router /examples/groups/{group_id}/accounts/{account_id} [get]
```

### Example value of struct

```go
type Account struct {
    ID   int    `json:"id" example:"1"`
    Name string `json:"name" example:"account name"`
    PhotoUrls []string `json:"photo_urls" example:"http://test/image/1.jpg,http://test/image/2.jpg"`
}
```

### Description of struct

```go
type Account struct {
    // ID this is userid
    ID   int    `json:"id"
}
```

### Override swagger type of a struct field

```go
type TimestampTime struct {
    time.Time
}

///implement encoding.JSON.Marshaler interface
func (t *TimestampTime) MarshalJSON() ([]byte, error) {
    bin := make([]byte, 16)
    bin = strconv.AppendInt(bin[:0], t.Time.Unix(), 10)
    return bin, nil
}

func (t *TimestampTime) UnmarshalJSON(bin []byte) error {
    v, err := strconv.ParseInt(string(bin), 10, 64)
    if err != nil {
        return err
    }
    t.Time = time.Unix(v, 0)
    return nil
}
///

type Account struct {
    // Override primitive type by simply specifying it via `swaggertype` tag
    ID     sql.NullInt64 `json:"id" swaggertype:"integer"`

    // Override struct type to a primitive type 'integer' by specifying it via `swaggertype` tag
    RegisterTime TimestampTime `json:"register_time" swaggertype:"primitive,integer"`

    // Array types can be overridden using "array,<prim_type>" format
    Coeffs []big.Float `json:"coeffs" swaggertype:"array,number"`
}
```

### Add extension info to struct field

```go
type Account struct {
    ID   int    `json:"id"   extensions:"x-nullable,x-abc=def"` // extensions fields must start with "x-"
}
```

generate swagger doc as follows:

```go
"Account": {
    "type": "object",
    "properties": {
        "id": {
            "type": "string",
            "x-nullable": true,
            "x-abc": "def"
        }
    }
}
```

### How to using security annotations

General API info.

```go
// @securityDefinitions.basic BasicAuth

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information
```

Each API operation.

```go
// @Security ApiKeyAuth
```

Make it AND condition

```go
// @Security ApiKeyAuth
// @Security OAuth2Application[write, admin]
```

## About the Project
This project was inspired by [yvasiyarov/swagger](https://github.com/yvasiyarov/swagger) but we simplified the usage and added support a variety of [web frameworks](#supported-web-frameworks). Gopher image source is [tenntenn/gopher-stickers](https://github.com/tenntenn/gopher-stickers). It has licenses [creative commons licensing](http://creativecommons.org/licenses/by/3.0/deed.en).