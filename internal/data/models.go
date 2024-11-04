package data

import "database/sql"

type LongUrl struct {
	Value string `json:"value"`
}

type UrlMapping struct {
	ID int64
	LongUrl string
	ShortUrl string
}

type Models struct {
	UrlMappings UrlMappingModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		UrlMappings: UrlMappingModel{DB: db},
	}
}