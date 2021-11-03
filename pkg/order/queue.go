package order

import "encoding/json"

type Queue struct {
	Orders []Order
}

func (q *Queue) Set(o Order) {
	length := len(q.Orders)
	if length > 12 {
		q.Orders = append(q.Orders[1:], o)
	}

	q.Orders = append(q.Orders, o)
}

func (q *Queue) Get() ([]byte, error) {
	j, err := json.Marshal(q.Orders)
	if err != nil {
		return j, ErrParseFailer
	}
	return j, nil
}

func NewQueue() *Queue {
	return &Queue{Orders: make([]Order, 0)}
}
