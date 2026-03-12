FROM mcr.microsoft.com/devcontainers/go:1-1.25-bookworm

# Fix: Remove the broken Yarn repository that causes GPG errors during feature installation
RUN rm -f /etc/apt/sources.list.d/yarn.list \
    && apt-get update \
    && apt-get install -y curl git
