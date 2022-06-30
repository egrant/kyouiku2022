FROM debian:bullseye

RUN apt update && apt install -y \
  bpytop \
  && rm -rf /var/lib/apt/lists/*
