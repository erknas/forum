package conv

import (
	"fmt"
	"strconv"
)

func ID(stdID string) (int, error) {
	id, err := strconv.Atoi(stdID)
	if err != nil {
		return -1, fmt.Errorf("invalid ID %s", stdID)
	}

	if id <= 0 {
		return id, fmt.Errorf("invalid ID %d", id)
	}

	return id, nil
}
