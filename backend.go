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

func game(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Game is running!")
}

func main() {
    http.HandleFunc("/", game)
    fmt.Println("Server on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
