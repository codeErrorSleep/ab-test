package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateABTestList(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		input    map[string]int
		expected ABTestBucketList
		err      error
	}{
		{
			name: "Normal distribution",
			input: map[string]int{
				"A": 30,
				"B": 70,
			},
			expected: ABTestBucketList{
				{
					Name:  "A",
					Start: 1,
					End:   300,
				},
				{
					Name:  "B",
					Start: 301,
					End:   1000,
				},
			},
			err: nil,
		},
		{
			name: "Distribution ratio exception",
			input: map[string]int{
				"A": 30,
				"B": 80,
			},
			expected: nil,
			err:      errors.New("The sum of input percentages must be 100"),
		},
		{
			name: "Extreme ratio",
			input: map[string]int{
				"A": 1,
				"B": 99,
			},
			expected: ABTestBucketList{
				{
					Name:  "A",
					Start: 1,
					End:   10,
				},
				{
					Name:  "B",
					Start: 11,
					End:   1000,
				},
			},
			err: nil,
		},
		{
			name: "Multiple buckets",
			input: map[string]int{
				"A": 25,
				"B": 25,
				"C": 25,
				"D": 25,
			},
			expected: ABTestBucketList{
				{
					Name:  "A",
					Start: 1,
					End:   250,
				},
				{
					Name:  "B",
					Start: 251,
					End:   500,
				},
				{
					Name:  "C",
					Start: 501,
					End:   750,
				},
				{
					Name:  "D",
					Start: 751,
					End:   1000,
				},
			},
			err: nil,
		},
		{
			name: "Boundary value",
			input: map[string]int{
				"A": 0,
				"B": 100,
			},
			expected: ABTestBucketList{
				{
					Name:  "B",
					Start: 1,
					End:   1000,
				},
			},
			err: nil,
		},
		{
			name:     "Input is empty",
			input:    map[string]int{},
			expected: nil,
			err:      errors.New("Input cannot be empty"),
		},
		{
			name: "The sum of input percentages is not 100",
			input: map[string]int{
				"A": 30,
				"B": 60,
			},
			expected: nil,
			err:      errors.New("The sum of input percentages must be 100"),
		},
		{
			name: "Input percentage is negative",
			input: map[string]int{
				"A": -30,
				"B": 130,
			},
			expected: nil,
			err:      errors.New("Input percentage cannot be negative"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := CreateABTestList(test.input)
			if test.err != nil {
				assert.EqualError(err, test.err.Error())
			} else {
				assert.Nil(err)
			}
			assert.Equal(ret, test.expected)
		})
	}
}

func TestHashBucket(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name       string
		bucketList ABTestBucketList
		input      string
		expected   string
		err        error
	}{
		{
			name: "Normal case",
			bucketList: ABTestBucketList{
				{
					Name:  "A",
					Start: 1,
					End:   500,
				},
				{
					Name:  "B",
					Start: 501,
					End:   1000,
				},
			},
			input:    "test",
			expected: "A",
			err:      nil,
		},
		{
			name:       "Bucket list is empty",
			bucketList: ABTestBucketList{},
			input:      "test",
			expected:   "",
			err:        errors.New("Not initialized"),
		},
		{
			name: "Normal case",
			bucketList: ABTestBucketList{
				{
					Name:  "A",
					Start: 1,
					End:   500,
				},
				{
					Name:  "B",
					Start: 501,
					End:   1000,
				},
			},
			input:    "test",
			expected: "A",
			err:      nil,
		},
		{
			name:       "Bucket list is empty",
			bucketList: ABTestBucketList{},
			input:      "test",
			expected:   "",
			err:        errors.New("Not initialized"),
		},
		{
			name: "Single character input string",
			bucketList: ABTestBucketList{
				{
					Name:  "A",
					Start: 1,
					End:   500,
				},
				{
					Name:  "B",
					Start: 501,
					End:   1000,
				},
			},
			input:    "a",
			expected: "A",
			err:      nil,
		},
		{
			name: "Unicode character input string",
			bucketList: ABTestBucketList{
				{
					Name:  "A",
					Start: 1,
					End:   500,
				},
				{
					Name:  "B",
					Start: 501,
					End:   1000,
				},
			},
			input:    "你好",
			expected: "B",
			err:      nil,
		},
		{
			name: "Very long input string",
			bucketList: ABTestBucketList{
				{
					Name:  "A",
					Start: 1,
					End:   500,
				},
				{
					Name:  "B",
					Start: 501,
					End:   1000,
				},
			},
			input:    strings.Repeat("a", 10000),
			expected: "A",
			err:      nil,
		},
		{
			name: "Input string with special characters",
			bucketList: ABTestBucketList{
				{
					Name:  "A",
					Start: 1,
					End:   500,
				},
				{
					Name:  "B",
					Start: 501,
					End:   1000,
				},
			},
			input:    "!@#$%^&*()",
			expected: "A",
			err:      nil,
		},
		{
			name: "Empty input string",
			bucketList: ABTestBucketList{
				{
					Name:  "A",
					Start: 1,
					End:   500,
				},
				{
					Name:  "B",
					Start: 501,
					End:   1000,
				},
			},
			input:    "",
			expected: "A",
			err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := test.bucketList.HashBucket(test.input)
			if test.err != nil {
				assert.EqualError(err, test.err.Error())
			} else {
				assert.Nil(err)
			}
			assert.Equal(ret, test.expected)
		})
	}
}
func BenchmarkCreateABTestList(b *testing.B) {
	abTestConfigMap := map[string]int{
		"A": 30,
		"B": 70,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = CreateABTestList(abTestConfigMap)
	}
}

func BenchmarkHashBucket(b *testing.B) {
	abTestBucketList := ABTestBucketList{
		{
			Name:  "A",
			Start: 1,
			End:   300,
		},
		{
			Name:  "B",
			Start: 301,
			End:   1000,
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = abTestBucketList.HashBucket("test")
	}
}
