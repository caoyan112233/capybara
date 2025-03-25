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
	// /api/v1//data/
	// [api v1 "" data ""]
	if path == "" {
		return []string{}
	}
	segs := strings.Split(path, "/")
	ans := make([]string, 0)
	for _, s := range segs {
		if s != "" {
			ans = append(ans, s)
		}
	}
	return ans
}

func checkPath(path string) bool {
	return strings.HasPrefix(path, "/")
}
