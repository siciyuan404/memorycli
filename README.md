# MemoryCLI

一个用 Golang 构建的 CLI 工具，通过 skill 机制让多个 AI 编程助手（Claude Code、Codex、iFlow、QwenCode、GeminiCLI 等）实现记忆获取和共享。

## 功能特性

- 存储记忆到 Redis
- 检索相关记忆
- 在多个 Agent 间共享记忆
- 列出和管理记忆
- 支持多种记忆类型（知识、偏好、对话、任务）

## 安装

### 从 Release 下载

从 [Releases](https://github.com/siciyuan404/memorycli/releases) 页面下载对应平台的二进制文件。

### 从源码构建

```bash
git clone https://github.com/siciyuan404/memorycli.git
cd memorycli
go build -o bin/memorycli .
```

## 快速开始

### 1. 启动 Redis

```bash
docker run -d --name memorycli-redis -p 6379:6379 redis/redis-stack-server:latest
```

### 2. 使用 MemoryCLI

```bash
# 存储记忆
memorycli store "用户偏好 TypeScript" --type preference --importance 0.8

# 检索记忆
memorycli recall "用户偏好" --limit 5

# 共享记忆
memorycli share mem_abc123 --to agent-gpt,agent-claude

# 列出记忆
memorycli list --type knowledge

# 删除记忆
memorycli forget mem_abc123
```

## CLI 命令

### store - 存储记忆

```bash
memorycli store <content> [flags]

Flags:
  -t, --type string         记忆类型 (knowledge|preference|conversation|task) (默认 "knowledge")
  -i, --importance float    重要性评分 (0-1) (默认 0.5)
      --tags string         标签 (逗号分隔)
      --ttl int             过期时间(秒)，0 表示永不过期
  -s, --share               共享给所有 Agent
```

### recall - 检索记忆

```bash
memorycli recall <query> [flags]

Flags:
  -l, --limit int       返回数量 (默认 10)
  -t, --type string     记忆类型
  -a, --agent string    指定 Agent ID
```

### share - 共享记忆

```bash
memorycli share <memory_id> --to <agent_ids>

Flags:
  --to strings   目标 Agent ID (逗号分隔，必填)
```

### list - 列出记忆

```bash
memorycli list [flags]

Flags:
  -a, --agent string    指定 Agent ID
  -t, --type string     记忆类型
  -l, --limit int       返回数量 (默认 20)
  -s, --shared          只显示共享记忆
```

### forget - 删除记忆

```bash
memorycli forget <memory_id>
```

## AI 工具集成

### Claude Code

将 `skill/claude-code.md` 复制到 `.claude/skills/memory.md`

### 其他 AI 工具

参考 `skill/memory.md` 创建对应的 skill 文件。

## 配置

### 环境变量

- `REDIS_URL`: Redis 连接地址，默认 `redis://localhost:6379`
- `MEMORYCLI_AGENT_ID`: 当前 Agent ID

### 配置文件

`~/.memorycli/config.yaml`

```yaml
redis:
  url: redis://localhost:6379
  password: ""
  db: 0

agent:
  id: ""

memory:
  default_ttl: 0
  max_results: 100
```

## 记忆类型

| 类型 | 说明 |
|------|------|
| knowledge | 知识类记忆，如项目信息、技术栈、API 文档等 |
| preference | 用户偏好，如代码风格、框架选择等 |
| conversation | 对话记录，重要的对话上下文 |
| task | 任务相关，如待办事项、项目进度等 |

## Docker 部署

```bash
docker-compose up -d
```

## 开发

```bash
# 安装依赖
go mod download

# 运行测试
go test -v ./...

# 构建
go build -o bin/memorycli .
```

## License

MIT
