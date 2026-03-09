# Gtodo 项目结构与代码解析

本文档详细解释项目的目录结构、各模块职责以及核心代码逻辑，帮助你理解整个项目的设计思路。

## 目录结构总览

```
Gtodo/
├── main.go                        # 程序入口
├── go.mod                         # Go 模块定义与依赖声明
├── go.sum                         # 依赖校验文件（自动生成）
├── LICENSE                        # 开源协议
├── README.md                      # 使用文档
├── cmd/                           # 命令层 —— 所有 CLI 子命令的定义
│   ├── root.go                    # 根命令，Cobra 框架初始化
│   ├── add.go                     # add 子命令：添加事项
│   ├── list.go                    # list 子命令：列出事项
│   ├── done.go                    # done 子命令：完成事项
│   └── delete.go                  # delete 子命令：删除事项
├── internal/                      # 内部包 —— 不对外暴露
│   ├── model/
│   │   └── task.go                # 数据模型：Task 结构体与状态常量
│   └── storage/
│       └── json_storage.go        # 持久化层：JSON 文件读写
└── docs/
    └── architecture.md            # 本文档
```

## 分层架构

项目采用经典的三层分离架构：

```
┌─────────────────────────────────┐
│           main.go               │  入口：调用 cmd.Execute()
├─────────────────────────────────┤
│           cmd/ 层               │  命令层：解析参数、调度逻辑
│  root.go  add.go  list.go ...   │
├─────────────────────────────────┤
│         internal/ 层            │  核心层：数据模型 + 持久化
│   model/task.go                 │  ← 定义 Task 结构
│   storage/json_storage.go       │  ← 读写 JSON 文件
└─────────────────────────────────┘
```

**为什么这样分层？**

- `cmd/` 只关心"用户输入了什么、该调用谁"，不关心数据怎么存。
- `internal/model/` 只定义"数据长什么样"，不关心命令和存储。
- `internal/storage/` 只关心"怎么读写文件"，不关心上层命令。
- 这样每一层职责单一，修改一个不会影响另一个。例如将来想把 JSON 存储换成 SQLite，只需改 `storage/` 即可。

---

## 各文件详解

### 1. `main.go` — 程序入口

```go
package main

import "github.com/namezzy/gtodo/cmd"

func main() {
    cmd.Execute()
}
```

**作用：** 整个程序的入口点，只做一件事 —— 调用 `cmd.Execute()` 启动 Cobra 命令框架。

**知识点：** Go 程序从 `main` 包的 `main()` 函数开始执行。这里把所有逻辑交给 `cmd` 包，保持入口文件极简。

---

### 2. `cmd/root.go` — 根命令

```go
var rootCmd = &cobra.Command{
    Use:   "gtodo",
    Short: "A brief description of your application",
}

func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}
```

**作用：**

- 定义 `gtodo` 这个根命令（即不带子命令时的行为）。
- `Execute()` 是 Cobra 的标准入口，解析命令行参数并路由到对应的子命令。
- `init()` 中可以定义全局 flag。

**知识点：** [Cobra](https://github.com/spf13/cobra) 是 Go 生态最流行的 CLI 框架，kubectl、Hugo、GitHub CLI 都在用。每个子命令（add、list 等）通过 `rootCmd.AddCommand()` 注册到根命令上。

---

### 3. `internal/model/task.go` — 数据模型

```go
type Status string

const (
    Todo Status = "todo"
    Done Status = "done"
)

type Task struct {
    ID          int       `json:"id"`
    Description string    `json:"description"`
    Priority    string    `json:"priority"`
    CreatedAt   time.Time `json:"created_at"`
    DoneAt      time.Time `json:"done_at,omitempty"`
    Status      Status    `json:"status"`
}
```

**作用：** 定义核心数据结构。

**字段说明：**

| 字段          | 类型        | 说明                                   |
| ------------- | ----------- | -------------------------------------- |
| `ID`          | `int`       | 任务唯一标识，自增                     |
| `Description` | `string`    | 任务描述文本                           |
| `Priority`    | `string`    | 优先级：`high` / `medium` / `low`      |
| `CreatedAt`   | `time.Time` | 创建时间                               |
| `DoneAt`      | `time.Time` | 完成时间（`omitempty` 未完成时不序列化）|
| `Status`      | `Status`    | 状态：`todo` 或 `done`                 |

**知识点：**

- `json:"..."` 是 Go 的 struct tag，控制 JSON 序列化时的字段名。
- `omitempty` 表示零值时省略该字段，所以未完成的任务不会有 `done_at` 字段。
- `Status` 用自定义类型（而非裸 `string`），可以利用常量约束取值范围，代码更安全。

---

### 4. `internal/storage/json_storage.go` — 持久化层

```go
type JSONStorage struct {
    path string
    mu   sync.Mutex
}
```

**核心方法：**

#### `NewJSONStorage()` — 构造函数

```go
func NewJSONStorage() (*JSONStorage, error) {
    home, _ := os.UserHomeDir()
    path := filepath.Join(home, ".gtodo", "tasks.json")
    os.MkdirAll(filepath.Dir(path), 0755)
    return &JSONStorage{path: path}, nil
}
```

自动在用户主目录下创建 `~/.gtodo/` 目录，数据文件为 `tasks.json`。

#### `Load()` — 读取全部任务

```go
func (s *JSONStorage) Load() ([]model.Task, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    // 读取文件 → JSON 反序列化 → 按 ID 排序 → 返回
}
```

- 文件不存在时返回空切片（而非报错），这样首次使用不需要手动创建文件。
- 读取后按 ID 排序，保证输出和 NextID 的正确性。

#### `Save()` — 保存全部任务

```go
func (s *JSONStorage) Save(tasks []model.Task) error {
    // JSON 序列化（带缩进） → 写入文件
}
```

使用 `json.MarshalIndent` 生成格式化 JSON，方便人工查看。

#### `NextID()` — 生成下一个 ID

```go
func (s *JSONStorage) NextID(tasks []model.Task) int {
    if len(tasks) == 0 {
        return 1
    }
    return tasks[len(tasks)-1].ID + 1
}
```

取已排序列表最后一个 ID + 1，简单有效。

**知识点：**

- `sync.Mutex` 互斥锁，防止并发读写冲突（CLI 场景实际很少并发，但这是良好习惯）。
- `os.UserHomeDir()` 获取用户主目录，跨平台兼容。
- `filepath.Join()` 跨平台拼接路径（Windows 用 `\`，Linux 用 `/`）。

---

### 5. `cmd/add.go` — 添加事项

```go
var addCmd = &cobra.Command{
    Use:   "add [description]",
    Short: "添加一个新待办事项",
    Args:  cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        // 1. 初始化存储
        // 2. 加载现有任务
        // 3. 构造新 Task（NextID 自增、当前时间、默认状态 todo）
        // 4. 追加到列表并保存
    },
}
```

**流程：** 初始化存储 → 加载 → 构造新任务 → 追加 → 保存。

**知识点：**

- `cobra.MinimumNArgs(1)` 要求至少传 1 个参数，否则 Cobra 自动报错。
- `-p` flag 在 `init()` 中通过 `StringVarP` 绑定到 `priority` 变量，P 后缀表示支持短选项。

---

### 6. `cmd/list.go` — 列出事项

```go
Run: func(cmd *cobra.Command, args []string) {
    // 1. 加载任务
    // 2. 按 --all flag 过滤（默认只显示未完成）
    // 3. 用 tablewriter 渲染彩色表格
}
```

**流程：** 加载 → 过滤 → 表格渲染。

**知识点：**

- `tablewriter` v1.1.3 使用新 API：`Header()` 设置表头，`Append()` 逐行添加，`Render()` 输出。
- `fatih/color` 给优先级和状态加颜色：红色=高、黄色=中、绿色=低。
- `--all` 是 bool flag，通过 `BoolVar` 绑定。

---

### 7. `cmd/done.go` — 完成事项

```go
Run: func(cmd *cobra.Command, args []string) {
    // 1. 解析 ID 参数
    // 2. 加载任务
    // 3. 遍历找到对应任务，修改 Status 和 DoneAt
    // 4. 保存
}
```

**流程：** 解析 ID → 加载 → 查找并修改 → 保存。

**注意：** 用 `for i := range tasks` 而非 `for _, t := range tasks`，因为需要修改切片中的元素。`range` 的值拷贝不会反映回原切片，用索引才能原地修改。

---

### 8. `cmd/delete.go` — 删除事项

```go
Run: func(cmd *cobra.Command, args []string) {
    // 1. 解析 ID 参数
    // 2. 加载任务
    // 3. 构建新切片，跳过目标 ID
    // 4. 保存新切片
}
```

**流程：** 解析 ID → 加载 → 过滤掉目标 → 保存。

**知识点：** 删除采用"重建切片"模式 —— 遍历时跳过要删除的元素，将其余元素追加到新切片。这比直接操作索引删除更安全、更易读。

---

## 依赖库说明

| 库                                   | 作用             |
| ------------------------------------ | ---------------- |
| `github.com/spf13/cobra`            | CLI 命令框架     |
| `github.com/olekukonko/tablewriter` | 终端表格渲染     |
| `github.com/fatih/color`            | 终端彩色文字输出 |

---

## 数据流示意

以 `gtodo add "写周报" -p high` 为例：

```
用户输入
   │
   ▼
main.go → cmd.Execute()
   │
   ▼
Cobra 路由 → addCmd.Run()
   │
   ├─ storage.NewJSONStorage()     // 初始化存储路径
   ├─ storage.Load()               // 读取 ~/.gtodo/tasks.json
   ├─ storage.NextID(tasks)        // 计算新 ID
   ├─ 构造 model.Task{...}        // 组装任务对象
   ├─ tasks = append(tasks, task)  // 追加到列表
   └─ storage.Save(tasks)          // 写回 JSON 文件
   │
   ▼
输出：已添加任务 #1: 写周报 (优先级: high)
```

---

## 扩展方向

如果你想继续练习，以下是一些扩展思路：

1. **添加 `edit` 命令** —— 修改已有事项的描述或优先级
2. **支持截止日期** —— 在 Task 中加 `Deadline` 字段，list 时高亮过期事项
3. **换用 SQLite** —— 替换 JSONStorage 为数据库存储，学习 `database/sql`
4. **添加单元测试** —— 为 storage 和 model 写 `_test.go`
5. **支持配置文件** —— 用 `viper`（Cobra 配套库）管理配置
