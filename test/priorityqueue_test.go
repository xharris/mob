package test

import (
	"mob/pkg/priorityqueue"
	"testing"
)

func TestPriorityQueueOrder(t *testing.T) {
	q := priorityqueue.New[string]()
	q.Set("ur mom", 4)
	q.Set("jack", 2)
	q.Set("jill", 1)
	q.Set("bob", 0)
	q.Set("bob", 3) // should replace previous priority

	values := q.Values()
	expected := []string{"jill", "jack", "bob", "ur mom"}

	if len(values) != len(expected) {
		t.Fatalf("incorrect priority queue size, expected=%d, got=%d", len(expected), len(values))
	}

	for v := range values {
		if values[v] != expected[v] {
			t.Fatalf("priority queue wrong order, expected=%v, got=%v", expected, values)
		}
	}
}
