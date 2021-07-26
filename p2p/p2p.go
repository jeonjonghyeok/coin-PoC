package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jeonjonghyeok/coin/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(w http.ResponseWriter, r *http.Request) {
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleErr(err)
	peer := initPeer(conn, ip, openPort)
	peer.inbox <- []byte("hello from 3000!")
}

func AddPeer(address, port, openPort string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	peer.inbox <- []byte("hello from 4000!")

}
