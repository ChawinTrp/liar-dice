package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

const codex = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message represents the standard JSON format for all communications
type Message struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"` // Use RawMessage to delay parsing
}

// Payload Structs (What's inside 'Data')
type CreateRoomData struct {
	// No data needed for creation, but good to have the struct ready
}

type JoinRoomData struct {
	RoomCode string `json:"roomCode"`
}

type PlaceBidData struct {
	Quantity int `json:"quantity"`
	Face     int `json:"face"`
}

func RandomCode() string {
	b := make([]byte, 4)
	for i := range b {
		b[i] = codex[rand.Intn(len(codex))]
	}
	return string(b)
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	player := &Player{
		ID:   "player_" + RandomCode(),
		Conn: conn,
	}

	defer conn.Close()

	for {
		// 1. Read the raw message
		_, rawMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received from %s: %s", player.ID, rawMsg)

		// 2. Parse the "Envelope" (Action)
		var msg Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Println("JSON Parse Error:", err)
			continue
		}

		// 3. Route based on Action
		switch msg.Action {
		case "create_room":
			handleCreateRoom(player)
		case "join_room":
			var data JoinRoomData
			// FIX: Explicitly cast msg.Data to []byte
			if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
				log.Println("Invalid JoinRoom data:", err)
				continue
			}
			handleJoinRoom(player, data)
		case "place_bid":
			// Similar logic for bid...
			log.Println("Bid placed (logic pending)")
		default:
			log.Println("Unknown action:", msg.Action)
		}
	}
}

// Placeholder handlers
func handleCreateRoom(p *Player) {
	room := CreateRoom()
	room.Mutex.Lock()
	room.Players[p.ID] = p
	p.RoomCode = room.Code // Assuming Player has this field
	room.Mutex.Unlock()

	// Send success response
	response := map[string]string{"event": "room_created", "code": room.Code}
	p.Conn.WriteJSON(response)
}

func handleJoinRoom(p *Player, data JoinRoomData) {
	// Logic to find room and add player...
	log.Printf("Player %s trying to join room %s", p.ID, data.RoomCode)
}
