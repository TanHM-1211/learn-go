package messagequeue

import (
	"fmt"
	"testing"
)

func TestQueue_Get(t *testing.T) {
	cap := 10
	q := NewQueue(cap)

	for i := 0; i < cap+10 && !q.Full(); i++ {
		node := &Node{data: &Message{fmt.Sprintf("%d", i)}}
		if err := q.Put(node); err != nil {
			t.Errorf("Can not put %v", i)
		}
	}

	node := &Node{data: &Message{fmt.Sprintf("%d", cap)}}
	if err := q.Put(node); err == nil {
		t.Errorf("successfully put data into full queue")
	}

	for i := 0; q.Empty(); i++ {
		node, err := q.Get()
		if err != nil || node.data.Content != fmt.Sprintf("%d", i) {
			t.Errorf("Error while getting: got %v, want %v", i, i)
		}
	}

}
