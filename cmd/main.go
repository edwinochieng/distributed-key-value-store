package main

import (
	"fmt"
	"net/http"

	"github.com/edwinochieng/distributed-key-value-store/api"
	"github.com/edwinochieng/distributed-key-value-store/internal/storage"
)

func main() {
    store := storage.NewStore()
    err := store.LoadFromFile()
    if err != nil {
        fmt.Println("Failed to load data:", err)
    }

    api := api.NewAPI(store)

    http.HandleFunc("/set", api.SetHandler)
    http.HandleFunc("/get", api.GetHandler)
    http.HandleFunc("/delete", api.DeleteHandler)

    fmt.Println("Server running on :8080...")
    http.ListenAndServe(":8080", nil)
}
