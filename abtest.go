package main

import (
	"errors"
)

// 创建配置
// map A:0.2 B:0.3

// 生成 1-1000
// 可以生成list A:[1:200] B:[201:500]

const (
	// 分配的坑位
	Position = 1000
)

type ABTestBucket struct {
	Name  string
	Start int
	End   int
}

func CreateABTestList(abTestConfigMap map[string]float64) (abTestBucketList []ABTestBucket, err error) {
	// 检验百分比是否到100了
	sum := 0
	for _, v := range abTestConfigMap {
		sum += int(v * 100)
	}
	if sum != 100 {
		return nil, errors.New("分配的比例异常")
	}

	// 计算每个数据的对应桶的范围
	start := 1
	for name, percentage := range abTestConfigMap {
		end := start + int(percentage*Position) - 1
		abTestBucketList = append(abTestBucketList, ABTestBucket{
			Name:  name,
			Start: start,
			End:   end,
		})
		start = end + 1
	}

	return abTestBucketList, nil
}

// hash 获取
// hash到对应的饭位置里面,这里直接遍历就OK了
