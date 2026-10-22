#!/bin/sh
# Install GoNime from GitHub Releases.
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/dinhtuandev/GoNime/main/install.sh | sh
#   VERSION=v1.2.3 DIR=$HOME/.local/bin sh install.sh
set -eu

REPO="dinhtuandev/GoNime"
BIN="gonime"
VERSION="${VERSION:-latest}"
DIR="${DIR:-}"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
	linux) os=linux ;;
	darwin) os=darwin ;;
	*)
		echo "GoNime: unsupported OS '$os'" >&2
		exit 1
		;;
esac

arch=$(uname -m)
case "$arch" in
	x86_64 | amd64) arch=amd64 ;;
	arm64 | aarch64) arch=arm64 ;;
	*)
		echo "GoNime: unsupported architecture '$arch'" >&2
		exit 1
		;;
esac

if [ "$VERSION" = "latest" ]; then
	tag=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" |
		sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)
else
	tag="$VERSION"
fi

if [ -z "$tag" ]; then
	echo "GoNime: cannot resolve a release version" >&2
	exit 1
fi

ver=${tag#v}
asset="${BIN}_${ver}_${os}_${arch}.tar.gz"
base="https://github.com/$REPO/releases/download/$tag"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading $asset ($tag)"
curl -fsSL "$base/$asset" -o "$tmp/$asset"
curl -fsSL "$base/checksums.txt" -o "$tmp/checksums.txt"

if command -v sha256sum >/dev/null 2>&1; then
	sha() { sha256sum "$1" | awk '{print $1}'; }
elif command -v shasum >/dev/null 2>&1; then
	sha() { shasum -a 256 "$1" | awk '{print $1}'; }
else
	echo "GoNime: sha256sum or shasum is required" >&2
	exit 1
fi

expected=$(awk -v a="$asset" '$2 == a {print $1}' "$tmp/checksums.txt")
if [ -z "$expected" ]; then
	echo "GoNime: no checksum published for $asset" >&2
	exit 1
fi

actual=$(sha "$tmp/$asset")
if [ "$expected" != "$actual" ]; then
	echo "GoNime: checksum mismatch for $asset" >&2
	exit 1
fi

tar -xzf "$tmp/$asset" -C "$tmp" "$BIN"

if [ -z "$DIR" ]; then
	if [ -w /usr/local/bin ]; then
		DIR=/usr/local/bin
	else
		DIR="$HOME/.local/bin"
	fi
fi

mkdir -p "$DIR"
install -m 0755 "$tmp/$BIN" "$DIR/$BIN"

echo "Installed $DIR/$BIN"
case ":$PATH:" in
*":$DIR:"*) ;;
*) echo "Add $DIR to your PATH to run 'gonime'." ;;
esac

"$DIR/$BIN" version
