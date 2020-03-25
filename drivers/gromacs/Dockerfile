FROM rust:latest
RUN cargo new app
WORKDIR /app
COPY Cargo.toml .
COPY Cargo.lock .
RUN cargo build
COPY . .
RUN cargo build
