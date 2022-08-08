#!/usr/bin/env bash
set -euo pipefail

mkdir -p tmp

readonly systems=( windows linux darwin  )
readonly archs=( amd64 arm64 )

for system in "${systems[@]}"; do
    for arch in "${archs[@]}"; do
        if [[ "$system" == "windows" ]] && [[ "$arch" == "arm64" ]]; then
            echo "windows/arm64 not yet supported by tinygo"
            continue
        fi
        echo "building $system/$arch"
        # GOOS=${system} GOARCH=${arch} tinygo build -o tmp/"${system}_${arch}" cmd/epoch/main.go
        docker run --rm -w /src -v "$(pwd)":/src -e GOOS="${system}" -e GOARCH="${arch}" tinygo/tinygo tinygo build -o "tmp/${system}_${arch}" cmd/epoch/main.go
    done
done
