package main

import "fmt"

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
	result, err := abTestBucketList.HashBucket("iii", DefaultHashFunc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

}
