FROM rust:latest AS builder
RUN rustup toolchain install nightly \
    && rustup default nightly
ENV USER=root
RUN mkdir -p /foldy/drivers \
    && cd /foldy/drivers \
    && cargo new gromacs \
    && rm gromacs/Cargo.toml
WORKDIR /foldy
COPY sal sal
COPY drivers/gromacs/Cargo.toml drivers/gromacs/Cargo.toml
COPY drivers/gromacs/Cargo.lock drivers/gromacs/Cargo.lock
RUN cd drivers/gromacs && cargo build
COPY . .

FROM gromacs/gromacs:latest
WORKDIR /foldy
COPY --from=builder /foldy/drivers/gromacs/target/debug/gromacs-driver gromacs-driver
CMD ["./gromacs-driver"]
