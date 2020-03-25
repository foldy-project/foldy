#!/bin/bash
set -euo pipefail
cd definitions
./build.sh
cd ..
cargo test