#!/usr/bin/env bash
set -euo pipefail

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$ROOT_DIR/bin"
OUTPUT="$BIN_DIR/mpris-tg-status"
PACKAGE="./cmd/mpris-tg-status"

trap 'printf "\n${RED}Build interrupted.${NC}\n"; exit 1' INT

print_step() {
    echo -e "${BLUE}${BOLD}==>${NC} ${BOLD}$1${NC}"
}

check_env() {
    print_step "Environment check..."

    if ! command -v go &> /dev/null; then
        echo -e "${RED}Error: Go is not installed.${NC}"
        exit 1
    fi

    GO_VER_MINOR=$(go version | sed -n 's/.*go1\.\([0-9]*\).*/\1/p')
    if [[ -z "$GO_VER_MINOR" || "$GO_VER_MINOR" -lt 23 ]]; then
        echo -e "${YELLOW}Warning: Go 1.23+ is required. Current version: $(go version)${NC}"
        echo -e "${YELLOW}Trying to build anyway...${NC}"
    fi

    if [[ ! -d "$ROOT_DIR/cmd/mpris-tg-status" ]]; then
        echo -e "${RED}Error: source directory cmd/mpris-tg-status not found.${NC}"
        exit 1
    fi

    mkdir -p "$BIN_DIR"
}

cleanup() {
    if [[ -f "$OUTPUT" ]]; then
        print_step "Cleaning up old binary..."
        rm -f "$OUTPUT"
    fi
}

build() {
    print_step "Compiling MPRIS Telegram Business Status..."

    cd "$ROOT_DIR"
    go build -buildvcs=false -ldflags="-s -w" -o "$OUTPUT" "$PACKAGE" &
    pid=$!

    frames="/ | \ -"
    while kill -0 "$pid" 2>/dev/null; do
        for frame in $frames; do
            printf "\r  ${BLUE}%s${NC} Building..." "$frame"
            sleep 0.1
        done
    done

    if wait "$pid"; then
        printf "\r  ${GREEN}✓${NC} Build complete! \n"
    else
        printf "\r  ${RED}✗${NC} Build failed! \n"
        exit 1
    fi
}

echo -e "${BLUE}---------------------------------------${NC}"
echo -e " ${BOLD}MPRIS Telegram Business Status Build${NC}"
echo -e "${BLUE}---------------------------------------${NC}"

check_env
cleanup
build

echo -e "${BLUE}---------------------------------------${NC}"
if [[ -f "$OUTPUT" ]]; then
    SIZE=$(du -h "$OUTPUT" | cut -f1)
    echo -e "${GREEN}${BOLD}Success!${NC}"
    echo -e "Binary: ${BOLD}$OUTPUT${NC}"
    echo -e "Size:   ${BOLD}$SIZE${NC}"
    if command -v file &> /dev/null; then
        echo -e "Type:   ${BOLD}$(file -b "$OUTPUT" | cut -d, -f1)${NC}"
    fi
else
    echo -e "${RED}${BOLD}Final binary not found.${NC}"
    exit 1
fi
