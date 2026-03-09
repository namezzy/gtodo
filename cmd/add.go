package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/namezzy/gtodo/internal/model"
	"github.com/namezzy/gtodo/internal/storage"
	"github.com/spf13/cobra"
)

var priority string

var addCmd = &cobra.Command{
	Use:   "add [description]",
	Short: "添加一个新待办事项",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		desc := args[0] // 简单起见只取第一个参数，实际可 strings.Join(args, " ")

		sto, err := storage.NewJSONStorage()
		if err != nil {
			fmt.Fprintln(os.Stderr, "存储初始化失败:", err)
			os.Exit(1)
		}

		tasks, err := sto.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "读取任务失败:", err)
			os.Exit(1)
		}

		task := model.Task{
			ID:          sto.NextID(tasks),
			Description: desc,
			Priority:    priority,
			CreatedAt:   time.Now(),
			Status:      model.Todo,
		}

		tasks = append(tasks, task)

		if err := sto.Save(tasks); err != nil {
			fmt.Fprintln(os.Stderr, "保存失败:", err)
			os.Exit(1)
		}

		fmt.Printf("已添加任务 #%d: %s (优先级: %s)\n", task.ID, task.Description, task.Priority)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&priority, "priority", "p", "medium", "优先级: high | medium | low")
}
