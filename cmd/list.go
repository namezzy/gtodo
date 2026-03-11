// cmd/list.go   ← v1.1.3 兼容写法（简化，不要太多高级设定）
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/namezzy/gtodo/internal/model"
	"github.com/namezzy/gtodo/internal/storage"
)

var showAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出待办事项",
	Run: func(cmd *cobra.Command, args []string) {
		st, err := storage.NewStorage()
		if err != nil {
			fmt.Fprintln(os.Stderr, "无法初始化存储:", err)
			os.Exit(1)
		}
		defer st.Close()

		tasks, err := st.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "读取任务失败:", err)
			os.Exit(1)
		}

		var displayTasks []model.Task
		if showAll {
			displayTasks = tasks
		} else {
			for _, t := range tasks {
				if t.Status == model.Todo {
					displayTasks = append(displayTasks, t)
				}
			}
		}

		if len(displayTasks) == 0 {
			msg := "目前没有待办事项"
			if showAll {
				msg = "目前没有任何任务"
			}
			fmt.Println(msg)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header("ID", "优先级", "事项", "创建时间", "状态")

		// 颜色与内容
		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		for _, t := range displayTasks {
			pri := t.Priority
			switch pri {
			case "high":
				pri = red("高")
			case "medium":
				pri = yellow("中")
			case "low":
				pri = green("低")
			default:
				pri = t.Priority
			}

			status := "待办"
			if t.Status == model.Done {
				status = green("已完成")
			}

			table.Append(
				fmt.Sprintf("%d", t.ID),
				pri,
				t.Description,
				t.CreatedAt.Format("2006-01-02 15:04"),
				status,
			)
		}

		if err := table.Render(); err != nil {
			fmt.Fprintln(os.Stderr, "渲染表格失败:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&showAll, "all", false, "显示所有事项（包含已完成）")
}
