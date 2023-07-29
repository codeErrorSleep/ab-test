# ab-test
A simple ab test library

## how to use

通过配置map初始化对应的比例配置

```go
// 初始化配置
abTestConfig := map[string]float64{
  "A": 0.3, "B": 0.7,
}

// 创建ab test
abTestBucketList, err := CreateABTestList(abTestConfig)
if err != nil {
  fmt.Println(err)
}

// 通过这种方式就可以hash了
fmt.Println(abTestBucketList.HashBucket("iii"))


```
