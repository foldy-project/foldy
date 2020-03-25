FROM rust:latest AS builder
ENV USER=root
RUN cargo new app
WORKDIR /app
COPY Cargo.toml .
COPY Cargo.lock .
RUN cargo build
COPY . .
RUN cargo build
RUN cargo test

FROM gromacs/gromacs:latest
WORKDIR /foldy
COPY --from=builder /app/target/debug/gromacs-driver gromacs-driver
CMD ["./gromacs-driver"]
