
## ws模块用法
```shell script
protoc --go_out=. ws/msg.proto

//lib js   
protoc --js_out=library=msg_pb_libs,binary:ws/js  ws/msg.proto

//commonjs
cd ws
protoc --js_out=import_style=commonjs,binary:js  msg.proto

cd js
npm i -g google-protobuf
npm i -g browserify
browserify msg_pb.js -o  msg_pb_dist.js
```

https://www.npmjs.com/package/google-protobuf
