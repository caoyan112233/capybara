package capybara

import (
	"strings"
)

// 辅助函数：合并前缀和路径，处理斜杠问题
func joinPath(prefix, path string) string {
	if prefix == "" {
		return path
	}
	if path == "" {
		return prefix
	}
	return strings.TrimSuffix(prefix, "/") + "/" + strings.TrimPrefix(path, "/")
}

func splitPath(path string) []string {
	ans := make([]string, 0)
	start := -1
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if start != -1 {
				ans = append(ans, path[start:i])
				start = -1
			}
		} else {
			if start == -1 {
				start = i
			}
		}
	}
	if start != -1 {
		ans = append(ans, path[start:])
	}
	return ans
}

func checkPath(path string) bool {
	return strings.HasPrefix(path, "/")
}
