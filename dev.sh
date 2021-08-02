#!/usr/bin/env bash
set -e

go build -o ./dist/apitesting-dev main.go
./dist/apitesting-dev --config ./config.toml
