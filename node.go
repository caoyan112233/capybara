package capybara

import (
	"strings"
)

type node struct {
	name      string           // 当前结点的名字
	childrens map[string]*node // 当前结点的子结点
	isWild    bool             // 是否为通配符结点
	handler   HandlerFunc      // 当前路由的路由函数
	method    string           // 当前路由的请求方法
	fullPath  string           // 当前结点的完整路由
}

// 插入路径
func (n *node) insertRoute(path string, method string, handler HandlerFunc) {
	// 如果路由第一个字符不是/，我们就添加一个
	if path[0] != '/' {
		path = "/" + path
	}
	// 如果是空的路由，默认设置成 /
	if path == "" {
		path = "/"
	}

	segments := splitPath(path)
	// 将目前结点作为一个根节点
	currNode := n
	
	for i := 0; i < len(segments); i++ {
		if _, exists := currNode.childrens[segments[i]]; !exists {
			currNode.childrens[segments[i]] = &node{
				name:      segments[i],
				childrens: make(map[string]*node),
				isWild:    strings.HasPrefix(segments[i], ":") || segments[i][0] == '*',
			}

		}
		currNode = currNode.childrens[segments[i]]
	}
	currNode.handler = handler
	currNode.method = method
	currNode.fullPath = path
}

// 找结点路径
func (n *node) FindRoute(path string) (*node, map[string]string) {
	if path == "" {
		return nil, nil
	}
	segments := splitPath(path)
	currNode := n
	params := make(map[string]string)
	for i, seg := range segments {
		if child, exists := currNode.childrens[seg]; exists {
			currNode = child
			continue
		}
		for key, child := range currNode.childrens {
			if child.isWild {
				if strings.HasPrefix(key, ":") {
					params[key[1:]] = seg
				} else if strings.HasPrefix(key, "*") {
					// 遇到了通配符，捕获剩余的路径
					longPath := strings.Join(segments[i:], "/")
					params[key[1:]] = longPath
				}
				currNode = child
			}
		}
	}
	if currNode.handler == nil {
		return nil, nil
	}
	return currNode, params
}

// 初始化单个结点
func InitNode() *node {
	return &node{
		name:      "",
		method:    "",
		childrens: make(map[string]*node),
		isWild:    false,
	}
}
