package priorityqueue

import (
	"log"
	"slices"
)

type PriorityQueue[T comparable] struct {
	values []value[T]
}

type value[D comparable] struct {
	data     D
	priority int
}

func New[T comparable]() (p PriorityQueue[T]) {
	return
}

func (p *PriorityQueue[T]) sort() {
	slices.SortStableFunc(p.values, func(a value[T], b value[T]) int {
		return a.priority - b.priority
	})
}

func (p *PriorityQueue[T]) Set(data T, priority int) {
	var values []value[T]
	found := false
	for _, value := range p.values {
		if value.data == data {
			value.priority = priority
			found = true
		}
		values = append(values, value)
	}
	p.values = values
	if !found {
		newValue := value[T]{
			data:     data,
			priority: priority,
		}
		p.values = append(p.values, newValue)
	}
	p.sort()
}

func (p *PriorityQueue[T]) Has(data T) bool {
	for _, value := range p.values {
		if value.data == data {
			return true
		}
	}
	return false
}

func (p *PriorityQueue[T]) Values() (values []T) {
	for _, value := range p.values {
		values = append(values, value.data)
	}
	return
}

// pop the highest priority value
func (p *PriorityQueue[T]) Pop() T {
	if p.Empty() {
		log.Fatalf("priority queue is empty")
	}
	first := p.values[0].data
	var values []value[T]
	for i, value := range p.values {
		if i != 0 {
			values = append(values, value)
		}
	}
	p.values = values
	return first
}

func (p *PriorityQueue[T]) Remove(data T) {
	if p.Empty() {
		log.Fatalf("priority queue is empty")
	}
	var values []value[T]
	for _, value := range p.values {
		if value.data != data {
			values = append(values, value)
		}
	}
	p.values = values
}

func (p *PriorityQueue[T]) Empty() bool {
	return len(p.values) == 0
}

func (p *PriorityQueue[T]) Size() int {
	return len(p.values)
}

func (p *PriorityQueue[T]) Get(data T) (d T, priority int, ok bool) {
	for _, value := range p.values {
		if value.data == data {
			d = value.data
			priority = value.priority
			ok = true
			return
		}
	}
	return
}
