package main

import (
	"fmt"
	"github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/middleware"
	"encoding/json"
)


type Message struct {
	Username	string	`json:"username"`
	Message		string	`json:"message"`
	Receiver	string	`json:"receiver"`
}

var (
	upgrader = websocket.Upgrader{}
    clients = make(map[string]*websocket.Conn)
    broadcast = make(chan Message)
)


func handleConnections(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
    	fmt.Println("Nisam mogao da upalim websockete")
	}
	defer ws.Close()

	for {
		// Read
		msg := Message{}
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, 0) {
				fmt.Println("Ugasi se korisnik: ")
				delete(clients, "")
			}
			break
		}
		fmt.Println("poruka: ", message)
		fmt.Println("greska: ", err)
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println("ne mogu da procitam JSON")
		}

		fmt.Printf("%s\n", msg)

		clients[msg.Username] = ws


		broadcast <- msg

	}
	return nil
}

func handleMessages() {
	for {
		msg := <- broadcast
		sender := msg.Username
		receiver := msg.Receiver
		text := sender + ": " + msg.Message

		if receiverWs, ok := clients[receiver]; ok {
			// Write
			err := receiverWs.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
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
	e.GET("/ws", handleConnections)

	go handleMessages()


	e.Logger.Fatal(e.Start(":1323"))
}
