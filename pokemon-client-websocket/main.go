package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
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
	totalLatencyWs int64
	messageCountWs int64
)

func getUpdateInterval() time.Duration {
	intervalStr := os.Getenv("UPDATE_INTERVAL_MS")
	intervalMs, err := strconv.Atoi(intervalStr)
	if err != nil || intervalMs <= 0 {
		return 1000 * time.Millisecond
	}
	return time.Duration(intervalMs) * time.Millisecond
}

func main() {
	rand.Seed(time.Now().UnixNano())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	wsHost := "websocket-server:8080"
	u := url.URL{Scheme: "ws", Host: wsHost, Path: "/ws"}
	log.Printf("WS Client: Connexion à %s", u.String())

	var c *websocket.Conn
	var err error
	var resp *http.Response

	for {
		c, resp, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err == nil {
			log.Println("WS Client: Connecté au serveur WebSocket !")
			if resp != nil {
				resp.Body.Close()
			}
			break
		}
		log.Printf("WS Client: Échec connexion WebSocket: %v", err)
		if resp != nil {
			log.Printf("WS Client: Réponse HTTP reçue: %d", resp.StatusCode)
			resp.Body.Close()
		}
		log.Println("WS Client: Nouvel essai dans 5s...")
		time.Sleep(5 * time.Second)
	}
	defer c.Close()
	defer func() {
		log.Println("WS Client: Envoi du message de fermeture au serveur...")
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}()

	pokemonID := "ws_pokemon_" + randomString(5)
	pokemonName := randomName()
	log.Printf("WS Client: Nouveau Pokémon connecté : %s (%s)", pokemonName, pokemonID)

	done := make(chan struct{})
	updateInterval := getUpdateInterval()

	go func() {
		ticker := time.NewTicker(updateInterval)
		defer ticker.Stop()
		defer log.Println("WS Client: Goroutine d'envoi terminée.")

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				x, y := randomMove()
				hp := rand.Intn(100) + 1
				action := randomAction()
				t1 := time.Now().UnixMilli()

				p := Pokemon{
					ID:        pokemonID,
					Name:      pokemonName,
					X:         x,
					Y:         y,
					HP:        int32(hp),
					Action:    action,
					TargetID:  "",
					Timestamp: t1,
				}

				err := c.WriteJSON(p)
				if err != nil {
					log.Println("WS Client: Erreur envoi:", err)
					close(done)
					return
				}
			}
		}
	}()

	go func() {
		defer close(done)
		defer log.Println("WS Client: Goroutine de réception terminée.")
		for {
			var pokemons []Pokemon
			err := c.ReadJSON(&pokemons)
			if err != nil {
				log.Println("WS Client: Erreur réception:", err)
				return
			}

			t3 := time.Now().UnixMilli()

			for _, p := range pokemons {
				if p.ID == pokemonID {
					latency := t3 - p.Timestamp
					if p.Timestamp > 0 {
						atomic.AddInt64(&totalLatencyWs, latency)
						atomic.AddInt64(&messageCountWs, 1)
					}
					break
				}
			}
		}
	}()

	go func() {
		avgTicker := time.NewTicker(10 * time.Second)
		defer avgTicker.Stop()
		for {
			select {
			case <-done:
				return
			case <-avgTicker.C:
				count := atomic.LoadInt64(&messageCountWs)
				if count > 0 {
					avg := float64(atomic.LoadInt64(&totalLatencyWs)) / float64(count)
					log.Printf("WS Client [%s]: Latence Moyenne RTT = %.2f ms (%d mesures)", pokemonID, avg, count)
				}
			}
		}
	}()

	select {
	case <-done:
		log.Println("WS Client: Arrêt dû à une erreur de communication.")
	case <-interrupt:
		log.Println("WS Client: Interruption reçue (Ctrl+C). Fermeture...")
		close(done)
		_ = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		<-time.After(500 * time.Millisecond)
	}

	count := atomic.LoadInt64(&messageCountWs)
	if count > 0 {
		avg := float64(atomic.LoadInt64(&totalLatencyWs)) / float64(count)
		log.Printf("WS Client [%s]: === Latence Moyenne Finale RTT = %.2f ms (%d mesures) ===", pokemonID, avg, count)
	}
	log.Println("WS Client: Client Pokémon WebSocket terminé.")
}

func randomMove() (float64, float64) {
	return rand.Float64() * 800, rand.Float64() * 600
}

func randomAction() string {
	actions := []string{"ATTACK", "DEFEND", "MOVE", "REST"}
	return actions[rand.Intn(len(actions))]
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func randomName() string {
	names := []string{"Pikachu", "Bulbizarre", "Salameche", "Carapuce", "Evoli"}
	return names[rand.Intn(len(names))]
}
