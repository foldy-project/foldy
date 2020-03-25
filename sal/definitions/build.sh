#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
pushd $GOPATH/src/github.com/pacedotdev/oto
  go install .
popd
templates=$GOPATH/src/github.com/pacedotdev/oto/otohttp/templates
out=../src
gen() {
  set -euo pipefail
  name=$1
  oto -template $templates/rust/${name}.rs.plush \
    -out ./${name}.rs \
    -ignore Ignorer \
    ./definitions.go
  docker run -v C:/Users/tlhavlik/go/src/github.com/foldy-project/foldy/sal/definitions:/data rust:rustfmt rustfmt /data/${name}.rs
  mv ${name}.rs $out
}
gen types
gen server
gen blocking_client

