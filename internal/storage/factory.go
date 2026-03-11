package storage

import (
	"fmt"
	"os"
)

// NewStorage 根据环境变量创建对应的存储实例
// GTODO_STORAGE: "json" (默认) 或 "mysql"
// GTODO_MYSQL_DSN: MySQL 连接字符串 (mysql 模式必须)
func NewStorage() (Storage, error) {
	backend := os.Getenv("GTODO_STORAGE")
	if backend == "" {
		backend = "json"
	}

	switch backend {
	case "json":
		return NewJSONStorage()
	case "mysql":
		dsn := os.Getenv("GTODO_MYSQL_DSN")
		if dsn == "" {
			return nil, fmt.Errorf("使用 MySQL 存储时必须设置 GTODO_MYSQL_DSN 环境变量")
		}
		return NewMySQLStorage(dsn)
	default:
		return nil, fmt.Errorf("不支持的存储后端: %s (可选: json, mysql)", backend)
	}
}
