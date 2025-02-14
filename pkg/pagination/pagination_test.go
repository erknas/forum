package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		page           *int32
		pageSize       *int32
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "Valid Inputs",
			page:           int32Ptr(2),
			pageSize:       int32Ptr(20),
			expectedLimit:  20,
			expectedOffset: 20,
		},
		{
			name:           "Nil Page",
			page:           nil,
			pageSize:       int32Ptr(15),
			expectedLimit:  15,
			expectedOffset: 0,
		},
		{
			name:           "Nil PageSize",
			page:           int32Ptr(3),
			pageSize:       nil,
			expectedLimit:  10,
			expectedOffset: 20,
		},
		{
			name:           "Negative Page",
			page:           int32Ptr(-1),
			pageSize:       int32Ptr(10),
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "Negative PageSize",
			page:           int32Ptr(2),
			pageSize:       int32Ptr(-5),
			expectedLimit:  10,
			expectedOffset: 10,
		},
		{
			name:           "Both Nil",
			page:           nil,
			pageSize:       nil,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "Zero Page",
			page:           int32Ptr(0),
			pageSize:       int32Ptr(10),
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "Zero PageSize",
			page:           int32Ptr(2),
			pageSize:       int32Ptr(0),
			expectedLimit:  10,
			expectedOffset: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, offset := New(tt.page, tt.pageSize)

			assert.Equal(t, tt.expectedLimit, limit)
			assert.Equal(t, tt.expectedOffset, offset)
		})
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}
