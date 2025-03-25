package capybara

import (
	"strings"
	"testing"
)

func TestNodeInsertAndFind(t *testing.T) {
	root := InitNode()
	testHandler := func(c Context) {}
	// 测试基础路由
	root.insertRoute("/user", "GET", testHandler)
	if n, _ := root.FindRoute("/user"); n == nil {
		t.Error("基础路由查找失败")
	}

	// 测试参数路由
	root.insertRoute("/user/:id", "GET", testHandler)
	if _, params := root.FindRoute("/user/123"); params["id"] != "123" {
		t.Error("参数路由解析失败")
	}

	// 测试通配符路由
	root.insertRoute("/static/*filepath", "GET", testHandler)
	_, params := root.FindRoute("/static/css/style.css")
	if params["filepath"] != "css/style.css" {
		t.Error("通配符路由解析失败")
	}
}

// 测试路由冲突
func TestRouteConflict(t *testing.T) {
	root := InitNode()
	testHandler := func(c Context) {}
	root.insertRoute("/user/delete", "GET", testHandler)
	root.insertRoute("/user/:action", "POST", testHandler)
	if n, _ := root.FindRoute("/user/delete"); n == nil {
		t.Error("静态路由被参数路由覆盖")
	}
}

// 测试方法覆盖
func TestMethodOverride(t *testing.T) {
	root := InitNode()
	getHandler := func(c Context) {}
	postHandler := func(c Context) {}

	root.insertRoute("/login", "GET", getHandler)
	root.insertRoute("/login", "POST", postHandler)

	// 验证方法存储
	if n, _ := root.FindRoute("/login"); n.method != "POST" {
		t.Error("方法覆盖异常")
	}
}

// 测试多层嵌套路由
func TestNestedRoutes(t *testing.T) {
	root := InitNode()
	testHandler := func(c Context) {}

	routes := []string{
		"/api/v1/users",
		"/api/v1/users/:id",
		"/api/v1/products",
	}

	for _, route := range routes {
		root.insertRoute(route, "GET", testHandler)
	}

	for _, route := range routes {
		n, _ := root.FindRoute(route)
		if n == nil {
			t.Errorf("嵌套路由 %s 查找失败", route)
		}
	}
}

// 测试通配符优先级
func TestWildcardPriority(t *testing.T) {
	root := InitNode()
	testHandler := func(c Context) {}

	root.insertRoute("/:version/user", "GET", testHandler)
	root.insertRoute("/v1/*catchall", "GET", testHandler)

	// 验证精确匹配优先
	if n, _ := root.FindRoute("/v1/user"); n == nil {
		t.Error("精确匹配优先级异常")
	}
}

// 测试路径分割
func TestPathSplitting(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"/user/profile", []string{"user", "profile"}},
		{"/api/v1//data/", []string{"api", "v1", "data"}},
		{"//debug/pprof/", []string{"debug", "pprof"}},
	}

	for _, tc := range testCases {
		result := splitPath(tc.input)
		if strings.Join(result, ",") != strings.Join(tc.expected, ",") {
			t.Errorf("路径分割错误: 输入 %s 期望 %v 得到 %v", tc.input, tc.expected, result)
		}
	}
}

// 测试节点初始化
func TestNodeInitialization(t *testing.T) {
	n := InitNode()
	if n.childrens == nil {
		t.Error("子节点映射初始化失败")
	}
	if n.method != "" {
		t.Error("节点方法初始化异常")
	}
}

// 测试错误处理
func TestErrorHandling(t *testing.T) {
	root := InitNode()

	// 测试不存在的路由
	if n, _ := root.FindRoute("/not/exist"); n != nil {
		t.Error("不存在路由错误处理异常")
	}

	// 测试非法路径
	root.insertRoute("invalid_path", "GET", func(c Context) {})
	if n, _ := root.FindRoute("invalid_path"); n != nil {
		t.Error("非法路径处理异常")
	}
}

// 测试参数覆盖
func TestParamOverride(t *testing.T) {
	root := InitNode()
	testHandler := func(c Context) {}

	root.insertRoute("/:category/:id", "GET", testHandler)
	_, params := root.FindRoute("/books/123")
	if params["category"] != "books" || params["id"] != "123" {
		t.Error("多参数解析失败")
	}
}
