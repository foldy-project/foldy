FROM rust:latest
RUN rustup component add rustfmt --toolchain 1.42.0-x86_64-unknown-linux-gnu
