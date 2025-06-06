package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

/*Postgres SQL-хранилище*/
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	user     string
	password string
	dbname   string
}

func (c *PostgresConfig) ParseAndInit(postgresCngPath string) (err error) {
	var barr []byte

	if barr, err = os.ReadFile(postgresCngPath); err == nil {
		if err = json.Unmarshal(barr, c); err != nil {
			return
		}
	} else {
		return
	}

	c.user = os.Getenv("POSTGRES_USER")
	c.password = os.Getenv("POSTGRES_PASSWORD")
	c.dbname = "urlcut"

	return
}

// Хранить данные конфигурации удобно для отладки или дампа структуры
type PSQL struct {
	postgresCngPath string
	cacheSize       int

	cache *Cache
	db    *sql.DB
}

func (p *PSQL) GetFullUrl(ctx context.Context, cutUrl string) (fullUrl string, err error) {
	if p.cache.Contains(cutUrl) {
		fullUrl, err = p.cache.Get(cutUrl)
		return
	}

	var rows *sql.Rows
	request := fmt.Sprintf("SELECT fullUrl FROM \"CutToFull\" WHERE cutUrl = '%s'", cutUrl)
	rows, err = p.db.QueryContext(ctx, request)
	if err != nil {
		return
	}
	defer rows.Close()

	/*А надо ли проверять ctx? или выходит по err?*/

	for rows.Next() {
		if err = rows.Scan(&fullUrl); err != nil {
			return
		}
	}

	if len(fullUrl) == 0 {
		err = fmt.Errorf("Элемент по ключу '" + cutUrl + "' не найден")
		return
	}

	p.cache.Add(cutUrl, fullUrl)

	return
}

func (p *PSQL) StoreCutUrl(ctx context.Context, cutUrl string, fullUrl string) (err error) {
	request := fmt.Sprintf(`	INSERT INTO "CutToFull" (cutUrl, fullUrl) VALUES ('%s', '%s') 
											ON CONFLICT (cutUrl) DO UPDATE  SET cutUrl = '%s', fullUrl = '%s'`,
		cutUrl, fullUrl, cutUrl, fullUrl)
	_, err = p.db.QueryContext(ctx, request)
	if err != nil {
		return
	}

	p.cache.Add(cutUrl, fullUrl)

	return
}

func (p *PSQL) Close() {
	p.db.Close()
}
