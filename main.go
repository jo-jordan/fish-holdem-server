package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jo-jordan/fish-holdem-server/entity/domain"
	"github.com/jo-jordan/fish-holdem-server/entity/inbound"
	"github.com/jo-jordan/fish-holdem-server/entity/outbound"
	"github.com/jo-jordan/fish-holdem-server/misc/global"
	"github.com/jo-jordan/fish-holdem-server/service/player"
	"github.com/jo-jordan/fish-holdem-server/service/table"
	"github.com/jo-jordan/fish-holdem-server/util"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func game(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	_, msg, err := c.ReadMessage()
	log.Printf("recv: %s", msg)
	baseInbound, err := inbound.UnmarshalBaseInbound(msg)
	if err != nil {
		log.Printf("Data format is wrong")
		c.Close()
		return
	}
	token := r.Header.Get("Game-Token")
	playerDO, playerExists := global.PlayerMap[token]
	if !playerExists {
		// TODO
		log.Printf("Player[%s] is not login: id: %d\n", token, playerDO.ID)
		c.Close()
		return
	}

	playerDO.Conn = c
	global.PlayerMap[token] = playerDO
	switch baseInbound.ReqType {
	case "MatchTable":
		{
			tableInfo, playerInfo := table_service.MatchTable(&playerDO)

			tableDO := global.TableMap[tableInfo.TableID]
			for _, p := range tableDO.PlayerListBySeat {
				if p.Conn == nil {
					return
				}
				err := p.Conn.WriteJSON(tableInfo)
				if err != nil {
					log.Println(err)
					continue
				}

				err = p.Conn.WriteJSON(playerInfo)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("resp: Player[%s] joined room[%d]\n", p.Username, tableDO.TableID)
			}
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

	userConn := domain.PlayerDO{Conn: c, Username: info.Username, ID: util.GenID()}

	token := fmt.Sprintf("Token-%s-%d", info.Username, userConn.ID)
	global.PlayerMap[token] = userConn
	result := outbound.LoginResultInfo{Success: true, Token: token}
	resultData, err := result.Marshal()
	log.Printf("Login successful: username: %s, token: %s\n", userConn.Username, token)
	err = c.WriteMessage(mt, resultData)
}

func main() {
	global.Init()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/login", login)
	http.HandleFunc("/game", game)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
