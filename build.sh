#!/usr/bin/env bash
set -e


build() {
    echo "build $name:$version ..."
    # CGO_ENABLED=0 完全使用静态编译
    go build -a -o dist/$name cmd/main.go
    docker build --tag $name:$version .
    docker save $name:$version > dist/$name-$version.tar
}


version=$(echo "$(git --no-pager tag)" | sed -n '$p')
name="apitesting"
build