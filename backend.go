package main 

import(
  "fmt"
  "net/http"
  "encoding/json"
  "math/rand"
)

type Enemy struct {

    ID int `json:"id"`
    Floor int `json:"floor"`
    X int `json:"x"`
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
    Player: Player{Floor: 1, Alive: true, Score: 0},
}

var enemyID = 0
func spawnEnemiesOnFloor(floor int) {
   count := rand.Intn(3) + 1 // 1 to 3 enemies per floor
    for i := 0; i < count; i++ {
        enemyID++
        state.Enemies = append(state.Enemies, Enemy{
            ID: enemyID,
            Floor: floor,
            X: rand.Intn(600) + 100, // random position on floor
        }) //close append and Enemy
    }
}
func moveUp(w http.ResponseWriter, r *http.Request) {
    if !state.Player.Alive {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(state)
        return
    }
    // climb up one floor
    state.Player.Floor++
    state.Player.Score++

// spawn enemies on the new floor
    spawnEnemiesOnFloor(state.Player.Floor)

    // check if player landed on an enemy
    checkCollision()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(state)
}

func checkCollision() {
    for _, e := range state.Enemies {
        if e.Floor == state.Player.Floor {
            // player spawns at x=50, check if any enemy is close
            if e.X < 100 {
                state.Player.Alive = false
                fmt.Println("Game over! Score:", state.Player.Score)
            }
        }
    }
}

func status(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(state)
}

func reset(w http.ResponseWriter, r *http.Request) {
    state = GameState{
        Player: Player{Floor: 1, Alive: true, Score: 0},
    }
    state.Enemies = []Enemy{}
    enemyID = 0
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(state)
}

func main() {
    http.HandleFunc("/up", moveUp)
    http.HandleFunc("/status", status)
    http.HandleFunc("/reset", reset)
    fmt.Println("Server on http://localhost:8080")
    fmt.Println("game is ruining")
    http.ListenAndServe(":8080", nil)
}
