package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jo-jordan/fish-holdem-server/entity/domain"
	"github.com/jo-jordan/fish-holdem-server/entity/inbound"
	"github.com/jo-jordan/fish-holdem-server/entity/outbound"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var id int64 = 0
var connMap = make(map[string]domain.PlayerDO)
var tableMap = make(map[int64]domain.TableDO)

func game(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	mt, msg, err := c.ReadMessage()
	log.Printf("recv: %s", msg)

	gameData, err := makeGameInfo()
	err = c.WriteMessage(mt, gameData)

	playerData, err := makePlayerList()
	err = c.WriteMessage(mt, playerData)

	log.Printf("Table match successful: id: %d\n", id)
}

func login(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	remoteAddr := c.RemoteAddr().String()

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

	id = id + 1
	userConn := domain.PlayerDO{Conn: c, Username: info.Username, UserID: id}

	connMap[remoteAddr] = userConn

	token := fmt.Sprintf("Token-%s", info.Username)
	result := outbound.LoginResultInfo{Success: true, Token: token}
	resultData, err := result.Marshal()
	log.Printf("Login successful: username: %s, token: %s\n", userConn.Username, token)
	err = c.WriteMessage(mt, resultData)
}

func makePlayerList() ([]byte, error) {
	pi := outbound.PlayerInfo{
		DataType: "player_info",
		PlayerList: []outbound.PlayerList{
			{
				ID:          1,
				Name:        "lz",
				Balance:     100,
				Bet:         30,
				Status:      "wait",
				Role:        "Player",
				IsOperator:  true,
				CardsInHand: []string{"201", "303"},
			},
			{
				ID:          2,
				Name:        "xjp",
				Balance:     10000,
				Bet:         200,
				Status:      "wait",
				Role:        "Small Blind",
				IsOperator:  false,
				CardsInHand: []string{"101", "106"},
			},
			{
				ID:          3,
				Name:        "lq",
				Balance:     9000,
				Bet:         190,
				Status:      "wait",
				Role:        "Big Blind",
				IsOperator:  false,
				CardsInHand: []string{"401", "412"},
			},
		},
	}

	data, err := pi.Marshal()
	return data, err
}

func makeGameInfo() ([]byte, error) {
	id = id + 1
	gi := outbound.TableInfo{
		TableID:      id,
		TotalPot:     1000,
		Status:       "",
		Countdown:    20,
		BetRate:      "10/20",
		CardsOnTable: []string{"312", "412", "109"},
		DataType:     "game_info",
	}

	gameData, err := gi.Marshal()

	return gameData, err
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/login", login)
	http.HandleFunc("/game", game)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
