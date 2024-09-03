//go:build !solution

package treeiter

func DoInOrder[T interface {
	Left() *T
	Right() *T
}](tree *T, f func(*T)) {
	if tree == nil {
		return
	}

	DoInOrder((*tree).Left(), f)
	f(tree)
	DoInOrder((*tree).Right(), f)
}
