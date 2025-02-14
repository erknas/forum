package conv

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedID  int
		expectedErr error
	}{
		{
			name:        "Valid ID",
			input:       "123",
			expectedID:  123,
			expectedErr: nil,
		},
		{
			name:        "Invalid string ID",
			input:       "abc",
			expectedID:  -1,
			expectedErr: errors.New("invalid ID abc"),
		},
		{
			name:        "Zero ID",
			input:       "0",
			expectedID:  0,
			expectedErr: errors.New("invalid ID 0"),
		},
		{
			name:        "Negative ID",
			input:       "-1",
			expectedID:  -1,
			expectedErr: errors.New("invalid ID -1"),
		},
		{
			name:        "Empty String",
			input:       "",
			expectedID:  -1,
			expectedErr: errors.New("invalid ID "),
		},
		{
			name:        "Whitespace",
			input:       " 123 ",
			expectedID:  -1,
			expectedErr: errors.New("invalid ID  123 "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ID(tt.input)

			assert.Equal(t, tt.expectedID, id)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
