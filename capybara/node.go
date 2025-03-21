package capybara

// **** route tree
type node struct {
	children  map[string]*node
	wildChild *node
	isWild    bool
	paramKey  string
	handler   map[string]HandlerFunc
}

// 增加一条路由
func (n *node) addRoute(path string, method string, handler HandlerFunc) {
	segs := splitPath(path)
	for _, seg := range segs {
		if seg[0] == ':' || seg[0] == '*' {
			n.wildChild = &node{
				children:  make(map[string]*node),
				wildChild: new(node),
				isWild:    true,
				paramKey:  seg[1:],
				handler:   make(map[string]HandlerFunc)}
			n = n.wildChild
		} else {
			if n.children == nil {
				n.children = make(map[string]*node)
			}
			if _, ok := n.children[seg]; !ok {
				n.children[seg] = &node{
					children:  make(map[string]*node),
					wildChild: new(node),
					isWild:    false,
					paramKey:  seg,
					handler:   make(map[string]HandlerFunc)}
			}
			n = n.children[seg]
		}
	}
	n.handler[method] = handler
}

// 匹配一条路由
func (n *node) matchRoute(path string) *node {
	segs := splitPath(path)
	current := n
	for _, seg := range segs {
		if child, ok := current.children[seg]; ok {
			current = child
			continue
		}
		if current.wildChild != nil {
			wildChild := current.wildChild
			if wildChild.isWild {
				current = wildChild
				break
			} else {
				current = wildChild
				continue
			}
		}
		return nil
	}
	return current
}
