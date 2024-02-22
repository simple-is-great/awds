package logic

import (
	"fmt"
	"strconv"
)



func (logic *Logic) GetFullEndpoint(endpoint string, startIdx int, endIdx int) string {
	return fmt.Sprintf("%s/%d-%d", endpoint, startIdx, endIdx)
}

func (logic *Logic) HandleResponse(response map[string]interface{}) (float64, error) {
	responseStr, ok := response["results"].(string)
	if !ok {
		return float64(-1), fmt.Errorf("pod response 'results' is not a string")
	}

	result, err := strconv.ParseFloat(responseStr, 64)
	if err != nil {
		return float64(-1), fmt.Errorf("failed to parse results to float64: %v", err)
	}

	return result, nil
}
