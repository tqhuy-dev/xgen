#!/usr/bin/env bash
# Tìm symbol theo --name trong code_base.json (cùng thư mục script).
# Mặc định in JSON array "compact" (ít field, line_code cắt ngắn) để tiết kiệm token.
# --full: in nguyên object từ index (docstring, constant_value, calls_to, ...).
#
# Usage:
#   ./find_code_destination_in_json.sh --name <SymbolName>
#   ./find_code_destination_in_json.sh --name <SymbolName> --full

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
JSON_FILE="${SCRIPT_DIR}/cb.json"
# Độ dài tối đa line_code khi chế độ compact (ký tự Unicode; jq xử lý UTF-8).
LINE_CODE_MAX="${LINE_CODE_MAX:-120}"

usage() {
  echo "Usage: $(basename "$0") --name <symbol_name> [--full|--compact]" >&2
  echo "  Đọc symbols từ ${JSON_FILE} và in mảng JSON các phần tử có .name khớp." >&2
  echo "  Mặc định: --compact (chỉ name, kind, file_path, line_start, line_end, address, line_code đã cắt)." >&2
  echo "  --full: toàn bộ field như trong index (có thể rất lớn)." >&2
}

NAME=""
MODE="compact"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --name)
      if [[ $# -lt 2 ]]; then
        echo "error: --name cần một giá trị" >&2
        usage
        exit 1
      fi
      NAME="$2"
      shift 2
      ;;
    --full)
      MODE="full"
      shift
      ;;
    --compact)
      MODE="compact"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "error: tùy chọn không hợp lệ: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$NAME" ]]; then
  echo "error: thiếu --name" >&2
  usage
  exit 1
fi

if [[ ! -f "$JSON_FILE" ]]; then
  echo "error: không thấy file: $JSON_FILE" >&2
  exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
  echo "error: cần cài jq (brew install jq)" >&2
  exit 1
fi

if [[ "$MODE" == "full" ]]; then
  jq --arg name "$NAME" '[.symbols[] | select(.name == $name)]' "$JSON_FILE"
else
  jq --arg name "$NAME" --argjson max "$LINE_CODE_MAX" '
    [.symbols[] | select(.name == $name) | {
      name,
      kind,
      file_path,
      line_start,
      line_end,
      address,
      line_code: (
        if .line_code == null then null
        elif (.line_code | length) <= $max then .line_code
        else (.line_code[0:$max] + "…")
        end
      )
    }]
  ' "$JSON_FILE"
fi
