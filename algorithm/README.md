

<!-- toc -->

- [algorithm 算法模块](#algorithm-%E7%AE%97%E6%B3%95%E6%A8%A1%E5%9D%97)
  * [crc16_test.go crc16算法](#crc16_testgo-crc16%E7%AE%97%E6%B3%95)
    + [TestCrc16](#testcrc16)
    + [TestCrc16s](#testcrc16s)
  * [descartes_test.go 笛卡尔组合](#descartes_testgo-%E7%AC%9B%E5%8D%A1%E5%B0%94%E7%BB%84%E5%90%88)
    + [TestDescartes](#testdescartes)

<!-- tocstop -->

# algorithm 算法模块
## crc16_test.go crc16算法
### TestCrc16
```go

t.Log(Crc16([]byte("abcdefg汉字")))
```
### TestCrc16s
```go

t.Log(Crc16s("abcdefg汉字") == Crc16([]byte("abcdefg汉字")))
```
## descartes_test.go 笛卡尔组合
### TestDescartes
```go

result := DescartesCombine([][]string{{"A", "B"}, {"1", "2", "3"}, {"a", "b", "c", "d"}})
for _, item := range result {
	t.Log(item)
}
```
