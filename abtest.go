package main

import (
	"errors"
	"fmt"
	"hash/fnv"
	"sort"
)

// 创建配置
// map A:0.2 B:0.3

// 生成 1-1000
// 可以生成list A:[1:200] B:[201:500]

const (
	// 分配的坑位
	Position = 1000
)

// 初始化的配置
type Conf struct {
	Name       string
	Percentage float64
}

type ConfList []Conf

func (c ConfList) Len() int           { return len(c) }
func (c ConfList) Less(i, j int) bool { return c[i].Name < c[j].Name }
func (c ConfList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

// 分配的桶
type ABTestBucket struct {
	Name  string
	Start int
	End   int
}

type ABTestBucketList []ABTestBucket

func CreateABTestList(abTestConfigMap map[string]float64) (abTestBucketList ABTestBucketList, err error) {
	// 检验百分比是否到100了
	sum := 0
	for _, v := range abTestConfigMap {
		sum += int(v * 100)
	}
	if sum != 100 {
		return nil, errors.New("分配的比例异常")
	}

	confList := ConfList{}

	for name, percentage := range abTestConfigMap {
		confList = append(confList, Conf{Name: name, Percentage: percentage})
	}
	// 排序保证每次分配到的都是同一个桶
	sort.Sort(confList)

	// 计算每个数据的对应桶的范围
	start := 1
	for i := range confList {
		end := start + int(confList[i].Percentage*Position) - 1
		abTestBucketList = append(abTestBucketList, ABTestBucket{
			Name:  confList[i].Name,
			Start: start,
			End:   end,
		})
		start = end + 1
	}

	return abTestBucketList, nil
}

// hash 获取
// hash到对应的饭位置里面,这里直接遍历就OK了
func (abtestBucketList ABTestBucketList) HashBucket(value string) (result string, err error) {
	if len(abtestBucketList) == 0 {
		return "", errors.New("未初始化")
	}

	h := fnv.New32a()
	h.Write([]byte(value))
	hashValue := h.Sum32()

	// 将哈希值映射到 1 到 maxRange 的范围内
	hashPosition := int(hashValue%uint32(Position)) + 1

	for i := range abtestBucketList {
		if hashPosition <= abtestBucketList[i].End && hashPosition >= abtestBucketList[i].Start {
			return abtestBucketList[i].Name, nil
		}
	}
	return abtestBucketList[0].Name, nil
}

func main() {

	abTestConfig := map[string]float64{
		"A": 0.3, "B": 0.7,
	}
	abTestBucketList, err := CreateABTestList(abTestConfig)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(abTestBucketList.HashBucket("iii"))

	h := fnv.New32a()
	h.Write([]byte("Bbbb"))
	hashValue := h.Sum32()
	fmt.Println(hashValue)
}
