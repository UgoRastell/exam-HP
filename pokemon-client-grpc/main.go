package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	pb "pokemon-client-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	totalLatencyGrpc int64
	messageCountGrpc int64
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

	serverAddr := "grpc-server:50051"
	log.Printf("gRPC Client: Connexion à %s", serverAddr)

	var conn *grpc.ClientConn
	var err error

	for {
		conn, err = grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			log.Println("gRPC Client: Connecté au serveur gRPC !")
			break
		}
		log.Printf("gRPC Client: Échec connexion gRPC: %v", err)
		log.Println("gRPC Client: Nouvel essai dans 5s...")
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()

	client := pb.NewBattleServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.BattleStream(ctx)
	if err != nil {
		log.Fatalf("gRPC Client: Erreur ouverture du stream: %v", err)
	}
	log.Println("gRPC Client: Stream gRPC établi.")

	pokemonID := "grpc_pokemon_" + randomString(5)
	pokemonName := randomName()
	log.Printf("gRPC Client: Connexion du Pokémon gRPC: %s (%s)", pokemonName, pokemonID)

	done := make(chan struct{})
	updateInterval := getUpdateInterval()

	go func() {
		ticker := time.NewTicker(updateInterval)
		defer ticker.Stop()
		defer log.Println("gRPC Client: Goroutine d'envoi gRPC terminée.")

		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				log.Println("gRPC Client: Contexte annulé, arrêt de l'envoi gRPC.")
				return
			case <-ticker.C:
				x, y := randomMove()
				hp := rand.Intn(100) + 1
				action := randomAction()
				t1 := time.Now().UnixMilli()

				update := &pb.PokemonUpdate{
					Id:        pokemonID,
					Name:      pokemonName,
					X:         x,
					Y:         y,
					Hp:        int32(hp),
					Action:    action,
					TargetId:  "",
					Timestamp: t1,
				}

				if err := stream.Send(update); err != nil {
					log.Printf("gRPC Client: Erreur envoi update gRPC: %v", err)
					close(done)
					return
				}
			}
		}
	}()

	go func() {
		defer close(done)
		defer log.Println("gRPC Client: Goroutine de réception gRPC terminée.")
		for {
			statusUpdate, err := stream.Recv()
			if err == io.EOF {
				log.Println("gRPC Client: Stream gRPC fermé par le serveur (EOF).")
				return
			}
			if err != nil {
				if status.Code(err) == codes.Canceled {
					log.Println("gRPC Client: Réception gRPC annulée via contexte.")
				} else {
					log.Printf("gRPC Client: Erreur réception gRPC: %v", err)
				}
				return
			}

			t3 := time.Now().UnixMilli()

			for _, p := range statusUpdate.Pokemons {
				if p.Id == pokemonID {
					latency := t3 - p.Timestamp
					if p.Timestamp > 0 {
						atomic.AddInt64(&totalLatencyGrpc, latency)
						atomic.AddInt64(&messageCountGrpc, 1)
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
				count := atomic.LoadInt64(&messageCountGrpc)
				if count > 0 {
					avg := float64(atomic.LoadInt64(&totalLatencyGrpc)) / float64(count)
					log.Printf("gRPC Client [%s]: Latence Moyenne RTT = %.2f ms (%d mesures)", pokemonID, avg, count)
				}
			}
		}
	}()

	select {
	case <-done:
		log.Println("gRPC Client: Arrêt dû à une erreur de communication gRPC.")
	case <-interrupt:
		log.Println("gRPC Client: Interruption reçue (Ctrl+C). Fermeture gRPC...")
		cancel()
		if err := stream.CloseSend(); err != nil {
			log.Printf("gRPC Client: Erreur lors de CloseSend: %v", err)
		}
	}

	<-time.After(500 * time.Millisecond)
	count := atomic.LoadInt64(&messageCountGrpc)
	if count > 0 {
		avg := float64(atomic.LoadInt64(&totalLatencyGrpc)) / float64(count)
		log.Printf("gRPC Client [%s]: === Latence Moyenne Finale RTT = %.2f ms (%d mesures) ===", pokemonID, avg, count)
	}
	log.Println("gRPC Client: Client Pokémon gRPC terminé.")
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
