package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID       string          // Unique ID
	Conn     *websocket.Conn // The active connection
	Dice     []int           // Their current dice (e.g., [3, 5, 1, 6])
	RoomCode string          // Which room they are in
}

type GameRoom struct {
	Code      string             // The 4-letter code (e.g., "ABCD")
	Players   map[string]*Player // A slice of players
	TurnIndex int
	LastBid   Bid        // The current bid on the table
	Mutex     sync.Mutex // CRITICAL: To prevent race conditions
}

func CreateRoom() *GameRoom {
	code := RandomCode()

	// Ensure code uniqueness (simple check)
	for Rooms[code] != nil {
		code = RandomCode()
	}

	room := &GameRoom{
		Code:      code,
		Players:   make(map[string]*Player), // Initialize the map
		TurnIndex: 0,
		LastBid:   Bid{},
		// Mutex doesn't need explicit initialization, zero value is usable
	}

	Rooms[code] = room
	return room
}

type Bid struct {
	Quantity int
	Face     int // 1-6
	PlayerID string
}

func (r *GameRoom) NextTurn() {
	r.TurnIndex = (r.TurnIndex + 1) % len(r.Players)
}

func (r *GameRoom) RollDice() {
	// Loop through r.Players
	// Generate random dice for each
	// Send "RoundStart" event with their specific dice
}

func (r *GameRoom) PlaceBid(playerID string, quantity int, face int) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	// 1. Check if it's this player's turn
	// 2. Check if the bid is valid (higher than r.LastBid)
	// 3. Update r.LastBid
	// 4. r.NextTurn()
	// 5. Broadcast "NewBid" event
	return nil
}

func (r *GameRoom) Challenge(playerID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	// 1. Check turn
	// 2. Reveal all dice (Broadcast "Reveal" event)
	// 3. Calculate result (Wild 1s logic)
	// 4. Remove a die from the loser
	// 5. Check for elimination (0 dice)
	// 6. Start new round or End Game
}

var Rooms = make(map[string]*GameRoom)
