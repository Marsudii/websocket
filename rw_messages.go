package main

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)

func main() {
	app := fiber.New()

	// WebSocket route
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// READ MESSAGES
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("❌ Error reading:", err)
				break
			}
			broadcast <- msg
		}

	}))

	// CALL WRITE MESSAGES FUNCTION
	go handleBroadcast()

	log.Println("🚀 Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

// WRITE MESSAGES
func handleBroadcast() {
	for {
		msg := <-broadcast
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("❌ Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
