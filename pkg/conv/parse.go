package conv

import (
	"fmt"
	"strconv"
)

func ID(stdID string) (int, error) {
	id, err := strconv.Atoi(stdID)
	if err != nil {
		return -1, fmt.Errorf("invalid ID")
	}

	if id <= 0 {
		return -1, fmt.Errorf("invalid ID")
	}

	return id, nil
}
