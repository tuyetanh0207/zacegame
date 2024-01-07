package main

import (
	"encoding/json"
	"log"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/websocket"

	// "strconv"
	"math/rand"
	"sort"
	"strconv"
	"time"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a connected client.
type Client struct {
	ID       int
	Conn     *websocket.Conn
	Position int
	Score    int
	Send     chan Message // Channel for sending messages to the client
	Status   string
}
type Message interface {
	GetType() string
}
type Competitor struct {
	ID             int
	Position       int
	Score          int
	Status         string
	Direction      string
	BulletCooldown int
	Color          string
}
type Bullet struct {
	ID        int
	Position  int
	Direction string
	ClientID  int
}
type Matrix struct {
	Positions []int `json:"matrix"`
}
type AssigningPositionMessage struct {
	Type        string              `json:"type"`
	ClientInfo  Competitor          `json:"clientInfo"`
	Matrix      Matrix              `json:"matrix"`
	Competitors map[int]*Competitor `json:"competitors"`
	Bullets     []Bullet            `json:"bullets"`
}

func (m AssigningPositionMessage) GetType() string {
	return m.Type
}

type UpdatingCompetitorInfoMessage struct {
	Type       string     `json:"type"`
	ClientInfo Competitor `json:"clientInfo"`
}

func (m UpdatingCompetitorInfoMessage) GetType() string {
	return m.Type
}

type UpdatingBulletInfoMessage struct {
	Type       string     `json:"type"`
	ClientInfo Competitor `json:"clientInfo"`
	BulletInfo Bullet     `json:"bulletInfo"`
}

func (m UpdatingBulletInfoMessage) GetType() string {
	return m.Type
}

type UpdatingStatusMessage struct {
	Type          string `json:"type"`
	StatusContent string `json:"statusContent"`
}

func (m UpdatingStatusMessage) GetType() string {
	return m.Type
}



var bullets = make([]Bullet, 0)
var nextClientID = 1
var nextCompetitorID = 1
var matrix Matrix
var gridWidth = 32
var gridHeight = 16
var (
	clientsMutex sync.RWMutex
	clients = make(map[int]*Client)

	competitorsMutex sync.RWMutex
	competitors = make(map[int]*Competitor)
)

// func randomPositionForClient () int{

// }

func isPositionOccupiedByCompetitor(position int, exceptedClientID int) (bool, int) {
	for _, comp := range competitors {
		if comp.Position == position && comp.ID != exceptedClientID {
			return true, comp.ID // Position is occupied
		}
	}
	return false, 0 // Position is not occupied
}

func containsElement(arr []int, target int) bool {
	index := sort.Search(len(arr), func(i int) bool {
		return arr[i] >= target
	})

	return index < len(arr) && arr[index] == target
}
func isCurrentPositionOccupiedByWall(currPosition int) bool {
	// log.Printf("Direction: %s\n", direction)
	return containsElement(matrix.Positions, currPosition) || currPosition < 0 || currPosition >= 32*16
}
func isNewPositionOccupiedByWall(currPosition int, direction string) bool {
	// log.Printf("Direction: %s\n", direction)
	newPosition := determineNewPositionByDirection(currPosition, direction)
	// log.Printf("New Position: %d\n", newPosition)

	switch direction {
	case "Up":
		if newPosition >= 0 && !isCurrentPositionOccupiedByWall(newPosition) {
			return false
		}
	case "Down":
		if newPosition < 32*16 && !isCurrentPositionOccupiedByWall(newPosition) {
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
func isClientCanMoveToPosition(clientInfo Competitor) bool {
	ok, _ := isPositionOccupiedByCompetitor(clientInfo.Position, clientInfo.ID)
	if !isCurrentPositionOccupiedByWall(clientInfo.Position) && !ok {
		return true
	}
	return false
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
func determineOldPositionByDirection(newPosition int, direction string) int {
	switch direction {
	case "Up":
		return newPosition + 32
	case "Down":
		return newPosition - 32
	case "Left":
		return newPosition + 1
	case "Right":
		return newPosition - 1
	}
	return newPosition
}

func randomPositionForClient() (Position int, Direction string) {
	randomX := rand.Intn(gridWidth)
	randomY := rand.Intn(gridHeight)
	position := randomX + randomY*gridWidth
	for {
		ok, _ := isPositionOccupiedByCompetitor(position, -1)
		if !ok && !isCurrentPositionOccupiedByWall(position) {
			break
		}
		randomX = rand.Intn(gridWidth)
		randomY = rand.Intn(gridHeight)
		position = randomX + randomY*gridWidth

	}

	// Generate a random direction (e.g., "north", "south", "east", "west")
	directions := []string{"Left", "Right", "Up", "Down"}
	randomDirection := directions[rand.Intn(len(directions))]
	for {

		if !isNewPositionOccupiedByWall(position, randomDirection) {
			break
		}
		randomDirection = directions[rand.Intn(len(directions))]

	}
	return position, randomDirection
}

func convertRawToCompetitor(rawStr string) Competitor {
	var competitor Competitor
	var raw map[string]interface{}
	err := json.Unmarshal([]byte(rawStr), &raw)
	if err != nil {
		return competitor
	}
	// Convert map values to struct fields
	if id, ok := raw["ID"].(float64); ok {
		competitor.ID = int(id)
	}

	if position, ok := raw["Position"].(float64); ok {
		competitor.Position = int(position)
	}

	if score, ok := raw["Score"].(float64); ok {
		competitor.Score = int(score)
	}

	if status, ok := raw["Status"].(string); ok {
		competitor.Status = status
	}

	if direction, ok := raw["Direction"].(string); ok {
		competitor.Direction = direction
	}
	return competitor

}
func updateOneCompetitorinArray(competitor Competitor) int {

	competitorsMutex.RLock()
	defer competitorsMutex.RUnlock()
	for i, comp := range competitors {
		if comp.ID == competitor.ID {
			competitors[i] = &competitor
			return 1
		}
	}
	return 0
}
func updateOneBulletInArray(bullet Bullet) int {
	for i, bull := range bullets {
		if bull.ID == bullet.ID {
			bullets[i] = bullet
			return 1
		}
	}
	return 0
}
func isAllowedToShoot(clientInfo Competitor) bool {
	position := clientInfo.Position
	// direction := clientInfo.Direction
	bulletCooldown := clientInfo.BulletCooldown
	ok, _ := isPositionOccupiedByCompetitor(position, clientInfo.ID)
	// log.Println("need ")
	// log.Println("ok, kk", ok, kk)
	// log.Println("isCurrentPositionOccupiedByWall", isCurrentPositionOccupiedByWall(position))
	log.Println("bulletCooldown", bulletCooldown)
	if !ok && !isCurrentPositionOccupiedByWall(position) && bulletCooldown == 0 {

		return true
	}
	return false
}
func isAllowedMoveBullet(bulletInfo Bullet) (string, int) {
	position := bulletInfo.Position
	direction := bulletInfo.Direction
	ok, competitorID := isPositionOccupiedByCompetitor(position, bulletInfo.ClientID)
	if ok {
		log.Println("is encountered competitor", competitorID)
		return "isEncounteringCompetitor", competitorID
	}
	if isCurrentPositionOccupiedByWall(position) || (direction == "Left" && position%32 == 0) || (direction == "Right" && position%32 == 31) || position < 0 {
		//log.Println("positionOccupied", position,position % 32  );
		return "isEncounteringWall", -1
	}
	if !ok && !isCurrentPositionOccupiedByWall(position) {
		return "ok", -1
	}

	return "unknown", -1
}
func removeBulletByID(bullets []Bullet, idToRemove int) []Bullet {
	var updatedBullets []Bullet

	for _, bullet := range bullets {
		if bullet.ID != idToRemove {
			updatedBullets = append(updatedBullets, bullet)
		}
	}
	return updatedBullets
}
var counthandleClientMessage = 0
func handleClientMessage(client *Client, message []byte) {
	log.Println("counthandle client message", counthandleClientMessage)
	var rawData map[string]interface{}
	if err := json.Unmarshal(message, &rawData); err != nil {
		log.Printf("Error unmarshalling JSON from client: %v\n", err)
		return
	}
	messageType, ok := rawData["type"].(string)
	if !ok {
		log.Printf("Error extracting message type from client\n")
		return
	}
	var clientInfo Competitor
	var bulletInfo Bullet
	clientInfoData, ok := rawData["clientInfo"].(map[string]interface{})
	if !ok {
		log.Printf("Error extracting clientInfo from clientRequestMoving message\n")
		return
	}

	clientInfoBytes, err := json.Marshal(clientInfoData)
	if err != nil {
		log.Printf("Error marshaling clientInfo: %v\n", err)
		return
	}

	if err := json.Unmarshal(clientInfoBytes, &clientInfo); err != nil {
		log.Printf("Error unmarshalling JSON from clientRequestMoving message: %v\n", err)
		return
	}
	// Handle different message types
	switch messageType {
	case "clientRequestMoving":

		if isClientCanMoveToPosition(clientInfo) {
			moveOneClientMessage := UpdatingCompetitorInfoMessage{
				Type:       "moveOneClient",
				ClientInfo: clientInfo,
			}
			updateOneCompetitorinArray(clientInfo)
			BroadcastMessage(clientInfo.ID, moveOneClientMessage, false)
			break
	
		} else {
			oldPosition := determineOldPositionByDirection(clientInfo.Position, clientInfo.Direction)
			clientInfo.Position = oldPosition
			moveOneClientMessage := UpdatingCompetitorInfoMessage{
				Type:       "moveOneClient",
				ClientInfo: clientInfo,
			}
			err := client.Conn.WriteJSON(moveOneClientMessage)
			if err != nil {
				log.Printf("Error sending message to client %d: %v\n", client.ID, err)
				break
			}

		}
		
	case "clientRequestShooting":
		//
		bulletInfoData, ok := rawData["bulletInfo"].(map[string]interface{})
		if !ok {
			log.Printf("Error extracting clientInfo from clientRequestMoving message\n")
			break
		}

		bulletInfoBytes, err := json.Marshal(bulletInfoData)
		if err != nil {
			log.Printf("Error marshaling clientInfo: %v\n", err)
			break
		}
		if err := json.Unmarshal(bulletInfoBytes, &bulletInfo); err != nil {
			log.Printf("Error unmarshalling JSON from clientRequestMoving message: %v\n", err)
			break
		}
		if isAllowedToShoot(clientInfo) {
			log.Println("updateScoreOfOneClient shoot before",competitors[clientInfo.ID].Score )
			clientInfo.Score = competitors[clientInfo.ID].Score - 1
			clientInfo.BulletCooldown = competitors[clientInfo.ID].BulletCooldown + 4
			competitors[clientInfo.ID].Score = clientInfo.Score
			competitors[clientInfo.ID].BulletCooldown  = clientInfo.BulletCooldown
			log.Println("updateScoreOfOneClient shoot after before",competitors[clientInfo.ID].Score )
			hasNewBulletMessage := UpdatingBulletInfoMessage{
				Type:       "hasNewBullet",
				ClientInfo: clientInfo,
				BulletInfo: bulletInfo,
			}

			bullets = append(bullets, bulletInfo)
			log.Println("BulletInfo", bullets)
			//
			
			updateScoreOfShooterMessage := UpdatingCompetitorInfoMessage{
				Type:       "updateScoreOfOneClient",
				ClientInfo: clientInfo,
			}
			statusContent := "User " + strconv.Itoa(clientInfo.ID) + " has just shot the bullet of code " + strconv.Itoa(bulletInfo.ID)
			hasNewBulletStatusMessage := UpdatingStatusMessage{
				Type:          "updateStatus",
				StatusContent: statusContent,
			}
			BroadcastMessage(clientInfo.ID, hasNewBulletStatusMessage, false)
			log.Println("Updating status", statusContent)
			BroadcastMessage(clientInfo.ID, hasNewBulletMessage, false)
			BroadcastMessage(clientInfo.ID, updateScoreOfShooterMessage, false)
			//
			for {

				isAllowedMoveBullet, encountedCompetitorID := isAllowedMoveBullet(bulletInfo)
				newBulletPosition := determineNewPositionByDirection(bulletInfo.Position, bulletInfo.Direction)
				bulletInfo.Position = newBulletPosition
				if isAllowedMoveBullet == "ok" {
					clientInfo.BulletCooldown = competitors[bulletInfo.ClientID].BulletCooldown -1
					competitors[bulletInfo.ClientID].BulletCooldown = clientInfo.BulletCooldown
					updateOneBulletInArray(bulletInfo)
					moveOneBulletMessage := UpdatingBulletInfoMessage{
						Type:       "moveOneBullet",
						ClientInfo: clientInfo,
						BulletInfo: bulletInfo,
					}

					BroadcastMessage(bulletInfo.ClientID, moveOneBulletMessage, false)
					
				}
				time.Sleep(200 * time.Millisecond)
				if isAllowedMoveBullet == "isEncounteringCompetitor" {
					log.Print(encountedCompetitorID)
					// random new position
					randomPosition, randomDirection := randomPositionForClient()
					killedCompetitor := competitors[encountedCompetitorID]
					killedCompetitor.Position = randomPosition
					killedCompetitor.Direction = randomDirection
					killedCompetitor.Score = competitors[killedCompetitor.ID].Score -5
					competitors[killedCompetitor.ID].Score = killedCompetitor.Score
					updateKilledCompetitorPositionMessage := UpdatingCompetitorInfoMessage{
						Type:       "moveOneClient",
						ClientInfo: *killedCompetitor,
					}
					

					// update score  updateScoreOfOneClient
					log.Println("updateScoreOfOneClient before",competitors[clientInfo.ID].Score )
					clientInfo.Score = competitors[clientInfo.ID].Score + 11
					competitors[clientInfo.ID].Score = clientInfo.Score
					log.Println("updateScoreOfOneClient after",competitors[clientInfo.ID].Score )
				
					updateScoreOfWinnerMessage := UpdatingCompetitorInfoMessage{
						Type:       "updateScoreOfOneClient",
						ClientInfo: clientInfo,
					}

					updateScoreOfLoserMessage := UpdatingCompetitorInfoMessage{
						Type:       "updateScoreOfOneClient",
						ClientInfo: *killedCompetitor,
					}
					BroadcastMessage(bulletInfo.ClientID, updateScoreOfWinnerMessage, false)
					BroadcastMessage(bulletInfo.ClientID, updateScoreOfLoserMessage, false)
					BroadcastMessage(encountedCompetitorID, updateKilledCompetitorPositionMessage, false)
					// remove one bullet
					bullets = removeBulletByID(bullets, bulletInfo.ID)
					log.Println("remove Bullets", bullets)
					removeOneBulletMessage := UpdatingBulletInfoMessage{
						Type:       "removeOneBullet",
						ClientInfo: clientInfo,
						BulletInfo: bulletInfo,
					}

					BroadcastMessage(bulletInfo.ClientID, removeOneBulletMessage, false)
					competitors[clientInfo.ID].BulletCooldown = 0
					statusContent := "User " + strconv.Itoa(clientInfo.ID) + "has just killed user " + strconv.Itoa(killedCompetitor.ID)
					killEnemyStatusMessage := UpdatingStatusMessage{
						Type:          "updateStatus",
						StatusContent: statusContent,
					}

					BroadcastMessage(clientInfo.ID, killEnemyStatusMessage, false)
					log.Println("Updating status", statusContent)
					return
				}
				if isAllowedMoveBullet == "isEncounteringWall" {
					//removeOneBullet
					log.Println("remove Bullets yet", bullets)
					log.Println("bullet ID", bulletInfo.ID)
					bullets = removeBulletByID(bullets, bulletInfo.ID)
					log.Println("remove Bullets", bullets)
					removeOneBulletMessage := UpdatingBulletInfoMessage{
						Type:       "removeOneBullet",
						ClientInfo: clientInfo,
						BulletInfo: bulletInfo,
					}

					BroadcastMessage(bulletInfo.ClientID, removeOneBulletMessage, false)
					return
				}
			}
		} else {
			// because not allow shooting
			break;
		}

		
	case "clientRequestLoggingOut":

		handleClientDisconnect(clientInfo)
		
	default:
		log.Printf("Unknown message type from client: %s\n", messageType)
	}

	counthandleClientMessage = counthandleClientMessage + 1
}
func isCompetitorNil(comp *Competitor) bool {
    return comp == nil || (comp.ID == 0 && comp.Position == 0 && comp.Score == 0 && comp.Status == "" && comp.Direction == "" && comp.BulletCooldown == 0 && comp.Color == "")
}
func handleClientDisconnect(clientInfo Competitor) {
	//if(!isCompetitorNil(&clientInfo)){
		log.Printf("Client %d disconnected in handleConnections\n", clientInfo.ID)
		removeOneClientMessage := UpdatingCompetitorInfoMessage{
			Type:       "removeOneClient",
			ClientInfo: clientInfo,
		}
		
		if(clients[clientInfo.ID] != nil && competitors[clientInfo.ID] != nil){
			delete(clients, clientInfo.ID)
			delete(competitors, clientInfo.ID)
		}
		BroadcastMessage(clientInfo.ID, removeOneClientMessage, true)
		statusContent := "User " + strconv.Itoa(clientInfo.ID) + " has just logged out. "
		log.Println("Updating status", statusContent)
		logOutStatusMessage := UpdatingStatusMessage{
			Type:          "updateStatus",
			StatusContent: statusContent,
		}
		BroadcastMessage(clientInfo.ID, logOutStatusMessage, true)
	
		
	//}
	return;
}
func handleClient(client *Client) {
	//var wg sync.WaitGroup
	//wg.Add(1)
	defer func() {
		log.Printf("Client %d disconnected in handleClient\n", client.ID)
		if (competitors[client.ID]!=nil){
			handleClientDisconnect(*competitors[client.ID])
		}
		
		
		close(client.Send)
	}()

	go func() {
		//defer wg.Done()
		for message := range client.Send {
			err := client.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error sending message to client %d: %v\n", client.ID, err)
				return
			}
		}
	}()

	for {
		messageType, p, err := client.Conn.ReadMessage()
		log.Printf("Received from client %d: %s\n", client.ID, p)
		if err != nil {
			log.Println("Client  disconnected in handleConnections\n", client.ID, err)
			handleClientDisconnect(*competitors[client.ID])
			return
		}
		if messageType == websocket.TextMessage {
			handleClientMessage(client, p)
		}
		
	}
}

// handleConnections handles WebSocket connections.
func handleConnections(w http.ResponseWriter, r *http.Request) {

	log.Println("Handling WebSocket connection...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	//conn.SetReadDeadline(time.Now().Add(1 * time.Hour))
	randomPosition, randomDirection := randomPositionForClient()
	//randomPosition, randomDirection := 12, "Right";
	client := &Client{
		ID:       nextClientID,
		Conn:     conn,
		Score:    0,
		Position: randomPosition,
		Send:     make(chan Message),
	}
	competitor := &Competitor{
		ID:        nextCompetitorID,
		Score:     0,
		Position:  randomPosition,
		Status:    "Active",
		Direction: randomDirection,
	}

	nextClientID++
	nextCompetitorID++

	clients[client.ID] = client
	competitors[competitor.ID] = competitor
	log.Printf("Client %d connected\n", client.ID)
	assigningPositionMessage := AssigningPositionMessage{
		Type:        "assignPositionForNewClient",
		Matrix:      matrix,
		ClientInfo:  *competitor,
		Competitors: competitors,
		Bullets:     bullets,
	}
	err = conn.WriteJSON(assigningPositionMessage)
	if err != nil {
		log.Printf("Error sending message to client %d: %v\n", client.ID, err)
		return
	}
	updatingCompetitorInfoMessage := UpdatingCompetitorInfoMessage{
		Type:       "hasNewClient",
		ClientInfo: *competitor,
	}
	statusContent := "User " + strconv.Itoa(client.ID) + " has just joined with us!"
	updatingHasNewClientInStatusMessage := UpdatingStatusMessage{
		Type:          "updateStatus",
		StatusContent: statusContent,
	}
	
	go handleClient(client)
	BroadcastMessage(client.ID, updatingCompetitorInfoMessage, true)
	BroadcastMessage(client.ID, updatingHasNewClientInStatusMessage, true)
	
	log.Println("Updating status", statusContent)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v", r)
			return
		}
	}()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Client %d disconnected in handleConnections\n", client.ID, err)

				handleClientDisconnect(*competitors[client.ID])

				return
			} else {
				log.Printf("Error reading message from client %d: %v\n", client.ID, err)
				//handleClientDisconnect(*competitors[client.ID])
				return
			}
		}

		log.Printf("Received from client %d: %s\n", client.ID, p)

		if messageType == websocket.TextMessage {
			handleClientMessage(client, p)
		}
	}
}

func BroadcastMessage(exceptedClientID int, message Message, isExceptingOneClient bool) {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()

	competitorsMutex.RLock()
	defer competitorsMutex.RUnlock()
	log.Println("broadcasting message to client")
	log.Println("type message in broadcast message", message.GetType())
	for _, client := range clients {
		// log.Println("clients in broadcast", client)
		// Check if the client is valid
		if (isExceptingOneClient && client.ID == exceptedClientID){
			continue
		}

		select {
		case client.Send <- message:

			// Message sent successfully
		default:
			// Error sending message, close the channel and remove the client
			log.Println("Error sending message in Broadcast")
			//closeClient(client)
		}
	}

}


// Close the client channel and remove it from the maps
func closeClient(client *Client) {
	clientsMutex.Lock()
	delete(clients, client.ID)
	clientsMutex.Unlock()

	competitorsMutex.Lock()
	delete(competitors, client.ID)
	competitorsMutex.Unlock()

	//close(client.Send)
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
		//log.Println("Error reading config file:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", handleConnections)

	log.Println("Server is running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
