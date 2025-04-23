package storage


import(
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"encoding/json"
	"os"
	"fmt"
)


/*Postgres SQL-хранилище*/
type PConfig struct {
	host		string		`json:'host'`
	port		string 		`json:'port'`
	user		string
	password	string 
	dbname		string	
}

type PSQL struct {
	cache 		*Cache
	db			*sql.DB
}

func (p *PSQL) Init() {
	var err error
	var conf PConfig

	p.cache = NewCache(1000)

	/*json config parsing by path fron ENV*/

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    conf.host, conf.port, conf.user, conf.password, conf.dbname)

	p.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	//todo: defer s.db.Close() -> функция деструктор?
}

func (p *PSQL) GetFullUrl(cutUrl string) (fullUrl string, err error) {
	if p.cache.Contains(fullUrl) {
		fullUrl, err = p.cache.Get(cutUrl)
		return
	}

	/*request sql data*/

	return
}

func (p *PSQL) StoreCutUrl(cutUrl string, fullUrl string) (err error) {
	p.cache.Add(cutUrl, fullUrl)

	/*add to sql*/

	return
}