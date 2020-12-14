# tcp2ws

Forwards all tcp traffic to websocket server

```
go run cmd/tcp2ws/main.go 127.0.0.1:8080 ws://echo.websocket.org
listen 127.0.0.1:8081
2020/12/14 18:44:15 127.0.0.1:62995 handling
2020/12/14 18:44:15 127.0.0.1:62995 copy 174.129.224.73:80
```

```
$ nc localhost 8080
hello
hello
```
