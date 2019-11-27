# 数据下发服务

### 命令
```
./issue start|restart|stop 

-c 配置文件
-d 后台运行

TODO: restart
```

### 生成go protos
```
go get github.com/golang/protobuf/protoc-gen-go

protoc --go_out=plugins=grpc:. ./protos/*.proto
```

## 服务端

#### 支持http和grpc方式发布数据

http:

``` go
body := service.ReleaseBody{
    Action:  "xxx",
    UniqIds: []string{"x"},
    Data:    map[string]string{"x": "xxx"},
}

data, err := json.Marshal(body)
if err != nil {
    panic(err)
}

req := bytes.NewBuffer(data)

resp, err := http.Post("http://127.0.0.1/release","application/json", req)
``` 

rpc:

``` go
con, err := grpc.Dial("127.0.0.1", grpc.WithInsecure())
if err != nil {
    panic(err)
}

cli := protos.NewIReleaseServiceClient(con)

resp, err = cli.Release(
    context.Background(),
    &protos.ReleaseBody{
        Action:"xxx",
        UniqIds:[]string{"x"},
        Data: map[string]string{"x": "xxxx"},
    })
```

## 客户端

服务器默认会有30秒间隔的保活心跳包检测，客户端收到 "\_heartbeat\_" 消息时需回应 "NOP" 消息证明客户端存活
服务器会在检测5次都未收到客户端响应之后断开连接
客户端可以在1分钟内未收到服务器的心跳包断开之前的连接重新连接

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