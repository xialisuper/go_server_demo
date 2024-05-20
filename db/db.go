package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {

	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	fmt.Println("New database ready")

	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		ID:   len(dbStructure.Chirps) + 1,
		Body: body,
	}
	dbStructure.Chirps[chirp.ID] = chirp

	if err := db.writeDB(dbStructure); err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {

	// check if file exists
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		// create file
		fmt.Println("Creating new database file")
		f, err := os.Create(db.path)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {

	// read data from file database.json
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	var dbStructure DBStructure

	if len(data) == 0 {
		dbStructure = DBStructure{
			Chirps: make(map[int]Chirp),
		}
		return dbStructure, nil
	}

	if err := json.Unmarshal(data, &dbStructure); err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	if err := os.WriteFile(db.path, data, 0644); err != nil {
		return err
	}

	return nil
}
