# go socket chat
通过map简单实现单人多人聊天简单例子

### 服务端启动
```
go run .
```

### 客户端连接
```
nc 127.0.0.1 8000
```

### 聊天及相关指令
```
查询在线用户 /list
改名字 /rename|xx
私聊 /to|xx|msg
群聊，随意发送消息会广播所有在线用户
超时下线，两分钟未操作踢下线
```
