package search

type queue []string

func newQueue() queue {
	queue := make([]string, 0, 10)
	return queue
}

func (q *queue) push(el string) {
	*q = append(*q, el)
}

func (q *queue) pop() string {
	el := (*q)[0]
	*q = (*q)[1:]
	return el
}

func (q *queue) len() int {
	return len(*q)
}