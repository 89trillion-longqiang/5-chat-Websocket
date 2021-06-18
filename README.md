#Websocket + protobuf 实现聊天服务

#整体框架
```
运用websocket开发一个聊天服务
```

#目录结构
```
├── README.md               #介绍
├── go.mod
├── handle                  #处理router的请求
│   └── handle.go
├── log                     #日志文件
│   └── info.log
├── main.go                 #代码入口
├── module                  #数据模型
│   ├── client.go
│   ├── const.go
│   └── hub.go
├── newClient.go            #测试代码
├── protobuf                
│   ├── com.pb.go
│   └── com.proto
├── router                  #路由
│   └── router.go
└── service             
    └── service.go

```

#运行方式

```
代码运行方式 go run main.go
测试代码运行方式 go run newClient.go -user username  //-user后输入的为username
```

#使用方法
```
运行后，服务器接受客户端发来的请求，服务端进行对应的恢复。
```