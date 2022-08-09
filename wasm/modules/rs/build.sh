#!/bin/sh

set -o errexit

build() {
	if [ $# -lt 2 ]
	then
		echo "missing arguement in build" >&2
		exit 1
	fi
	dir=$1
	target=$2
	echo "$dir $target"
	(
		cd $1
		cargo build --target $2
	)
}

# to-upper
build to-upper-raw wasm32-unknown-unknown
build to-upper-wapc wasm32-unknown-unknown
build to-upper-wasi wasm32-wasi

# k8s
build k8s-raw wasm32-unknown-unknown
build k8s-wapc wasm32-unknown-unknown
build k8s-wasi wasm32-wasi
