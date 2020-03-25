#!/bin/bash
set -euo pipefail
docker build -t foldy/gromacs-driver:latest -f Dockerfile ../../
docker run -it --rm foldy/gromacs-driver:latest ./gromacs-driver test --help