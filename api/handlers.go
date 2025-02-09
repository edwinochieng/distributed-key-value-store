package api

import (
	"encoding/json"
	"net/http"

	"github.com/edwinochieng/distributed-key-value-store/internal/storage"
)

type API struct {
    Store *storage.Store
}

func NewAPI(store *storage.Store) *API {
    return &API{Store: store}
}

func (api *API) SetHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var data map[string]string
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    for key, value := range data {
        api.Store.Set(key, value)
    }

    api.Store.SaveToFile()
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Key-Value pair stored successfully"))
}

func (api *API) GetHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if key == "" {
        http.Error(w, "Key is required", http.StatusBadRequest)
        return
    }

    value, found := api.Store.Get(key)
    if !found {
        http.Error(w, "Key not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"key": key, "value": value})
}

func (api *API) DeleteHandler(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if key == "" {
        http.Error(w, "Key is required", http.StatusBadRequest)
        return
    }

    api.Store.Delete(key)
    api.Store.SaveToFile()

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Key deleted successfully"))
}
