package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jo-jordan/fish-holdem-server/entity/domain"
	"github.com/jo-jordan/fish-holdem-server/entity/inbound"
	"github.com/jo-jordan/fish-holdem-server/entity/outbound"
	"github.com/jo-jordan/fish-holdem-server/misc/global"
	player_service "github.com/jo-jordan/fish-holdem-server/service/player"
	table_service "github.com/jo-jordan/fish-holdem-server/service/table"
	"github.com/jo-jordan/fish-holdem-server/util"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	token string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.checkHeader()
		c.dispatch(message)
	}
}

func (c *Client) checkHeader() {
	_, ok := global.PlayerMap.Load(c.token)
	if !ok {
		c.conn.Close()
		log.Printf("Player[%s] is not login\n", c.token)
		return
	}
}

func (c *Client) dispatch(message []byte) {
	log.Printf("recv: %s", message)
	baseInbound, err := inbound.UnmarshalBaseInbound(message)
	if err != nil {
		log.Printf("Data format is wrong")
		return
	}

	playerAny, ok := global.PlayerMap.Load(c.token)
	if !ok {
		c.conn.Close()
		log.Printf("Player[%s] is not login\n", c.token)
		return
	}
	playerDO := playerAny.(domain.PlayerDO)

	switch baseInbound.ReqType {
	case "MatchTable":
		{
			tableInfo, playerInfo := table_service.MatchTable(&playerDO)

			tableData, err := tableInfo.Marshal()
			if err != nil {
				return
			}
			c.hub.broadcast <- tableData

			playerData, err := playerInfo.Marshal()
			if err != nil {
				return
			}
			c.hub.broadcast <- playerData

			log.Printf("resp: Player[%s] joined room[%d]\n", playerDO.Username, tableInfo.TableID)
		}
	case "ActionCall":
		{
			player_service.ActionCall()
		}
	case "ActionFold":
		{
		}
	case "ActionCheck":
		{
		}
	case "ActionRaise":
		{
		}
	case "QuitGame":
		{
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	mt, msg, err := c.ReadMessage()

	if err != nil {
		log.Println("login read message err:", err)
		c.Close()
		return
	}

	info, err := inbound.UnmarshalLoginInfo(msg)
	if err != nil {
		log.Println("login err:", err)
		c.Close()
		return
	}

	userConn := domain.PlayerDO{Username: info.Username, ID: util.GenID()}

	token := fmt.Sprintf("Token-%s-%d", info.Username, userConn.ID)
	global.PlayerMap.Store(token, userConn)
	result := outbound.LoginResultInfo{Success: true, Token: token}
	resultData, err := result.Marshal()
	log.Printf("Login successful: username: %s, token: %s\n", userConn.Username, token)
	err = c.WriteMessage(mt, resultData)
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), token: r.Header.Get("Game-Token")}
	client.hub.register <- client

	go client.readPump()
	go client.writePump()
}
