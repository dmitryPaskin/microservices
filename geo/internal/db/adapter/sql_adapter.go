package adapter

import (
	"fmt"
	"log"
	"microservices/geo/internal/models"

	"github.com/jmoiron/sqlx"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=SQLAdapterer
type SQLAdapterer interface {
	Select(query string) (models.Address, error)
	Insert(query string, lat string, lon string) error
}

type SQLAdapter struct {
	db *sqlx.DB
}

func NewSQLAdapter(db *sqlx.DB) *SQLAdapter {
	return &SQLAdapter{
		db: db,
	}
}

func (s *SQLAdapter) Select(query string) (models.Address, error) {

	rows, err := s.db.Query(`
		SELECT address.lat, address.lon
		FROM history_search_address
		JOIN search_history ON search_history.id = history_search_address.search_id
		JOIN address ON history_search_address.address_id = address.id
		WHERE levenshtein(search_history.query, $1) <= LENGTH($1) * 0.3;
		`, query)
	if err != nil {
		return models.Address{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		log.Println("совпадений нет")
		return models.Address{}, fmt.Errorf("совпадений нет")
	}

	var address models.Address

	if err = rows.Scan(&address.Lat, &address.Lon); err != nil {
		return models.Address{}, err
	}
	return address, nil
}

func (s *SQLAdapter) Insert(query string, lat string, lon string) error {
	_, err := s.db.Exec(`INSERT INTO search_history (query) VALUES ($1)`, query)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT INTO address (lat, lon) VALUES ($1, $2)`, lat, lon)
	if err != nil {
		return err
	}

	var searchID int
	err = s.db.QueryRow(`SELECT id FROM search_history ORDER BY id DESC LIMIT 1`).Scan(&searchID)
	if err != nil {
		return err
	}

	var addressID int
	err = s.db.QueryRow(`SELECT id FROM address ORDER BY id DESC LIMIT 1`).Scan(&addressID)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT INTO history_search_address (search_id, address_id) VALUES ($1, $2)`, searchID, addressID)
	if err != nil {
		return err
	}

	return nil
}
