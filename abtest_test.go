package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateABTestList(t *testing.T) {
	assert := assert.New(t)

	type expectedStruct struct {
		abTestBucketList ABTestBucketList
		err              error
	}

	tests := []struct {
		name     string
		input    map[string]float64
		expected expectedStruct
	}{
		{
			name: "参数不超过100%的",
			input: map[string]float64{
				"A": 0.4, "B": 0.3,
			},
			expected: expectedStruct{
				abTestBucketList: nil,
				err:              errors.New("分配的比例异常"),
			},
		},
		{
			name: "参数超过100%的",
			input: map[string]float64{
				"A": 0.9, "B": 0.3,
			},
			expected: expectedStruct{
				abTestBucketList: nil,
				err:              errors.New("分配的比例异常"),
			},
		},
		{
			name: "只有两个的",
			input: map[string]float64{
				"A": 0.7, "B": 0.3,
			},
			expected: expectedStruct{
				abTestBucketList: []ABTestBucket{
					{Name: "A", Start: 1, End: 700},
					{Name: "B", Start: 701, End: 1000},
				},
				err: nil,
			},
		},
		{
			name: "多个的",
			input: map[string]float64{
				"A": 0.2, "B": 0.3, "C": 0.1, "D": 0.4,
			},
			expected: expectedStruct{
				abTestBucketList: []ABTestBucket{
					{Name: "A", Start: 1, End: 200},
					{Name: "B", Start: 201, End: 500},
					{Name: "C", Start: 501, End: 600},
					{Name: "D", Start: 601, End: 1000},
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := CreateABTestList(test.input)
			if test.expected.err != nil {
				assert.EqualError(err, test.expected.err.Error())
			}
			assert.Equal(ret, test.expected.abTestBucketList)
		})

	}

}

func TestHashBucket(t *testing.T) {

	assert := assert.New(t)

	tests := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "正常初始化",
			input:    "asdfasdfasdf",
			expected: "B",
			err:      nil,
		},
	}

	abTestConfig := map[string]float64{
		"A": 0.3, "B": 0.7,
	}
	abTestBucketList, err := CreateABTestList(abTestConfig)
	if err != nil {
		assert.Nil(err)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := abTestBucketList.HashBucket(test.input)
			assert.Nil(err)
			assert.Equal(ret, test.expected)
		})
	}

}
