# Memory Skill

一个用于 AI Agent 记忆管理的 skill，支持记忆存储、检索、共享和同步。

## 功能

- 存储记忆到 Redis
- 检索相关记忆
- 在多个 Agent 间共享记忆
- 列出和管理记忆

## 工具

### memory_recall

检索与查询相关的记忆。

**参数:**
- `query` (string, required): 搜索查询
- `limit` (int, optional): 返回结果数量，默认 10
- `agent_id` (string, optional): 指定 Agent ID
- `memory_type` (string, optional): 记忆类型 (knowledge|preference|conversation|task)

**示例:**
```bash
memorycli recall "用户偏好 TypeScript" --limit 5
memorycli recall "项目配置" --type knowledge --json
```

### memory_store

存储新的记忆。

**参数:**
- `content` (string, required): 记忆内容
- `memory_type` (string, optional): 记忆类型，默认 "knowledge"
- `importance` (float, optional): 重要性评分 0-1，默认 0.5
- `tags` (string, optional): 标签列表，逗号分隔
- `ttl` (int, optional): 过期时间（秒），0 表示永不过期
- `share` (bool, optional): 是否共享给所有 Agent

**示例:**
```bash
memorycli store "用户正在开发 FastAPI 项目" --type knowledge --importance 0.8 --tags "project,python"
memorycli store "项目使用 FastAPI" --tags "project,python" --share
```

### memory_share

将记忆共享给其他 Agent。

**参数:**
- `memory_id` (string, required): 记忆 ID
- `target_agents` (string, required): 目标 Agent ID 列表，逗号分隔

**示例:**
```bash
memorycli share mem_abc123 --to agent-gpt,agent-claude
```

### memory_list

列出记忆。

**参数:**
- `agent_id` (string, optional): 指定 Agent ID
- `memory_type` (string, optional): 记忆类型
- `limit` (int, optional): 返回数量，默认 20
- `shared` (bool, optional): 只显示共享记忆

**示例:**
```bash
memorycli list --type knowledge --limit 10
memorycli list --shared
```

### memory_forget

删除记忆。

**参数:**
- `memory_id` (string, required): 记忆 ID

**示例:**
```bash
memorycli forget mem_abc123
```

## 配置

环境变量:
- `REDIS_URL`: Redis 连接地址，默认 `redis://localhost:6379`
- `MEMORYCLI_AGENT_ID`: 当前 Agent ID

配置文件: `~/.memorycli/config.yaml`

## 使用示例

### Claude Code 使用

```
用户: 帮我创建一个 FastAPI 项目

Claude Code 调用:
1. memorycli recall "FastAPI 项目 用户偏好"
2. 获取记忆: "用户偏好使用 TypeScript"
3. 创建项目后存储: memorycli store "创建了 FastAPI 项目..."
```

### 多 Agent 协作

```
Claude Code:
  memorycli store "项目使用 FastAPI + TypeScript" --share

Codex:
  memorycli recall "项目技术栈"
  → 获取共享记忆，了解项目背景

iFlow:
  memorycli recall "项目技术栈"
  → 继续基于已有记忆工作
```

## 记忆类型说明

- **knowledge**: 知识类记忆，如项目信息、技术栈、API 文档等
- **preference**: 用户偏好，如代码风格、框架选择等
- **conversation**: 对话记录，重要的对话上下文
- **task**: 任务相关，如待办事项、项目进度等

## 重要性评分

- 0.0 - 0.3: 低重要性，临时信息
- 0.4 - 0.6: 中等重要性，一般信息
- 0.7 - 0.9: 高重要性，关键信息
- 1.0: 最高重要性，核心记忆
