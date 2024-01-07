package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"os"
	"encoding/json"
	"io"
	// "strconv"
	"math/rand"
	// "time"
	"sort"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a connected client.
type Client struct {
	ID   int
	Conn *websocket.Conn
	Position int
	Score int
	Send  chan Message // Channel for sending messages to the client
	Status string
}
type Message interface {
    GetType() string
}
type Competitor struct {
	ID   int
	Position int
	Score int
	Status string
	Direction string
	BulletCooldown int
	Color string
}
type Bullet struct {
	ID int
	Position int
	Direction string
}
type Matrix struct {
	Positions []int `json:"matrix"`
}
type AssigningPositionMessage struct {
    Type    string `json:"type"`
	ClientInfo Competitor  `json:"clientInfo"`
	Matrix Matrix `json:"matrix"`
	Competitors map[int]*Competitor `json:"competitors"`
}
func (m AssigningPositionMessage) GetType() string {
	return m.Type
}
type UpdatingCompetitorInfoMessage struct {
	Type    string `json:"type"`
	ClientInfo Competitor  `json:"clientInfo"`
}
func (m UpdatingCompetitorInfoMessage) GetType() string {
	return m.Type
}
type UpdatingBulletInfoMessage struct {
	Type    string `json:"type"`
	ClientInfo Competitor  `json:"clientInfo"`
	BulletInfo Bullet `json:"bulletInfo"`
}
func (m UpdatingBulletInfoMessage) GetType() string {
	return m.Type
}
type UpdatingStatusMessage struct{
	Type    string `json:"type"`
	StatusContent string `json:"statusContent"`
}
func (m UpdatingStatusMessage) GetType() string {
	return m.Type
}
var clients = make(map[int]*Client)
var competitors = make(map[int]*Competitor)
var nextClientID = 1
var nextCompetitorID = 1
var matrix Matrix
var gridWidth = 32
var gridHeight = 16
// func randomPositionForClient () int{

// }

func isPositionOccupiedByCompetitor(position int) bool {
	if _, ok := competitors[position]; ok {
		// fmt.Printf("Competitor position: %d\n", competitor.Position)
		return true // Position is occupied
	}
	return false // Position is not occupied
}

func containsElement(arr []int, target int) bool {
	index := sort.Search(len(arr), func(i int) bool {
		return arr[i] >= target
	})

	return index < len(arr) && arr[index] == target
}
func isCurrentPositionOccupiedByWall(currPosition int) bool {
	// fmt.Printf("Direction: %s\n", direction)
	return containsElement(matrix.Positions, currPosition);
}
func isNewPositionOccupiedByWall(currPosition int, direction string) bool {
	// fmt.Printf("Direction: %s\n", direction)
	newPosition := determineNewPositionByDirection(currPosition, direction)
	// fmt.Printf("New Position: %d\n", newPosition)

	switch direction {
	case "Up":
		if newPosition >= 0 && !isCurrentPositionOccupiedByWall(newPosition){
			return false
		}
	case "Down":
		if newPosition < 32*16 && !isCurrentPositionOccupiedByWall(newPosition){
			return false
		}
	case "Left":
		if currPosition%32 != 0 && !isCurrentPositionOccupiedByWall(newPosition) {
			return false
		}
	case "Right":
		if currPosition%32 != 31 && !isCurrentPositionOccupiedByWall(newPosition) {
			return false
		}
	}
	return true
}
func determineNewPositionByDirection(currPosition int, direction string) int {
	switch direction {
	case "Up":
		return currPosition - 32
	case "Down":
		return currPosition + 32
	case "Left":
		return currPosition - 1
	case "Right":
		return currPosition + 1
	}
	return currPosition
}

func randomPositionForClient() (Position int, Direction string){
	randomX := rand.Intn(gridWidth)
	randomY := rand.Intn(gridHeight)
	position := randomX + randomY*gridWidth
	for {
		if(!isPositionOccupiedByCompetitor(position) && !isCurrentPositionOccupiedByWall(position)) {
			break;
		}
		randomX = rand.Intn(gridWidth)
		randomY = rand.Intn(gridHeight)
		position = randomX + randomY*gridWidth
		
	}

	// Generate a random direction (e.g., "north", "south", "east", "west")
	directions := []string{"Left", "Right", "Up", "Down"}
	randomDirection := directions[rand.Intn(len(directions))]
	for {

		if(!isNewPositionOccupiedByWall(position, randomDirection)) {
			break;
		} 
		randomDirection = directions[rand.Intn(len(directions))]
		

	}
	return position, randomDirection
}

	
func handleClient(client *Client) {
		defer func() {
		log.Printf("Client %d disconnected\n", client.ID)
		delete(clients, client.ID)
		delete(competitors, client.ID)
		close(client.Send)
	}()
	
	go func() {
		for {
			select {
			case message, ok := <-client.Send:
				if !ok {
					return
				}
		
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Error sending message to client %d: %v\n", client.ID, err)
					return
				}
			}
		}
	}()
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Client %d disconnected\n", client.ID)
			delete(clients, client.ID)
			delete(competitors, client.ID)
			return
		}

		// Unmarshal the received JSON into a generic map
		var rawData map[string]interface{}
		if err := json.Unmarshal(p, &rawData); err != nil {
			log.Printf("Error unmarshalling JSON from client %d: %v\n", client.ID, err)
			continue
		}

		// Extract the message type
		messageType, ok := rawData["type"].(string)
		if !ok {
			log.Printf("Error extracting message type from client %d\n", client.ID)
			continue
		}

		// Handle different message types
		switch messageType {
		case "move":
			// Handle move message
			// Example: client.Move(rawData["direction"].(string))
		case "shoot":
			// Handle shoot message
			// Example: client.Shoot(rawData["direction"].(string))
		// Add more cases for other message types as needed
		default:
			log.Printf("Unknown message type from client %d: %s\n", client.ID, messageType)
		}
	}
}
// handleConnections handles WebSocket connections.
func handleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling WebSocket connection...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	//conn.SetReadDeadline(time.Now().Add(1 * time.Hour))
	randomPosition, randomDirection := randomPositionForClient();
	//randomPosition, randomDirection := 12, "Right";
	client := &Client{
		ID:   nextClientID,
		Conn: conn,
		Score: 0,
		Position: randomPosition,
		Send:     make(chan Message),
	}
	competitor := &Competitor{
		ID:   nextCompetitorID,
		Score: 0,
		Position: randomPosition,
		Status: "Active",
		Direction: randomDirection,
	}
	
	// fmt.Println("client", client.ID, client.Score, client.Position)
	// fmt.Println("competitor", competitor.ID, competitor.Score)
	nextClientID++
	nextCompetitorID++
	// fmt.Println("nextClientID", nextClientID)
	// fmt.Println("nextCompetitorID", nextCompetitorID)
	clients[client.ID] = client
	competitors[competitor.ID] = competitor
	fmt.Printf("Client %d connected\n", client.ID)
	assigningPositionMessage :=AssigningPositionMessage{
		Type:"assignPositionForNewClient",
		Matrix: matrix,
		ClientInfo: *competitor,
		Competitors: competitors,
	}
	err = conn.WriteJSON(assigningPositionMessage)
	if err != nil {
		log.Printf("Error sending message to client %d: %v\n", client.ID, err)
		return
	}
	updatingCompetitorInfoMessage := UpdatingCompetitorInfoMessage{
		Type:"hasNewClient",
		ClientInfo: *competitor,
	}
	go handleClient(client)
	BroadcastMessage(client.ID, updatingCompetitorInfoMessage);
	
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Client %d disconnected\n", client.ID)
			delete(clients, client.ID)
			delete(competitors, client.ID)

			return
		}

		fmt.Printf("Received from client %d: %s\n", client.ID, p)

	}
	
	
}

func BroadcastMessage(exceptedClientID int, message Message) {
	fmt.Println("broadcasting message to client")
	for _, client := range clients {
		if client.ID != exceptedClientID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(clients, client.ID)
				delete(competitors, client.ID)
			}
		}
	}
}
		
func loadConfig(filename string) (Matrix, error) {
	var matrix Matrix

	file, err := os.Open(filename)
	if err != nil {
		return matrix, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		if err := decoder.Decode(&matrix); err == io.EOF {
			break
		} else if err != nil {
			return matrix, err
		}
	}

	return matrix, nil
}

// main function sets up the HTTP routes and starts the server.

func main() {
	var err error
	matrix, err = loadConfig("matrix.json")
	
	if err != nil {
		//fmt.Println("Error reading config file:", err)
		return
	}


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Server is running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
	
}
