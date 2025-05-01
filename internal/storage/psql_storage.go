package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
			log.Println("Ошибка в парсинге конфигурационного файла "+postgresCngPath, err.Error())
			return
		}
	} else {
		log.Println("Не удалось открыть конфигурационный файл:", err.Error())
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

func (p *PSQL) GetFullUrl(cutUrl string) (fullUrl string, err error) {
	if p.cache.Contains(cutUrl) {
		fullUrl, err = p.cache.Get(cutUrl)
		return
	}

	var rows *sql.Rows
	request := fmt.Sprintf("SELECT fullUrl FROM \"CutToFull\" WHERE cutUrl = '%s'", cutUrl)
	rows, err = p.db.Query(request)
	if err != nil {
		return
	}
	defer rows.Close()

	/*todo: если более 1, то коллизия*/
	for rows.Next() {
		if err = rows.Scan(&fullUrl); err != nil {
			return
		}
	}

	p.cache.Add(cutUrl, fullUrl)

	return
}

func (p *PSQL) StoreCutUrl(cutUrl string, fullUrl string) (err error) {
	request := fmt.Sprintf("INSERT INTO \"CutToFull\" (cutUrl, fullUrl) VALUES ('%s', '%s')", cutUrl, fullUrl)
	_, err = p.db.Query(request)
	if err != nil {
		return
	}

	p.cache.Add(cutUrl, fullUrl)

	return
}

func (p *PSQL) Close() {
	p.db.Close()
}
