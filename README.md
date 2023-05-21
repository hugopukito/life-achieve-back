# Life Achieve

## Services for dev

### Install sql and set database

sudo apt update

sudo apt install -y mysql-server

sudo mysql

ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'password';

FLUSH PRIVILEGES;

mysql -u root -p
password -> 'password'

CREATE DATABASE IF NOT EXISTS life_achieve;

exit

### Load database file -> not existing yet

mysql -u root -p life_achieve < life_achieve.sql

### Create dir with secret for jwt tokens

nano secret_jwt.txt

Then add your long password/secret in it

## Run api

### dev

go run main.go

OR

go build main.go && ./main

## Swagger Doc

http://localhost:{PORT}/swagger/index.html

### Install (if docs folder does not exist)

go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/swag/cmd/swag

### .bashrc

export PATH=$(go env GOPATH)/bin:$PATH

### Generate new doc

swag init

## Tests

run repository tests
-> go test test/repository/*

## Database schema

| User | type | value |
| ------ | ------ | ------ |
| id | binary(16)
| type | varchar(255) | 'employer' or 'candidate'
| first_name | varchar(255) |
| last_name | varchar(255) |
| email | varchar(255) |
| password | varchar(255) |
