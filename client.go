package gosocketio

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/kryptogo/golang-socketio/protocol"
	"github.com/kryptogo/golang-socketio/transport"
)

const (
	webSocketProtocol       = "ws://"
	webSocketSecureProtocol = "wss://"
	socketioUrl             = "/socket.io/?EIO=3&transport=websocket"
)

/**
Socket.io client representation
*/
type Client struct {
	methods
	Channel
}

/**
Get ws/wss url by host and port
*/
func GetUrl(host string, port int, secure bool) string {
	var prefix string
	if secure {
		prefix = webSocketSecureProtocol
	} else {
		prefix = webSocketProtocol
	}
	return prefix + host + ":" + strconv.Itoa(port) + socketioUrl
}

/**
connect to host and initialise socket.io protocol

The correct ws protocol url example:
ws://myserver.com/socket.io/?EIO=3&transport=websocket

You can use GetUrlByHost for generating correct url
*/
func Dial(urlStr, nsp string, tr transport.Transport) (*Client, error) {
	c := &Client{}
	c.initChannel()
	c.initMethods()

	var err error
	c.conn, err = tr.Connect(urlStr)
	if err != nil {
		return nil, err
	}

	go inLoop(&c.Channel, &c.methods)
	go outLoop(&c.Channel, &c.methods)
	go pinger(&c.Channel)

	if len(nsp) > 0 {
		nspMsg := fmt.Sprintf("4%d%s", protocol.MessageTypeOpen, nsp)
		u, _ := url.Parse(urlStr)
		if u != nil {
			nspMsg += "?" + u.RawQuery
		}

		fmt.Println("nspMsg:", nspMsg)
		err = c.conn.WriteMessage(nspMsg)
		if err != nil {
			return nil, err
		}
	}
	c.Channel.Namespace = nsp

	return c, nil
}

/**
Close client connection
*/
func (c *Client) Close() {
	closeChannel(&c.Channel, &c.methods)
}
