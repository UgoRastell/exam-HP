package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Pokemon struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	HP        int32   `json:"hp"`
	Action    string  `json:"action"`
	TargetID  string  `json:"target_id"`
	Timestamp int64   `json:"timestamp"`
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients   = sync.Map{}
	pokemons  = sync.Map{}
	broadcast = make(chan *Pokemon, 100)
)

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleBroadcast()

	log.Println("Serveur WebSocket lancé sur :8080 🚀")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erreur serveur WebSocket :", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erreur upgrade WebSocket: %v", err)
		return
	}
	defer ws.Close()

	log.Println("Nouveau client WebSocket connecté")
	clients.Store(ws, true)

	defer func() {
		log.Println("Client WebSocket déconnecté")
		clients.Delete(ws)

	}()

	var clientPokemonID string

	for {
		var p Pokemon
		err := ws.ReadJSON(&p)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erreur lecture WS (déconnexion client %s): %v", clientPokemonID, err)
			} else {
				log.Printf("Erreur lecture WS (client %s): %v", clientPokemonID, err)
			}
			if clientPokemonID != "" {
				pokemons.Delete(clientPokemonID)
				log.Printf("Pokemon %s retiré suite à déconnexion client.", clientPokemonID)
			}
			break
		}

		if clientPokemonID == "" {
			clientPokemonID = p.ID
			log.Printf("Client associé à Pokemon ID: %s", clientPokemonID)
		}

		pokemons.Store(p.ID, &p)

		broadcast <- &p
	}
}

func handleBroadcast() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {

		var currentPokemons []*Pokemon
		pokemons.Range(func(key, value interface{}) bool {
			pokemon, ok := value.(*Pokemon)
			if ok {
				currentPokemons = append(currentPokemons, pokemon)
			}
			return true
		})

		if len(currentPokemons) == 0 {
			continue
		}

		clients.Range(func(key, value interface{}) bool {
			client, ok := key.(*websocket.Conn)
			if !ok {
				return true
			}

			err := client.WriteJSON(currentPokemons)
			if err != nil {
				log.Printf("Erreur écriture WS vers client %v: %v", client.RemoteAddr(), err)
				client.Close()
				clients.Delete(client)
			}
			return true
		})
	}
}

// --- Older version ---
// func handleBroadcast() {
// 	for {
// 		p := <-broadcast // Attend une mise à jour
// 		// log.Printf("Broadcasting update pour Pokemon ID: %s", p.ID)

// 		var currentPokemons []*Pokemon
// 		pokemons.Range(func(key, value interface{}) bool {
// 			pokemon, ok := value.(*Pokemon)
// 			if ok {
// 				currentPokemons = append(currentPokemons, pokemon)
// 			}
// 			return true
// 		})

// 		clients.Range(func(key, value interface{}) bool {
// 			client, ok := key.(*websocket.Conn)
// 			if !ok { return true }

// 			err := client.WriteJSON(currentPokemons) // Envoie pkm avec son Timestamp T1
// 			if err != nil {
// 				log.Printf("Erreur écriture WS vers client %v: %v", client.RemoteAddr(), err)
// 				client.Close()
// 				clients.Delete(client)
// 			}
// 			return true
// 		})
// 	}
// }
