<div align="center">

# ✅ Gtodo

**一个简洁优雅的命令行待办事项管理工具**

使用 Go 构建 · 基于 Cobra 框架 · 彩色表格输出

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## ✨ 功能特性

- 📝 **快速添加** — 一行命令创建待办事项，支持优先级标记
- 📋 **彩色列表** — 终端表格展示，优先级以红/黄/绿色区分
- ✅ **完成标记** — 轻松将事项标记为已完成
- 🗑️ **删除事项** — 按 ID 精准删除
- 💾 **本地持久化** — JSON 文件存储，数据透明可备份

## 📦 安装

### 从源码构建

```bash
git clone https://github.com/namezzy/gtodo.git
cd gtodo
go build -o gtodo .
```

### 使用 `go install`

```bash
go install github.com/namezzy/gtodo@latest
```

## 🚀 快速上手

```bash
# 添加事项
gtodo add "完成项目文档"
gtodo add "修复紧急 Bug" -p high
gtodo add "整理书签" -p low

# 查看待办
gtodo list

# 查看全部（含已完成）
gtodo list --all

# 完成事项
gtodo done 1

# 删除事项
gtodo delete 2
```

## 📖 命令详解

### `gtodo add [描述]`

添加一条新的待办事项。

| 参数 | 缩写 | 默认值   | 说明                        |
| ---- | ---- | -------- | --------------------------- |
| `-p` | `-p` | `medium` | 优先级：`high` `medium` `low` |

```bash
gtodo add "学习 Go 并发编程" -p high
# 已添加任务 #1: 学习 Go 并发编程 (优先级: high)
```

### `gtodo list`

以表格形式列出事项，默认只显示未完成的。

| 参数    | 说明                     |
| ------- | ------------------------ |
| `--all` | 显示所有事项（包含已完成） |

输出示例：

```
┌────┬────────┬──────────────────────┬──────────────────┬────────┐
│ ID │ 优先级 │         事项         │     创建时间     │  状态  │
├────┼────────┼──────────────────────┼──────────────────┼────────┤
│  1 │   高   │ 修复紧急 Bug         │ 2026-03-09 15:04 │  待办  │
│  2 │   中   │ 完成项目文档         │ 2026-03-09 15:05 │  待办  │
│  3 │   低   │ 整理书签             │ 2026-03-09 15:06 │  待办  │
└────┴────────┴──────────────────────┴──────────────────┴────────┘
```

### `gtodo done <id>`

将指定事项标记为已完成。

```bash
gtodo done 1
# 已将事项 #1 标记为完成 ✓
```

### `gtodo delete <id>`

永久删除指定事项。

```bash
gtodo delete 2
# 已删除事项 #2
```

## 🏗️ 项目结构

```
gtodo/
├── main.go                  # 程序入口
├── cmd/                     # CLI 命令层（Cobra）
│   ├── root.go              # 根命令定义
│   ├── add.go               # add 子命令
│   ├── list.go              # list 子命令
│   ├── done.go              # done 子命令
│   └── delete.go            # delete 子命令
├── internal/                # 内部核心逻辑
│   ├── model/
│   │   └── task.go          # Task 数据模型
│   └── storage/
│       └── json_storage.go  # JSON 文件持久化
└── docs/
    └── architecture.md      # 架构与代码解析文档
```

> 详细的代码解析与学习笔记请参阅 [docs/architecture.md](docs/architecture.md)

## 💾 数据存储

所有任务以 JSON 格式保存在用户主目录下：

```
~/.gtodo/tasks.json
```

数据格式示例：

```json
[
  {
    "id": 1,
    "description": "学习 Go 并发编程",
    "priority": "high",
    "created_at": "2026-03-09T15:04:00Z",
    "status": "todo"
  }
]
```

## 🛠️ 技术栈

| 依赖 | 用途 |
| ---- | ---- |
| [spf13/cobra](https://github.com/spf13/cobra) | CLI 命令框架 |
| [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) | 终端表格渲染 |
| [fatih/color](https://github.com/fatih/color) | 终端彩色输出 |

## 📄 开源协议

本项目基于 [MIT License](LICENSE) 开源。
