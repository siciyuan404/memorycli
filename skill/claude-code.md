# Memory Skill

使用 memorycli 工具管理记忆。

## 工具调用

当需要存储或检索记忆时，执行以下命令：

```bash
# 检索记忆
memorycli recall "查询内容" --json

# 存储记忆
memorycli store "记忆内容" --type knowledge

# 共享记忆
memorycli share <memory_id> --to <agent_ids>

# 列出记忆
memorycli list --type knowledge

# 删除记忆
memorycli forget <memory_id>
```

## 常用场景

### 1. 用户偏好记忆
```bash
# 记住用户的代码偏好
memorycli store "用户喜欢简洁的代码，偏好类型注解" --type preference --importance 0.8

# 检索用户偏好
memorycli recall "用户偏好"
```

### 2. 项目信息记忆
```bash
# 记住项目信息
memorycli store "项目使用 FastAPI + TypeScript" --type knowledge --tags "project,fastapi,typescript"

# 记住 API 设计
memorycli store "用户偏好 RESTful API 设计风格" --type knowledge --importance 0.7
```

### 3. 协作记忆
```bash
# 共享重要信息给其他 Agent
memorycli share <memory_id> --to agent-gpt,agent-codex

# 查看所有共享记忆
memorycli list --shared
```

## 输出格式

建议使用 `--json` 标志获取结构化输出，便于解析：

```bash
memorycli recall "查询" --json
```

返回 JSON 格式：
```json
[
  {
    "id": "mem_abc123",
    "agent_id": "claude-code",
    "type": "preference",
    "content": "用户喜欢简洁的代码",
    "importance": 0.8,
    "tags": ["preference"],
    "created_at": "2026-03-11T00:00:00Z"
  }
]
```
