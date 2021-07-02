#!/usr/bin/env bash
set -e


apitesting() {
    if [ $usedocker == 1 ]; then
        docker-compose up -d --force-recreate apitesting
    else
        echo "building..."
        go build -o ./dist/apitesting-dev cmd/main.go
        ./dist/apitesting-dev run --config ./config.toml --server "http" --mode dev
    fi
}

usedocker=0
if [[ $* =~ "--docker" ]]; then
    usedocker=1
fi

$1