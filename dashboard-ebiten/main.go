package main

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	pb "dashboard-ebiten/proto"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

type PokemonGRPC struct {
	ID        string
	Name      string
	X         float64
	Y         float64
	HP        int32
	Action    string
	TargetID  string
	Timestamp int64
}

var (
	pokemons      = make(map[string]interface{})
	pokemonsMutex sync.RWMutex
	connWS        *websocket.Conn
	streamGRPC    pb.BattleService_BattleStreamClient
	connGRPC      *grpc.ClientConn
	mode          string
	images        map[string]*ebiten.Image
	lastKeyPress  time.Time
)

const (
	screenWidth       = 800
	screenHeight      = 600
	addClientCooldown = 300 * time.Millisecond
)

type Game struct{}

func runReceiver(ctx context.Context) {
	log.Printf("Démarrage du récepteur en mode: %s", mode)

	if mode == "ws" {
		for {
			select {
			case <-ctx.Done():
				log.Println("Contexte annulé, arrêt du récepteur WebSocket.")
				return
			default:
				if connWS == nil {
					log.Println("Connexion WebSocket non établie, attente...")
					time.Sleep(2 * time.Second)
					continue
				}
				var receivedPokemons []*Pokemon
				err := connWS.ReadJSON(&receivedPokemons)
				if err != nil {
					log.Printf("Erreur réception WS: %v. Tentative de reconnexion...", err)
					connWS.Close()
					connWS = nil
					go func() {
						u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
						for {
							select {
							case <-ctx.Done():
								log.Println("Arrêt de la tentative de reconnexion WS (contexte annulé).")
								return
							default:
							}

							var errReconnect error
							connWS, _, errReconnect = websocket.DefaultDialer.Dial(u.String(), nil)
							if errReconnect == nil {
								log.Println("Reconnexion WebSocket réussie !")
								return
							}
							log.Printf("Échec reconnexion WS: %v, nouvel essai dans 5s...", errReconnect)
							time.Sleep(5 * time.Second)
						}
					}()
					time.Sleep(1 * time.Second)
					continue
				}

				pokemonsMutex.Lock()
				pokemons = make(map[string]interface{})
				for _, p := range receivedPokemons {
					pCopy := *p
					pokemons[p.ID] = &pCopy
				}
				pokemonsMutex.Unlock()
			}
		}
	} else if mode == "grpc" {
		log.Println("Mode gRPC pour le dashboard non entièrement implémenté pour la reconnexion.")
		for {
			select {
			case <-ctx.Done():
				log.Println("Contexte annulé, arrêt du récepteur gRPC.")
				return
			default:
				if streamGRPC == nil {
					log.Println("Stream gRPC non établi, attente...")
					if connGRPC == nil {
						serverAddr := "localhost:50051"
						log.Printf("Tentative de reconnexion gRPC à %s", serverAddr)
						var errDial error
						connGRPC, errDial = grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
						if errDial != nil {
							log.Printf("Échec reconnexion gRPC: %v", errDial)
							connGRPC = nil
							time.Sleep(5 * time.Second)
							continue
						}
						log.Println("Reconnexion gRPC réussie.")
					}

					client := pb.NewBattleServiceClient(connGRPC)
					var errStream error
					streamGRPC, errStream = client.BattleStream(ctx)
					if errStream != nil {
						log.Printf("Erreur recréation stream gRPC: %v", errStream)
						streamGRPC = nil
						if connGRPC != nil {
							connGRPC.Close()
						}
						connGRPC = nil
						time.Sleep(5 * time.Second)
						continue
					}
					log.Println("Stream gRPC rétabli.")
				}

				status, err := streamGRPC.Recv()
				if err != nil {
					log.Printf("Erreur réception gRPC: %v. Tentative de reconnexion...", err)
					streamGRPC = nil
					time.Sleep(1 * time.Second)
					continue
				}

				pokemonsMutex.Lock()
				pokemons = make(map[string]interface{})
				for _, pProto := range status.Pokemons {
					pGRPC := PokemonGRPC{
						ID:        pProto.Id,
						Name:      pProto.Name,
						X:         pProto.X,
						Y:         pProto.Y,
						HP:        pProto.Hp,
						Action:    pProto.Action,
						TargetID:  pProto.TargetId,
						Timestamp: pProto.Timestamp,
					}
					pokemons[pGRPC.ID] = &pGRPC
				}
				pokemonsMutex.Unlock()
			}
		}
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if time.Since(lastKeyPress) > addClientCooldown {
			log.Println("Touche 'A' pressée (Action de Debug/Test)")
			lastKeyPress = time.Now()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 60, G: 60, B: 60, A: 255})

	pokemonsMutex.RLock()
	defer pokemonsMutex.RUnlock()

	count := len(pokemons)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mode: %s | Pokémons visibles: %d", mode, count), 10, 10)

	for _, pInterface := range pokemons {
		var id, name, action string
		var x, y float64
		var hp int32
		var ts int64

		switch p := pInterface.(type) {
		case *Pokemon:
			id, name, action, x, y, hp, ts = p.ID, p.Name, p.Action, p.X, p.Y, p.HP, p.Timestamp
		case *PokemonGRPC:
			id, name, action, x, y, hp, ts = p.ID, p.Name, p.Action, p.X, p.Y, p.HP, p.Timestamp
		default:
			log.Printf("Type inconnu dans la map pokemons: %T", pInterface)
			continue
		}

		img, ok := images[name]
		if ok {
			op := &ebiten.DrawImageOptions{}
			w, h := img.Size()
			scaleFactor := 0.6
			op.GeoM.Scale(scaleFactor, scaleFactor)
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			op.GeoM.Translate(x, y)
			screen.DrawImage(img, op)

			textOffsetY_Name := -float64(h)*scaleFactor/2 - 5
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s(%d)-%s", name, hp, action), int(x)-20, int(y)+int(textOffsetY_Name))

			timeStr := time.UnixMilli(ts).Format("15:04:05.000")
			textOffsetY_Time := float64(h)*scaleFactor/2 + 5
			ebitenutil.DebugPrintAt(screen, timeStr, int(x)-35, int(y)+int(textOffsetY_Time))

		} else {
			placeholderSize := 20.0 * 0.6
			ebitenutil.DrawRect(screen, x-placeholderSize/2, y-placeholderSize/2, placeholderSize, placeholderSize, color.RGBA{R: 255, A: 255})

			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s(%d)-%s", name, hp, action), int(x)-20, int(y)-15)
			timeStr := time.UnixMilli(ts).Format("15:04:05.000")
			ebitenutil.DebugPrintAt(screen, timeStr, int(x)-35, int(y)+15)
		}

		_ = id
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	envMode := os.Getenv("EBT_DASHBOARD_MODE")
	flag.StringVar(&mode, "mode", "ws", "Mode de connexion : ws ou grpc")
	flag.Parse()
	if envMode != "" {
		mode = envMode
	}

	log.Printf("Initialisation du Dashboard Ebiten en mode: %s", mode)

	images = make(map[string]*ebiten.Image)
	loadImages()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if mode == "ws" {
		wsHost := "localhost:8080" // Adresse locale pour go run
		u := url.URL{Scheme: "ws", Host: wsHost, Path: "/ws"}
		log.Printf("Tentative de connexion WS initiale à: %s", u.String())
		var err error
		// ***** CORRECTION TYPE ICI *****
		var resp *http.Response // Doit être http.Response
		connWS, resp, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Printf("Erreur connexion WS initiale: %v", err)
			if resp != nil {
				log.Printf("Code réponse HTTP: %d", resp.StatusCode)
				// Lire et fermer le corps si besoin de débugger plus
				// bodyBytes, _ := io.ReadAll(resp.Body) // Nécessite import "io"
				// log.Printf("Corps réponse: %s", string(bodyBytes))
				resp.Body.Close() // Fermer le corps de réponse
			}
			connWS = nil
			log.Println("(Le récepteur tentera de reconnecter en arrière-plan)")
		} else {
			log.Println("Connexion WebSocket initiale réussie.")
			if resp != nil {
				// Fermer le corps de la réponse même si la connexion WS est réussie
				resp.Body.Close()
			}
			// Laisser connWS ouvert, sera géré par runReceiver
		}
	} else if mode == "grpc" {
		serverAddr := "localhost:50051" // Adresse locale pour go run
		log.Printf("Tentative de connexion gRPC initiale à: %s", serverAddr)
		var err error
		connGRPC, err = grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Erreur connexion gRPC initiale: %v (le récepteur tentera de reconnecter)", err)
			connGRPC = nil
		} else {
			log.Println("Connexion gRPC initiale réussie.")
			// Laisser connGRPC ouvert, sera géré par runReceiver
		}
	} else {
		log.Fatalf("Mode inconnu : '%s'. Choisir ws ou grpc", mode)
	}

	go runReceiver(ctx)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Battle Royale Pokémon Dashboard")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	log.Println("Lancement de la boucle Ebiten...")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Printf("Erreur Ebiten: %v", err)
	}

	log.Println("Arrêt du Dashboard Ebiten.")
}

func loadImages() {
	names := []string{"Pikachu", "Bulbizarre", "Salameche", "Carapuce", "Evoli"}
	basePath := "images/"

	for _, name := range names {
		path := fmt.Sprintf("%s%s.png", basePath, name)
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Printf("Erreur chargement image %s: %v. Vérifier chemin.", path, err)
			continue
		}
		images[name] = img
		log.Printf("Image chargée pour: %s depuis %s", name, path)
	}
}
