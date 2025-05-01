package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PSQLOption func(*PSQL)

func WithPostgresCngPath(postgresCngPath string) PSQLOption {
	return func(p *PSQL) {
		p.postgresCngPath = postgresCngPath
	}
}

func WithCacheSize(cacheSize int) PSQLOption {
	return func(p *PSQL) {
		p.cacheSize = cacheSize
	}
}

func NewPSQL(opts ...PSQLOption) (p *PSQL, err error) {
	p = &PSQL{
		postgresCngPath: "/home/ilia/Desktop/temki/gonya/UrlCut/configs/postgres.json",
		cacheSize:       1000,
	}
	for _, opt := range opts {
		opt(p)
	}

	var conf PostgresConfig
	conf.ParseAndInit(p.postgresCngPath)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.user, conf.password, conf.dbname)

	p.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln("Не удалось установить соединение с БД Postgres", err.Error())
		return nil, err
	}

	p.cache = NewCache(p.cacheSize)

	return
}
