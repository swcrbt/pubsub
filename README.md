# 实时发布订阅服务

### 命令
```
./pubsub -c 配置文件
```

### 生成go protos
```
go get github.com/golang/protobuf/protoc-gen-go

protoc --go_out=plugins=grpc:. ./protos/*.proto
```

## 服务端发布数据

``` go
body := service.PublishMessage{
    Topics: []string{"pgygame_173060"},
    Action: "xxx",
    Body:   map[string]string{"a": "b", "b": "c"},
}

data, err := json.Marshal(body)
if err != nil {
    panic(err)
}

req := bytes.NewBuffer(data)

resp, err := http.Post("http://127.0.0.1/release","application/json", req)
``` 

## 客户端订阅数据

``` javascript
websocket = new WebSocket("ws://127.0.0.1:8888/subscribe?key=xxx");

websocket.onopen = function(evt) {
    console.log('open: ');
    console.dir(evt);
};
websocket.onclose = function(evt) {
    console.log('close: ');
    console.dir(evt);
};
websocket.onmessage = function(evt) {
    console.log('message: ');
    console.dir(evt);
};
websocket.onerror = function(evt) {
    console.log('error: ');
    console.dir(evt);
};
```