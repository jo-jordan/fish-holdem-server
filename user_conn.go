package main

import "github.com/gorilla/websocket"

type UserConn struct {
	UserID   int64
	Username string
	conn     *websocket.Conn
}
