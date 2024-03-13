package logic

import (
	"fmt"
)

func (logic *Logic) GetFullEndpoint(ip string, port string, endpoint string, startIdx int, endIdx int) string {
	return fmt.Sprintf("http://%s:%s/%s/%d-%d", ip, port, endpoint, startIdx, endIdx)
}

func (logic *Logic) HandleResponse(response map[string]interface{}, key string) (interface{}, error) {
	result, ok := response[key]
	if !ok {
		return fmt.Errorf("key '%s' is not found in response", key), nil
	}

	return result, nil
}

type Queue []string

//IsEmpty - check if queue is empty
func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

//Enqueue - append value to the queue
func (q *Queue) Enqueue (id string) {
	*q = append(*q, id)
	fmt.Printf("Enqueue: %v\n", id)
}

//Dequeue - pop first element from queue
func (q *Queue) Dequeue() (string, error) {
	if q.IsEmpty() {
		return "", fmt.Errorf("queue is empty")
	}
	data := (*q)[0] // get first element
	*q = (*q)[1:]   // remove first element
	fmt.Printf("Dequeue: %v\n", data)
	return data, nil
}