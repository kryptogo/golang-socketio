package main

import (
	"encoding/json"
	"log"
	"runtime"
	"strings"
	"time"

	gosocketio "github.com/kryptogo/golang-socketio"
	"github.com/kryptogo/golang-socketio/transport"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	url := "wss://api-v4.zerion.io/socket.io/?EIO=3&transport=websocket&api_token={your_api_token}"
	namespace := "namespace"
	c, err := gosocketio.Dial(url, "/"+namespace, transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	scopes := []string{"portfolio", "assets", "deposits", "loans", "staked-assets", "locked-assets"}
	for _, scope := range scopes {
		events := []string{"received", namespace, scope}
		event := strings.Join(events, " ")
		c.On(event, func(h *gosocketio.Channel, msg interface{}) {
			data, _ := json.Marshal(msg)
			log.Println("on received:", string(data))
		})
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.Emit("subscribe", map[string]interface{}{
		"scope": scopes,
		"payload": map[string]interface{}{
			"address":          "{YOUR_ADDRESS}",
			"currency":         "usd",
			"portfolio_fields": "all",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(60 * time.Second)
	c.Close()

	log.Println(" [x] Complete")
}
