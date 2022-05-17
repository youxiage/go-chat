package main

import (
    "fmt"
    "net"
    "strings"
)

type User struct {
    Name   string
    Addr   string
    conn   net.Conn
    server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
    remoteAddr := conn.RemoteAddr().String()
    return &User{
        Name:   remoteAddr,
        Addr:   remoteAddr,
        conn:   conn,
        server: server,
    }
}

// 给server对象的map添加user
func (u *User) Online() {
    u.server.l.Lock()
    u.server.OnlineMap[u.Name] = u
    u.server.l.Unlock()
}

func (u *User) Offline() {
    u.server.l.Lock()
    delete(u.server.OnlineMap, u.Name)
    u.server.l.Unlock()
}

func (u *User) SendMsg(msg string) {
    message := fmt.Sprintf("(ip:%s)[%s]:%s\n", u.Addr, u.Name, msg)
    u.conn.Write([]byte(message))
}

func (u *User) DoMessage(msg string) {
    msg = strings.TrimSpace(msg)
    if msg == `/list` {
        u.SelectUser()
    }

    // 改名 rename|xx
    if len(msg) > 7 && strings.HasPrefix(msg, `/rename|`) {
        newName := strings.Split(msg, "|")[1]
        u.Offline()
        u.Name = newName
        u.Online()
    }

    // 发送私聊
    if strings.HasPrefix(msg, `/to|`) {
        newName := strings.Split(msg, "|")[1]

        if toUser, ok := u.server.OnlineMap[newName]; ok {
            msg := strings.Split(msg, "|")[2]
            toUser.SendMsg(msg)
        } else {
            u.SendMsg("该用户不存在")
        }

    }

}

func (u *User) SelectUser() {
    msg := "\n"
    for _, user := range u.server.OnlineMap {
        msg += fmt.Sprintf("[%s]%s %s\n", user.Addr, user.Name, "在线")
    }
    u.conn.Write([]byte(msg))
}
