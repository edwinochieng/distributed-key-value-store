package storage

import (
	"encoding/json"
	"os"
)

const dbFile = "store.json"

func (s *Store) SaveToFile() error {
    s.mu.RLock()
    defer s.mu.RUnlock()

    file, err := os.Create(dbFile)
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

    file, err := os.Open(dbFile)
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
