.DEFAULT_GOAL := help

GO     ?= go
PYTHON ?= python3

VENV   := .venv
PYTEST := $(VENV)/bin/pytest

# make new 默认生成三门语言的骨架，用 LANGS=go,python 可以只要其中几门
LANGS ?= go,python,java

# 给了 ID 就只测那一道题，没给就全量。
# 解析失败（题号不存在）直接报错，免得静默退化成"测了全部"。
ifdef ID
ifneq ($(filter test test-go test-python test-java,$(MAKECMDGOALS)),)
PROBLEM_DIR := $(shell bash ./scripts/problem-dir.sh $(ID) 2>/dev/null)
ifeq ($(PROBLEM_DIR),)
$(error 找不到题号 $(ID) 对应的目录，先跑 make new ID=$(ID))
endif
endif
endif

.PHONY: help
help: ## 列出所有命令
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

## ---------- 初始化 ----------

.PHONY: init
init: ## 初始化仓库：检查工具链、装依赖、跑通测试、生成 README
	bash ./scripts/init.sh

## ---------- 写题 ----------

.PHONY: new
new: ## 新开一题，如 make new ID=1 [LANGS=go,python]
	@test -n "$(ID)" || { echo "用法: make new ID=1 [LANGS=go,python,java]"; exit 1; }
	cd ctl && $(GO) run . new $(ID) --langs $(LANGS)

.PHONY: readme
readme: ## 重新生成仓库根 README.md
	cd ctl && $(GO) run . build readme

## ---------- 测试 ----------

.PHONY: test
test: test-go test-python test-java ## 跑测试，加 ID=94 只测一道题

.PHONY: test-go
test-go: ## 跑 Go 题解测试并生成覆盖率，加 ID=94 只测一道题
	bash ./gotest.sh $(PROBLEM_DIR)

.PHONY: test-python
test-python: $(PYTEST) ## 跑 Python 题解测试，加 ID=94 只测一道题
	PYTEST=$(PYTEST) bash ./pytest.sh $(PROBLEM_DIR)

.PHONY: test-java
test-java: ## 编译并跑 Java 题解测试，加 ID=94 只测一道题
	bash ./javatest.sh $(PROBLEM_DIR)

# 只在缺失或依赖清单变动时重建虚拟环境，避免每次跑测试都重装
$(PYTEST): requirements-dev.txt
	$(PYTHON) -m venv $(VENV)
	$(VENV)/bin/pip install --quiet --upgrade pip
	$(VENV)/bin/pip install --quiet -r requirements-dev.txt
	@touch $(PYTEST)

## ---------- 代码质量 ----------

.PHONY: fmt
fmt: ## 格式化 Go 代码
	$(GO) fmt ./...

.PHONY: vet
vet: ## 静态检查 Go 代码
	$(GO) vet ./...

.PHONY: tidy
tidy: ## 整理 go.mod / go.sum
	$(GO) mod tidy

.PHONY: clean
clean: ## 清掉构建产物和缓存，不碰题解
	rm -rf out coverage.txt .pytest_cache ctl/.cache
	find . -name __pycache__ -type d -prune -exec rm -rf {} +
