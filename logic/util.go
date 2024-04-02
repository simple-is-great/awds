package logic

import (
	"fmt"
	"regexp"
	"strconv"
)

func (logic *Logic) GetFullEndpoint(ip string, port string, endpoint string, startIdx int, endIdx int) string {
	return fmt.Sprintf("http://%s:%s/%s/%d-%d", ip, port, endpoint, startIdx, endIdx)
}

// extractMetric parses the metric value from the metrics body using a regular expression.
func extractMetric(metricsBody, metricName string) (float64, error) {
    // Regular expression to match the metric line
    re := regexp.MustCompile(metricName + ` (\d+(\.\d+)?(e[+-]?\d+)?)`)
    matches := re.FindStringSubmatch(metricsBody)

    if len(matches) < 2 {
        return 0, fmt.Errorf("metric %s not found", metricName)
    }

    // Convert the string value to float64
    value, err := strconv.ParseFloat(matches[1], 64)
    if err != nil {
        return 0, fmt.Errorf("error parsing value for metric %s: %v", metricName, err)
    }

    return value, nil
}


type Queue []string

//IsEmpty - check if queue is empty
func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

//Enqueue - append value to the queue
func (q *Queue) Enqueue (id string) {
	*q = append(*q, id)
	// fmt.Printf("Enqueue: %v\n", id)
}

//Dequeue - pop first element from queue
func (q *Queue) Dequeue() (string, error) {
	if q.IsEmpty() {
		return "", fmt.Errorf("queue is empty")
	}
	data := (*q)[0] // get first element
	*q = (*q)[1:]   // remove first element
	// fmt.Printf("Dequeue: %v\n", data)
	return data, nil
}

// func calculateComputeTime(batchSize float64, elapsedTime float64, networkLatency float64) float64 {
// 	return elapsedTime - (batchSize / networkLatency)
// }