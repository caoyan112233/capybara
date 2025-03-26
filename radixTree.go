package capybara

import (
	"strings"
)

type RadixNode struct {
	path     string // "/abcd"
	children map[string]*RadixNode
	isEnd    bool
	// paramChild *RadixNode
	// anyChild   *RadixNode
}

type RadixTree struct {
	root *RadixNode
}

func (n *RadixNode) findChild(prefix string) *RadixNode {
	for childPrefix, child := range n.children {
		if strings.HasPrefix(prefix, childPrefix) {
			return child
		}
	}
	return nil
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (n *RadixNode) insert2(path string) {
	// 首次插入一个新节点
	if len(n.children) == 0 && n.path == "" {
		n.path = path
	} else {
		// 先求出公共前缀
		currentNode := n
		nowPathLen := len(path)
		savePathLen := len(currentNode.path)
		lcpLen := 0
		minLen := min(nowPathLen, savePathLen)

		for i := 0; i < minLen; i++ {
			if path[i] == currentNode.path[i] {
				lcpLen++
			} else {
				break
			}
		}
		// 然后执行节点插入逻辑
		// 情况1：  如果最长公共前缀的长度仍然小于原树中的路径长度，将原路径拆开
		if lcpLen < savePathLen {
			// 拆开原路径的非前缀部分
			newPath := currentNode.path[lcpLen:]
			newNode := &RadixNode{
				path:     newPath,
				children: make(map[string]*RadixNode),
			}
			// 修改原路径的节点信息
			currentNode.path = path[0:lcpLen]
			// 如果原路径有子节点，需要将子节点重新分配到i父节点中
			for key, child := range currentNode.children {
				delete(currentNode.children, key)
				newNode.children[key] = child
			}
			// 移除原路径节点的子节点
			// 将拆出的部分插入到前缀节点中
			currentNode.children[newPath] = newNode

			// 添加新节点
			newPath2 := path[lcpLen:]
			newNode2 := &RadixNode{
				path:     newPath2,
				children: make(map[string]*RadixNode),
			}
			currentNode.children[newPath2] = newNode2
		} else if lcpLen > savePathLen {
			// 将新来的路径拆开
			newPath := path[lcpLen:]
			newNode := &RadixNode{
				path:     newPath,
				children: make(map[string]*RadixNode),
			}
			currentNode.children[newPath] = newNode
		} else if lcpLen == savePathLen {
			// 原路径不用改变
			newPath := path[lcpLen:]
			if _, ok := currentNode.children[newPath]; ok {

			}
		}
	}
}
func (n *RadixNode) insert(path string) {
	current := n
	for {
		//  计算最长公共前缀
		lcpLen := 0
		max := min(len(current.path), len(path))
		for ; lcpLen < max && current.path[lcpLen] == path[lcpLen]; lcpLen++ {
		}
		// 场景1：完全匹配当前节点路径
		if lcpLen == len(current.path) {
			path = path[lcpLen:]
			if len(path) == 0 {
				current.isEnd = true
				return
			}
			// 查找或创建子节点
			child := current.findChild(path)
			if child == nil {
				child = &RadixNode{path: path, children: make(map[string]*RadixNode)}
				current.children[path] = child
			}
			current = child
		} else {
			// 场景2：拆分当前节点（路径压缩）
			commonPrefix := current.path[:lcpLen]
			remaining := current.path[lcpLen:]

			// 创建新父节点并继承原节点属性
			newParent := &RadixNode{
				path:     commonPrefix,
				children: map[string]*RadixNode{remaining: current},
			}
			*current = *newParent // 替换原节点

			// 插入新路径剩余部分
			newChild := &RadixNode{path: path[lcpLen:], children: make(map[string]*RadixNode)}
			current.children[newChild.path] = newChild
			current = newChild
		}
	}

}
