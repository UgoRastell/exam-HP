package main

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	pb "grpc-server/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBattleServiceServer
	pokemons sync.Map
}

func (s *server) BattleStream(stream pb.BattleService_BattleStreamServer) error {
	streamID := time.Now().UnixNano()
	log.Printf("Nouveau stream gRPC connecté: ID %d", streamID)

	var clientPokemonID string

	recvErrChan := make(chan error, 1)
	go func() {
		defer close(recvErrChan)
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Stream ID %d: Client a fermé la connexion (EOF)", streamID)
				recvErrChan <- nil
				return
			}
			if err != nil {
				log.Printf("Stream ID %d: Erreur réception gRPC: %v", streamID, err)
				recvErrChan <- err
				return
			}

			if clientPokemonID == "" && in.Id != "" {
				clientPokemonID = in.Id
				log.Printf("Stream ID %d associé à Pokemon ID: %s", streamID, clientPokemonID)
			}

			s.pokemons.Store(in.Id, in)
		}
	}()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	defer func() {
		log.Printf("Stream gRPC déconnecté: ID %d (Pokemon: %s)", streamID, clientPokemonID)
		if clientPokemonID != "" {
			s.pokemons.Delete(clientPokemonID)
			log.Printf("Pokemon %s retiré de l'état global.", clientPokemonID)
		}
	}()

	for {
		select {
		case <-ticker.C:
			var currentPokemons []*pb.PokemonUpdate
			s.pokemons.Range(func(key, value interface{}) bool {
				pokemon, ok := value.(*pb.PokemonUpdate)
				if ok {
					currentPokemons = append(currentPokemons, pokemon)
				}
				return true
			})

			if len(currentPokemons) == 0 {
				continue
			}

			statusUpdate := &pb.BattleStatus{Pokemons: currentPokemons}
			if err := stream.Send(statusUpdate); err != nil {
				log.Printf("Stream ID %d: Erreur envoi gRPC: %v", streamID, err)
				return err
			}

		case err := <-recvErrChan:
			// La goroutine de réception s'est terminée (erreur ou EOF)
			if err != nil {
				log.Printf("Stream ID %d: Arrêt suite à erreur de réception: %v", streamID, err)
				return err
			} else {
				log.Printf("Stream ID %d: Arrêt suite à fermeture client (EOF).", streamID)
				return nil
			}

		case <-stream.Context().Done():
			log.Printf("Stream ID %d: Contexte du stream terminé: %v", streamID, stream.Context().Err())
			return stream.Context().Err()
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Échec de l'écoute: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBattleServiceServer(s, &server{})

	log.Println("Serveur gRPC prêt sur le port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Erreur serveur gRPC: %v", err)
	}
}
