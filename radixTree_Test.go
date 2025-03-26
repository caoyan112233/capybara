package capybara

import "testing"

func TestRadixTreeInsert(t *testing.T) {
	root := &RadixNode{}
	root.insert("/abcd")
	root.insert("/abce")
	root.insert("/aecb")
	root.insert("/aecd")
	t.Error("精确匹配优先级异常")
}
