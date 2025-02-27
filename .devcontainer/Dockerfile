FROM mcr.microsoft.com/devcontainers/go:1.21

RUN apt-get update && apt-get install -y \
    git \
    curl \
    make \
    && rm -rf /var/lib/apt/lists/*

USER vscode
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download 