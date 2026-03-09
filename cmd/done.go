// cmd/done.go
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/namezzy/gtodo/internal/model"
	"github.com/namezzy/gtodo/internal/storage"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "将事项标记为已完成",
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

		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				if tasks[i].Status == model.Done {
					fmt.Printf("事项 #%d 已经是完成状态\n", id)
					return
				}
				tasks[i].Status = model.Done
				tasks[i].DoneAt = time.Now()
				found = true
				break
			}
		}

		if !found {
			fmt.Fprintf(os.Stderr, "找不到 ID 为 %d 的事项\n", id)
			os.Exit(1)
		}

		if err := storage.Save(tasks); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("已将事项 #%d 标记为完成 ✓\n", id)
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
