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
	Percentage int
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

func CreateABTestList(abTestConfigMap map[string]int) (abTestBucketList ABTestBucketList, err error) {
	if len(abTestConfigMap) == 0 {
		return nil, errors.New("Input cannot be empty")
	}

	// Check if the percentages add up to 100
	sum := 0
	for _, v := range abTestConfigMap {
		if v < 0 {
			return nil, errors.New("Input percentage cannot be negative")
		}

		sum += v
	}
	if sum != 100 {
		return nil, errors.New("The sum of input percentages must be 100")
	}

	confList := ConfList{}

	for name, percentage := range abTestConfigMap {
		confList = append(confList, Conf{Name: name, Percentage: percentage})
	}
	// Sort to ensure the same bucket is allocated each time
	sort.Sort(confList)

	// Calculate the range of each data's corresponding bucket
	start := 1
	for i := range confList {
		// If it's 0, no need to allocate
		if confList[i].Percentage == 0 {
			continue
		}

		end := start + int(confList[i].Percentage*10) - 1
		abTestBucketList = append(abTestBucketList, ABTestBucket{
			Name:  confList[i].Name,
			Start: start,
			End:   end,
		})
		start = end + 1
	}

	return abTestBucketList, nil
}

// Get hash
// Hash to the corresponding position, here you can directly traverse
func (abtestBucketList ABTestBucketList) HashBucket(value string) (result string, err error) {
	if len(abtestBucketList) == 0 {
		return "", errors.New("Not initialized")
	}

	h := fnv.New32a()
	h.Write([]byte(value))
	hashValue := h.Sum32()

	// Map the hash value to the range from 1 to maxRange
	hashPosition := int(hashValue%uint32(Position)) + 1

	for i := range abtestBucketList {
		if hashPosition <= abtestBucketList[i].End && hashPosition >= abtestBucketList[i].Start {
			return abtestBucketList[i].Name, nil
		}
	}
	return abtestBucketList[0].Name, nil
}

func main() {
	// 初始化配置
	abTestConfig := map[string]int{
		"A": 30, "B": 70,
	}

	// 创建ab test
	abTestBucketList, err := CreateABTestList(abTestConfig)
	if err != nil {
		fmt.Println(err)
	}

	// 通过这种方式就可以hash了
	result, err := abTestBucketList.HashBucket("iii")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

}
