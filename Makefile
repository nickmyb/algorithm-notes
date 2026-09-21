.DEFAULT_GOAL := help

# make init 探测到的本机工具链位置（go / python3 装了但没进 PATH 时才会有这个文件）。
# 必须在下面的 ?= 之前 include：local.mk 用 := 赋值，赋过就轮不到默认值了。
# 命令行的 make xxx GO=... 依然优先于这里的一切。
-include local.mk

GO     ?= go
PYTHON ?= python3
JAVAC  ?= javac
JAVA   ?= java

VENV   := .venv
PYTEST := $(VENV)/bin/pytest

# make new 默认生成三门语言的骨架，用 LANGS=go,python 可以只要其中几门
LANGS ?= go,python,java

# IDE 项目只初始化一次，不绑定题号；单独初始化可指定 IDE=idea/pycharm/goland。
IDE ?= all

# 给了 ID 就只测那一道题，没给就全量。
# 解析失败（题号不存在）直接报错，免得静默退化成"测了全部"。
ifneq ($(origin ID),undefined)
ifneq ($(filter test test-go test-python test-java,$(MAKECMDGOALS)),)
PROBLEM_DIR := $(shell bash ./scripts/problem-dir.sh "$(ID)")
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
	GO="$(GO)" PYTHON="$(PYTHON)" JAVAC="$(JAVAC)" bash ./scripts/init.sh

.PHONY: ide
ide: ## 初始化独立 IDE 项目（无须 ID）；可选 IDE=idea/pycharm/goland
	"$(PYTHON)" ./scripts/ide-init.py --ide "$(IDE)"

## ---------- 写题 ----------

.PHONY: new
new: ## 新开一题，如 make new ID=1 [LANGS=go,python]
	@test -n "$(ID)" || { echo "用法: make new ID=1 [LANGS=go,python,java]"; exit 1; }
	cd ctl && "$(GO)" run . new "$(ID)" --langs "$(LANGS)"

.PHONY: readme
readme: ## 重新生成仓库根 README.md
	cd ctl && "$(GO)" run . build readme

.PHONY: readme-anon
readme-anon: ## 生成不含个人数据的 README，给 template 分支用
	cd ctl && "$(GO)" run . build readme --anonymous

## ---------- 测试 ----------

.PHONY: test
test: $(PYTEST) ## 跑测试，加 ID=94 只测一道题
	GO="$(GO)" PYTEST="$(PYTEST)" JAVAC="$(JAVAC)" JAVA="$(JAVA)" bash ./scripts/test.sh $(PROBLEM_DIR)

.PHONY: test-go
test-go: ## 跑 Go 题解测试并生成覆盖率，加 ID=94 只测一道题
	GO="$(GO)" bash ./gotest.sh $(PROBLEM_DIR)

.PHONY: test-python
test-python: $(PYTEST) ## 跑 Python 题解测试，加 ID=94 只测一道题
	PYTEST="$(PYTEST)" bash ./pytest.sh $(PROBLEM_DIR)

.PHONY: test-java
test-java: ## 编译并跑 Java 题解测试，加 ID=94 只测一道题
	JAVAC="$(JAVAC)" JAVA="$(JAVA)" bash ./javatest.sh $(PROBLEM_DIR)

# 只在缺失或依赖清单变动时重建虚拟环境，避免每次跑测试都重装
$(PYTEST): requirements-dev.txt
	"$(PYTHON)" -m venv $(VENV)
	$(VENV)/bin/pip install --quiet --upgrade pip
	$(VENV)/bin/pip install --quiet -r requirements-dev.txt
	@touch $(PYTEST)

## ---------- 代码质量 ----------

.PHONY: fmt
fmt: ## 格式化 Go 代码
	"$(GO)" fmt ./...

.PHONY: vet
vet: ## 静态检查 Go 代码
	"$(GO)" vet ./...

.PHONY: tidy
tidy: ## 整理 go.mod / go.sum
	"$(GO)" mod tidy

.PHONY: clean
clean: ## 清掉构建产物和缓存，不碰题解
	rm -rf out coverage.txt .pytest_cache ctl/.cache
	find . -name __pycache__ -type d -prune -exec rm -rf {} +
