#!/usr/bin/env bash
# 即使某门语言失败，也跑完其余语言，最后给出整次运行的结论。
set -eu
cd "$(dirname "${BASH_SOURCE[0]}")/.."
TEST_LANGUAGE=All
TEST_BANNER=0
source ./scripts/test-common.sh

failed=0
for script in gotest.sh pytest.sh javatest.sh; do
    if ! bash "./$script" "$@"; then
        failed=1
    fi
done
exit "$failed"
