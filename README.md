# poe-api-go

[![GitHub issues](https://img.shields.io/github/issues/yslinear/poe-api-go)](https://github.com/yslinear/poe-api-go/issues)
[![GitHub forks](https://img.shields.io/github/forks/yslinear/poe-api-go)](https://github.com/yslinear/poe-api-go/network)
[![GitHub stars](https://img.shields.io/github/stars/yslinear/poe-api-go)](https://github.com/yslinear/poe-api-go/stargazers)
[![GitHub license](https://img.shields.io/github/license/yslinear/poe-api-go)](https://github.com/yslinear/poe-api-go/blob/master/LICENSE)

fetch [Path of Exile API](https://www.pathofexile.com/developer/docs/api-resources) with golang

## Requirements

* PostgreSQL
* github.com/jmoiron/sqlx
* github.com/joho/godotenv
* github.com/joho/godotenv/autoload
* github.com/lib/pq

## Installation

```bash
go get github.com/jmoiron/sqlx
go get github.com/joho/godotenv
go get github.com/joho/godotenv/autoload
go get github.com/lib/pq
```

## Usage

rename `env.example` to `.env`

```bash
go run main.go
```
