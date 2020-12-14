# tcp2ws

Proxies all tcp traffic to websocket server

```
$ go get github.com/matti/tcp2ws
```

```
$ tcp2ws 127.0.0.1:8080 ws://echo.websocket.org
listen 127.0.0.1:8080
2020/12/14 18:44:15 127.0.0.1:62995 handling
2020/12/14 18:44:15 127.0.0.1:62995 copy 174.129.224.73:80
```

```
$ nc localhost 8080
hello
hello
```
