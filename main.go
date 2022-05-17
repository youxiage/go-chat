package main

import (
    "fmt"
    "io"
    "net"
    "strings"
    "sync"
    "time"
)

func main() {
    server := NewServer("127.0.0.1", 8000)
    server.Start()
}

type Server struct {
    Ip        string
    Port      int
    OnlineMap map[string]*User

    l sync.RWMutex
}

func NewServer(ip string, port int) *Server {
    return &Server{
        Ip:        ip,
        Port:      port,
        OnlineMap: make(map[string]*User),
    }
}

func (s *Server) Start() {
    listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
    if err != nil {
        fmt.Println(err)
        return
    }
    defer listen.Close()

    for {
        conn, err := listen.Accept()
        if err != nil {
            fmt.Println(err)
            return
        }

        go s.Handler(conn)
    }
}

func (s *Server) Handler(conn net.Conn) {
    defer conn.Close()

    user := NewUser(conn, s)
    // 上线用户
    user.Online()
    // 广播
    s.BradCast(user, "上线了")

    buf := make([]byte, 4096)

    isAlive := make(chan bool)
    go func() {
        for {
            n, err := conn.Read(buf)
            // 客户端退出
            if n == 0 {
                user.Offline()
                user.conn.Close()
                return
            }

            if err != nil && err != io.EOF {
                fmt.Println(err)
                return
            }

            msg := string(buf[:n-1])

            // 是否指令方式
            if strings.HasPrefix(msg, `/`) {
                user.DoMessage(msg)
            } else if msg != "" {
                s.BradCast(user, msg)
            }

            isAlive<-true
        }

    }()

    for {
        select {
        case <-isAlive:

        case <-time.After(time.Minute * 2):
            user.conn.Write([]byte("系统通知，不活跃被踢出\n"))
            user.Offline()
            user.conn.Close()
        }
    }
}

// 广播
func (s *Server) BradCast(user *User, msg string) {
    //message := fmt.Sprintf("\n[%s]%s %s\n", user.Addr, user.Name, msg)
    message := fmt.Sprintf("(ip:%s)[%s]:%s\n", user.Addr, user.Name, msg)
    for _, user := range s.OnlineMap {
        user.conn.Write([]byte(message))
    }
}

