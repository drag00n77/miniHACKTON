package main 

import(
  "fmt"
  "net/http"
)

func game(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Game is running!")
}

func main() {
    http.HandleFunc("/", game)
    fmt.Println("Server on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
