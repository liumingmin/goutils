# algorithm 算法模块
## crc16_test.go crc16算法
### TestCrc16
```go

t.Log(Crc16([]byte("abcdefg")))
```
## descartes_test.go 笛卡尔组合
### TestDescartes
```go

result := DescartesCombine([][]string{{"A", "B"}, {"1", "2", "3"}, {"a", "b", "c", "d"}})
for _, item := range result {
	t.Log(item)
}
```
