package main

import (
	"Octave/golibs/settings"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var missedbeats int
var beatrecived = false
var beatcount int

// Websocket handler
func handleRequests() {
	settings := settings.Settings()
	myRouter := http.NewServeMux()
	myRouter.HandleFunc("/ws", wshandler)
	fs := http.FileServer(http.Dir("./assets/"))
	myRouter.Handle("/", fs)

	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%v", settings.WSport),
		Handler: myRouter,
	}
	Log.Infof("Listening on %s", server.Addr)
	Log.Infof("Websockets listening on %s/ws", server.Addr)
	Log.Error(server.ListenAndServe().Error())
}

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 0,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024, EnableCompression: false,
}

var activeconn *websocket.Conn

func wshandler(w http.ResponseWriter, r *http.Request) {
	Log.Info("Connection received, upgrading to websocket")
	conn, err := upgrader.Upgrade(w, r, nil)
	activeconn = conn
	if err != nil {
		Log.Error(err.Error())
	}
	go heartbeat(conn)
}

func heartbeat(conn *websocket.Conn) {
	for {
		if conn != activeconn {
			break
		}
		if TimeCheck(1000) {
			WSUpdate(conn)
		}
		if TimeCheck(500) {
			dobreak := WSHeartbeat(conn)
			if dobreak {
				break
			}
		}
		time.Sleep(time.Millisecond)

	}
	Log.Info("A connection was made from the same IP and port more recently than this one, closing connection")

	if err := conn.Close(); err != nil {
		Log.Error(err.Error())
	}
}

// Will return true every x milliseconds
func TimeCheck(every int64) bool {
	current := time.Now().UnixMilli()

	return current%every == 0
}

func WSUpdate(conn *websocket.Conn) {

	out := make(map[string]string)
	out["type"] = "update"
	out["time"] = strconv.Itoa(int(time.Now().Unix()))
	jsonform, err := json.Marshal(out)
	if err != nil {
		Log.Error(err.Error())
	}
	err = conn.WriteMessage(websocket.TextMessage, jsonform)
	if err != nil {
		Log.Error(err.Error())
	}
}

func WSHeartbeat(conn *websocket.Conn) bool {
	out := make(map[string]string)
	out["type"] = "heartbeat"
	out["time"] = strconv.Itoa(int(time.Now().Unix()))
	out["id"] = strconv.Itoa(beatcount)
	beatcount++
	jsonform, err := json.Marshal(out)
	if err != nil {
		Log.Error(err.Error())
	}
	err = conn.WriteMessage(websocket.TextMessage, jsonform)
	if err != nil {
		Log.Error(err.Error())
	}

	msgcont := make(map[string]string)
	var msg []uint8
	for {
		_, msg, err = conn.ReadMessage()

		if err != nil {
			Log.Error(err.Error())

			break
		}
		if len(msg) > 0 {

			break
		}
	}
	err = json.Unmarshal(msg, &msgcont)
	if err != nil {
		Log.Error(err.Error())
	}
	if msgcont["type"] == "heartbeat" {
		beatrecived = true
		missedbeats = 0
	}
	if !beatrecived {
		missedbeats++
	}
	if missedbeats > 1 {
		Log.WarningF("Lost connection to client, missed %d beats", missedbeats)

		return true
	}
	beatrecived = false

	return false
}
