# github.com/liumingmin/goutils
gotuils目标是快速搭建应用的辅助代码库,扫码加讨论群。

<img src="avatar.jpg" width="150" height="150" ></img>

<!-- toc -->

## ws模块用法
```shell script
protoc --go_out=. ws/msg.proto

//js  
protoc --js_out=library=protobuf,binary:ws/js  ws/msg.proto
```

## 常用工具库

|文件  |说明    |
|----------|--------|
|async.go|带超时异步调用|
|crc16.go |查表法crc16|
|crc16-kermit.go|算法实现crc16|
|csv_parse.go|csv解析封装|
|httputils.go|httpClient工具|
|math.go|数学库|
|models.go|反射创建对象|
|stringutils.go|字符串处理|
|struct.go|结构体工具(拷贝、合并)|
|tags.go|结构体tag工具 |                     
|utils.go|其他工具类 |  

​                     