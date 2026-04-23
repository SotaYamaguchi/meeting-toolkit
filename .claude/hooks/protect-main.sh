#!/bin/bash
# protect-main.sh
# PreToolUse hook: main ブランチでの git commit をブロックし、
# 自動でブランチを作成して切り替える。
#
# 入力: stdin から JSON (tool_name, tool_input)
# 出力: stdout に JSON (decision: "block"|"approve", reason)

set -euo pipefail

INPUT=$(cat)

TOOL_NAME=$(echo "$INPUT" | python3 -c "import sys,json; print(json.load(sys.stdin).get('tool_name',''))" 2>/dev/null || echo "")

if [ "$TOOL_NAME" != "Bash" ]; then
  echo '{"decision":"approve"}'
  exit 0
fi

COMMAND=$(echo "$INPUT" | python3 -c "import sys,json; print(json.load(sys.stdin).get('tool_input',{}).get('command',''))" 2>/dev/null || echo "")

# git commit コマンドかどうかを判定
if ! echo "$COMMAND" | grep -qE '(^|&&\s*|;\s*)git commit'; then
  echo '{"decision":"approve"}'
  exit 0
fi

# 現在のブランチを取得
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "")

if [ "$CURRENT_BRANCH" != "main" ]; then
  echo '{"decision":"approve"}'
  exit 0
fi

# main ブランチ上での commit を検出 → 自動ブランチ作成
# ステージされたファイルからブランチ名を推測
STAGED_FILES=$(git diff --cached --name-only 2>/dev/null || echo "")

if [ -z "$STAGED_FILES" ]; then
  # ステージされていない場合は git commit -a の可能性
  STAGED_FILES=$(git diff --name-only 2>/dev/null || echo "")
fi

# 変更内容からブランチ名のプレフィックスを決定
BRANCH_PREFIX="feat"
if echo "$STAGED_FILES" | grep -qE '(_test\.go|test_)'; then
  BRANCH_PREFIX="test"
elif echo "$STAGED_FILES" | grep -qE '(fix|bug)'; then
  BRANCH_PREFIX="fix"
elif echo "$STAGED_FILES" | grep -qE '(refactor)'; then
  BRANCH_PREFIX="refactor"
fi

# 変更されたディレクトリからトピック名を生成
TOPIC=$(echo "$STAGED_FILES" | head -1 | sed 's|/.*||' | tr '[:upper:]' '[:lower:]' | tr ' ' '-')
if [ -z "$TOPIC" ]; then
  TOPIC="update"
fi

DATE=$(date +%Y%m%d)
BRANCH_NAME="${BRANCH_PREFIX}/${TOPIC}-${DATE}"

# 同名ブランチが既にある場合はサフィックスを追加
COUNTER=1
ORIGINAL_NAME="$BRANCH_NAME"
while git rev-parse --verify "$BRANCH_NAME" >/dev/null 2>&1; do
  BRANCH_NAME="${ORIGINAL_NAME}-${COUNTER}"
  COUNTER=$((COUNTER + 1))
done

# ブランチを作成して切り替え
git checkout -b "$BRANCH_NAME" 2>/dev/null

echo "{\"decision\":\"approve\",\"reason\":\"main ブランチから '${BRANCH_NAME}' を作成して切り替えました。コミットを続行します。\"}"
