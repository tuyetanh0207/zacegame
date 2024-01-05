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
	// "math/rand"
	// "time"
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
// func randomPositionForClient () int{

// }
func handleClient(client *Client) {
    defer func() {
        log.Printf("Client %d disconnected\n", client.ID)
        delete(clients, client.ID)
        close(client.Send)
    }()
	message := Message{
		ClientID: client.ID, 
		Type:"AssignPosition",
		Content: "Nothings",
		Position: 39,
		Score: 0,
		Matrix: matrix,
		Competitors: competitors,
	}
	err := client.Conn.WriteJSON(message)
	if err != nil {
		log.Printf("Error sending message to client %d: %v\n", client.ID, err)
		return
	}

    // for {
    //     select {
    //     case message, ok := <-client.Send:
    //         if !ok {
    //             return
    //         }

    //         err := client.Conn.WriteJSON(message)
    //         if err !
	//= nil {
    //             log.Printf("Error sending message to client %d: %v\n", client.ID, err)
    //             return
    //         }
    //     }
    // }
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

	client := &Client{
		ID:   nextClientID,
		Conn: conn,
		Score: 0,
		Position: 39,
		Send:     make(chan Message),

	}
	competitor := &Competitor{
		ID:   nextCompetitorID,
		Score: 0,
		Position: 39,
		Status: "Active",
	}
	fmt.Println("client", client.ID, client.Score, client.Position)
	fmt.Println("competitor", competitor.ID, competitor.Score)

	
	nextClientID++
	nextCompetitorID++
	fmt.Println("nextClientID", nextClientID)
	fmt.Println("nextCompetitorID", nextCompetitorID)
	clients[client.ID] = client
	competitors[competitor.ID] = competitor
	fmt.Printf("Client %d connected\n", client.ID)
	message := Message{
        ClientID:    client.ID,
        Type:        "AssignPosition",
        Content:     "Nothings",
        Position:    39,
        Score:       0,
        Matrix:      matrix,
        Competitors: competitors,
    }
	
	err = client.Conn.WriteJSON(message)
    if err != nil {
        log.Printf("Error sending message to client %d: %v\n", client.ID, err)
        return
    }
	
	

	for _, otherClient := range clients {
		if otherClient.ID != client.ID {
			err := otherClient.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error broadcasting message to client %d: %v\n", otherClient.ID, err)
			}
		}
	}
	//go handleClient(client)
		
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Client %d disconnected\n", client.ID)
			delete(clients, client.ID)
			return
		}

		fmt.Printf("Received from client %d: %s\n", messageType, client.ID, p)

		//Broadcast the message to all connected clients
		// for _, otherClient := range clients {
		// 	if otherClient.ID != client.ID {
		// 		err := otherClient.Conn.WriteJSON(message)
		// 		if err != nil {
		// 			log.Printf("Error broadcasting message to client %d: %v\n", otherClient.ID, err)
		// 		}
		// 	}
		// }
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
