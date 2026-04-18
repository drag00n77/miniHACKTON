package main 

import(
  "fmt"
  "net/http"
  "encoding/json"
  "math/rand"
  "sync"
  "strings"
  "time"
)

type Enemy struct {

    ID int `json:"id"`
    Floor int `json:"floor"`
    X int `json:"x"`
    Speed int `json:"speed"`
}
type Player struct {
    Floor int  `json:"floor"`
    Alive bool `json:"alive"`
    Score int  `json:"score"`
}
type GameState struct {
    Player  Player  `json:"player"`
    Enemies []Enemy `json:"enemies"`
}
var state = GameState{
    Player: Player{Floor: 1, Alive: true, Score: 0},}
var enemyID = 0
var mu  sync.Mutex

func spawnEnemiesOnFloor(floor int) {
   count := rand.Intn(3) + 2
    for i := 0; i < count; i++ {
	enemyID++
	baseSpeed := 5 + (state.Player.Score / 5)
    if baseSpeed > 25 {
	baseSpeed = 25
}
state.Enemies = append(state.Enemies, Enemy{
      ID: enemyID,
      Floor: floor,
      X: rand.Intn(600) + 100,
      Speed: baseSpeed + rand.Intn(5),
        }) //close append and Enemy
    }
}

func checkCollision() {
    for _, e := range state.Enemies {
   	if e.Floor == state.Player.Floor && e.X < 120 && e.X > 20 {
	state.Player.Alive = false
        }
    }
}
func moveEnemies() {
	for {
	time.Sleep(50 * time.Millisecond)
	mu.Lock()
	if state.Player.Alive {
	alive := []Enemy{}
	for _, e := range state.Enemies {
	e.X -= e.Speed
	if e.X <= 0 {
	e.X = 700
	}
	alive = append(alive, e)
	}
	state.Enemies = alive
	checkCollision()
	}

	mu.Unlock()
	}
}

func spawnLoop() {
	for {
	mu.Lock()
	score := state.Player.Score
	mu.Unlock()
	wait := 2000 - (score * 20)
	if wait < 300 {
		wait = 300
	}
	time.Sleep(time.Duration(wait) * time.Millisecond)
	mu.Lock()
		if state.Player.Alive {
		floors := 1
		if score > 20 {
		  floors = 2
		}
		if score > 40 {
		floors = 3
		}
		for i := 0; i < floors; i++ {
		floor := rand.Intn(5) + 1
		spawnEnemiesOnFloor(floor)
			}
		}
		mu.Unlock()
	}
}


func gameHandler(w http.ResponseWriter, r *http.Request) {
mu.Lock()
defer mu.Unlock()

input := strings.TrimPrefix(r.URL.Path, "/")

if state.Player.Alive {
if input == "up" {
	state.Player.Floor++
	if state.Player.Floor > 5 {
	state.Player.Floor = 5
	}

	state.Player.Score++
	spawnEnemiesOnFloor(state.Player.Floor)
	checkCollision()
}
if input == "down" && state.Player.Floor > 1 {
	state.Player.Floor--
	checkCollision()
	}
}

if input == "reset" {
	state = GameState{Player: Player{Floor: 1, Alive: true}}
	enemyID = 0
}

w.Header().Set("Content-Type", "application/json")
w.Header().Set("Access-Control-Allow-Origin", "*")
json.NewEncoder(w).Encode(state)
}

func main() {
    go moveEnemies()
    go spawnLoop()
    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.HandleFunc("/up",  gameHandler)
    http.HandleFunc("/down", gameHandler)
    http.HandleFunc("/status", gameHandler)
    http.HandleFunc("/reset",  gameHandler)
    fmt.Println("Server on http://localhost:8080")
   // fmt.Println("game is ruining")
    http.ListenAndServe(":8080", nil)
}
