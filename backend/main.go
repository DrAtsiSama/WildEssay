package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dratsisama/wildrider-backend/database"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var db *gorm.DB

func initDB() {
	var err error
	db, err = database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database!")
}

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Message reçu: %s", p)
		if err := ws.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func main() {
	initDB()
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		sqlDB.Close()
	}()

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/users", createUser)
	log.Println("Serveur WebSocket et HTTP démarré sur :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Écoute échouée sur le port 8080: ", err)
	}
}
