package capybara

import "strings"

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
	if path == "" {
		return []string{}
	}
	return strings.Split(strings.TrimPrefix(path, "/"), "/")
}
