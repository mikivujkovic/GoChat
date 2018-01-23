package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"github.com/labstack/echo/middleware"
	"github.com/gorilla/websocket"
)


type Message struct {
	Username	string	`json:"username"`
	Message		string	`json:"message"`
	Receiver	string	`json:"receiver"`
}

var (
	upgrader = websocket.Upgrader{}
    clients = make(map[string]*websocket.Conn)
)

func hello(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "greska")
	}
	defer ws.Close()

    for {
		// Read
		msg := Message{}
		err = ws.ReadJSON(&msg)
		if err != nil {
			c.Logger().Error(err)
			fmt.Println("ne mogu da procitam JSON")
		}
		fmt.Printf("%s\n", msg)

		clients[msg.Username] = ws
		sender := msg.Username
		receiver := msg.Receiver
		text := sender + ": " + msg.Message
		if receiverWs, ok := clients[receiver]; ok {
			// Write
			err := receiverWs.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				c.Logger().Error(err)
				fmt.Println("greska u slanju poruke")
			}//do something here
		}

	}

}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.File("/", "public/index.html")
	e.GET("/ws", hello)
	e.Logger.Fatal(e.Start(":1323"))
}