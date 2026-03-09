// cmd/delete.go
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/namezzy/gtodo/internal/model"
	"github.com/namezzy/gtodo/internal/storage"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "删除指定事项",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "无效的 ID：%s\n", idStr)
			os.Exit(1)
		}

		storage, err := storage.NewJSONStorage()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		tasks, err := storage.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		newTasks := make([]model.Task, 0, len(tasks))
		found := false

		for _, t := range tasks {
			if t.ID == id {
				found = true
				continue
			}
			newTasks = append(newTasks, t)
		}

		if !found {
			fmt.Fprintf(os.Stderr, "找不到 ID 为 %d 的事项\n", id)
			os.Exit(1)
		}

		if err := storage.Save(newTasks); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("已删除事项 #%d\n", id)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
