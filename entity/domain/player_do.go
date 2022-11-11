package domain

import "github.com/gorilla/websocket"

type PlayerDO struct {
	UserID   int64
	Username string
	Conn     *websocket.Conn
}
