package capybara

import (
	"strings"
)

type node struct {
	path      string           // 当前结点的名字
	childrens map[string]*node // 当前结点的子结点
	isWild    bool             // 是否为通配符结点
	handler   HandlerFunc      // 当前路由的路由函数
	method    string           // 当前路由的请求方法
}

// 插入路径
func (n *node) insertRoute(path string, method string, handler HandlerFunc) {
	// 例子： /user/:id/post/:post_id
	// 先判断 user结点是否存在，如果不存在则创建一个user结点
	// currentNode := n
	segments := splitPath(path)
	currNode := n
	for i := 0; i < len(segments); i++ {
		if _, exists := currNode.childrens[segments[i]]; !exists {
			currNode.childrens[segments[i]] = &node{
				path:      segments[i],
				childrens: make(map[string]*node),
				isWild:    strings.HasPrefix(segments[i], ":") || segments[i][0] == '*',
			}
			currNode = currNode.childrens[segments[i]]
		}
	}
	currNode.handler = handler
	currNode.method = method
}

// 找结点路径
func (n *node) FindRoute(path string) (HandlerFunc, map[string]string, string) {
	segments := splitPath(path)
	currNode := n
	params := make(map[string]string)
	for _, seg := range segments {
		if child, exists := currNode.childrens[seg]; exists {
			currNode = child
			continue
		}

		for key, child := range currNode.childrens {
			if child.isWild {
				if strings.HasPrefix(key, ":") {
					params[key[1:]] = seg
				}
				currNode = child
			}
		}
	}
	return currNode.handler, params, currNode.method
}

// 初始化单个结点
func InitNode() *node {
	return &node{
		path:      "",
		method:    "",
		childrens: make(map[string]*node),
		isWild:    false,
	}
}
