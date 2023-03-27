package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	} else {
		Walk(t.Left, ch)
		ch <- t.Value
		Walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		Walk(t1, ch1)
		close(ch1)
	}()
	go func() {
		Walk(t2, ch2)
		close(ch2)
	}()

	for v1 := range ch1 {
		if v2, ok := <-ch2; ok && v2 != v1 {
			return false
		}
	}
	return true
}

func main() {
	if !Same(tree.New(1), tree.New(1)) || Same(tree.New(1), tree.New(2)) {
		panic(fmt.Errorf("wrong walk implemation"))
	}
	fmt.Println("\nSuccessful implemetation")
}
