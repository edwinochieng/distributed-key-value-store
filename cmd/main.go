package main

import (
	"fmt"

	"github.com/edwinochieng/distributed-key-value-store/internal/storage"
)

func main() {
    store := storage.NewStore()
    err := store.LoadFromFile()
    if err != nil {
        fmt.Println("Failed to load data:", err)
    }

    store.Set("foo", "bar")
    value, _ := store.Get("foo")
    fmt.Println("Value:", value)

    store.SaveToFile()
}
