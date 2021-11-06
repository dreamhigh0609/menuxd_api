<h1 align="center">
  <a href="http://devsstore.net"><img src="https://devsstore.net/img/logo.d7ca7ff6.png" alt="DevsStore Logo" width="300"></a>
</h1>

# MenuXD API Restful

Restful application for the Menu system. Complete administration system for restaurants and digital menu.

![PostgreSQL](https://img.shields.io/badge/DevsStore-PostgreSQL-darkblue.svg?logo=postgresql&longCache=true&style=flat) ![Go badge](https://img.shields.io/badge/DevsStore-golang-blue.svg?logo=go&longCache=true&style=flat)

## Getting Started

This project uses the **Go** programming language (Golang) and the **PostgreSQL** database engine.

### Prerequisites

[PostgreSQL](https://www.postgresql.org/) is required in version 9.6 or higher and [Go](https://golang.org/) at least in version 1.12

### Installing

The following dependencies are required:

* github.com/go-chi/chi

* github.com/go-chi/cors

* github.com/go-chi/jwtauth

* github.com/dgrijalva/jwt-go

* github.com/jinzhu/gorm

* github.com/joho/godotenv

* github.com/mailjet/mailjet-apiv3-go

* github.com/sethvargo/go-password

* github.com/stretchr/testify

* golang.org/x/crypto

* gopkg.in/olahol/melody.v1


#### Using GOPATH
```
go get github.com/go-chi/chi

go get github.com/go-chi/cors

go get github.com/go-chi/jwtauth

go get github.com/dgrijalva/jwt-go

go get github.com/jinzhu/gorm

go get github.com/joho/godotenv

go get github.com/mailjet/mailjet-apiv3-go

go get github.com/sethvargo/go-password

go get github.com/stretchr/testify

go get golang.org/x/crypto

go get gopkg.in/olahol/melody.v1

```

#### Using GOMODULE
```
go build ./cmd/menuxd
```

## Running the tests

```
go test ./...
```

## Deployment

Clone the repository
```
git clone git@gitlab.com:menuxd/api-rest.git
```
Enter the repository folder
```
cd api-rest
```
Build the binary
```
go build ./cmd/menuxd/
```
Run the program
```
# In Unix-like OS
./menuxd

# In Windows
menuxd.exe

# Debug Mode in Unix-like OS
./menuxd -debug

# Debug Mode in Windows
menuxd.exe -debug
```

### API Documentation
[Swagger](https://app.swaggerhub.com/apis/orlmonteverde/MenuxD/1.5.0)

## Built With

* [chi](https://github.com/go-chi/chi) - lightweight, idiomatic and composable router for building Go HTTP services.
* [gorm](https://gorm.io/) - The fantastic ORM library for Golang.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://gitlab.com/menuxd/api-rest/-/tags).

## Authors

* **Orlando Monteverde** - *Initial work* - [orlmonteverde](https://github.com/orlmonteverde)

See also the list of [contributors](https://gitlab.com/menuxd/api-rest/-/graphs/master) who participated in this project.

