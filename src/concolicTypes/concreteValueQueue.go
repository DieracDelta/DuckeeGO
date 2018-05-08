package concolicTypes

type ConcreteValueQueue struct {
	queue []*ConcreteValues
}

func initialConcreteValueQueue() *ConcreteValueQueue {
	return &ConcreteValueQueue{[]*ConcreteValues{newConcreteValues()}}
}

func (q *ConcreteValueQueue) isEmpty() bool {
	return len(q.queue) == 0
}

func (q *ConcreteValueQueue) enqueue(cv *ConcreteValues) {
	q.queue = append(q.queue, cv)
}

func (q *ConcreteValueQueue) dequeue() *ConcreteValues {
	res := q.queue[0]
	q.queue = q.queue[1:]
	return res
}