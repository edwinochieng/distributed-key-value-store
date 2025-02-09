package storage

import (
	"encoding/json"
	"os"
	"sync"
)

type Store struct {
    mu    sync.Mutex
    data  map[string]string
}

func NewStore() *Store {
    return &Store{
        data: make(map[string]string),
    }
}

func (s *Store) Set(key, value string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    value, exists := s.data[key]
    return value, exists
}

func (s *Store) Delete(key string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.data, key)
}

func (s *Store) SaveToFile() error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    file, err := os.Create("store.json")
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    return encoder.Encode(s.data)
}

func (s *Store) LoadFromFile() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    file, err := os.Open("store.json")
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    return decoder.Decode(&s.data)
}
