#!/usr/bin/env bash
# 交叉编译 sunexchange-entry 插件并打包为各平台 Release 资产。
# 产物：dist/sunexchange-entry-<os>-<arch>[.exe]
set -euo pipefail
cd "$(dirname "$0")"

export PATH="${PATH:+/usr/local/go/bin:}$PATH"
OUT="dist"
mkdir -p "$OUT"

# 清理旧产物
rm -f "$OUT"/sunexchange-entry-*

GOBIN="$(command -v go || true)"
if [ -z "$GOBIN" ]; then
  echo "未找到 go，请先安装 Go 1.21+" >&2
  exit 1
fi

build() {
  local os="$1" arch="$2" ext="$3"
  echo "构建 $os/$arch ..."
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
    go build -trimpath -ldflags="-s -w" \
    -o "$OUT/sunexchange-entry-$os-$arch$ext" .
}

build linux amd64 ""
build linux arm64 ""
build windows amd64 ".exe"
build darwin amd64 ""
build darwin arm64 ""

echo "完成，产物位于 $OUT/"
ls -lh "$OUT"