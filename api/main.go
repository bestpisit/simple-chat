package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan Message)           // Broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// Message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// MongoDB collection
var collection *mongo.Collection

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set.")
	}
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("chatdb").Collection("messages")

	// Create a simple file server
	fs := http.FileServer(http.Dir("../web/dist"))
	http.Handle("/", fs)

	// Configure WebSocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 8080 and log any errors
	log.Println("http server started on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	// Send chat history to the newly connected client
	sendChatHistory(ws)

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// Insert message into MongoDB
		_, insertErr := collection.InsertOne(context.TODO(), msg)
		if insertErr != nil {
			log.Printf("InsertOne error: %v", insertErr)
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func sendChatHistory(ws *websocket.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Find error: %v", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var msg Message
		if err := cursor.Decode(&msg); err != nil {
			log.Printf("Cursor Decode error: %v", err)
			return
		}
		if err := ws.WriteJSON(msg); err != nil {
			log.Printf("WriteJSON error: %v", err)
			return
		}
	}
}
