package data

import (
	"database/sql"
	"errors"
)

type UrlMappingModel struct {
	DB *sql.DB
}

func (m UrlMappingModel) Insert(urlMapping *UrlMapping) error {

	query := `
		INSERT INTO url_mapping(id, short_url, long_url)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	args := []interface{}{urlMapping.ID, urlMapping.ShortUrl, urlMapping.LongUrl}

	return m.DB.QueryRow(query, args...).Scan(&urlMapping.ID)
}

func (m UrlMappingModel) Get(shortUrl string) (*string, error) {
	
	query := `
		SELECT long_url 
		FROM url_mapping
		WHERE short_url = $1
	`

	var longUrl string

	err := m.DB.QueryRow(query, shortUrl).Scan(&longUrl)
	if err != nil {
		return nil, errors.New("url not found")
	}

	return &longUrl, nil
}