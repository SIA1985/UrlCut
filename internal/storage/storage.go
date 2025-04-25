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
	Host		string		`json:"host"`
	Port		int 		`json:"port"`
	user		string
	password	string 
	dbname		string	
}

func (c *PConfig) Init() {
	var err error

	//todo: переделать
	cngPath := "/home/ilia/Desktop/темки/гоня/UrlCut/configs/postgres.json"
	
	var barr []byte
	if barr, err = os.ReadFile(cngPath); err == nil {
        if err = json.Unmarshal(barr, c); err != nil {
            log.Panicln("Ошибка в парсинге конфигурационного файла " + cngPath, err.Error())
        }
    } else {
        log.Panicln("Не удалось открыть конфигурационный файл:", err.Error())
    }

	c.user = os.Getenv("POSTGRES_USER")
	c.password = os.Getenv("POSTGRES_PASSWORD")
	c.dbname = "urlcut"
}

type PSQL struct {
	cache 		*Cache
	db			*sql.DB
}

func (p *PSQL) Init() {
	var err error

	p.cache = NewCache(1000)

	var conf PConfig
	conf.Init()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    conf.Host, conf.Port, conf.user, conf.password, conf.dbname)

	p.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln("Не удалось установить соединение с БД Postgres", err.Error())
	}
	//todo: defer s.db.Close() -> функция деструктор?
}

func (p *PSQL) GetFullUrl(cutUrl string) (fullUrl string, err error) {
	if p.cache.Contains(fullUrl) {
		fullUrl, err = p.cache.Get(cutUrl)
		return
	}

	var rows *sql.Rows
	//todo: переделать
	rows, err = p.db.Query(fmt.Sprintf("SELECT fullUrl FROM \"CutToFull\" WHERE cutUrl = '%s'", cutUrl))
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&fullUrl); err != nil {
			return
		}
	}

	//todo: занесение в кэш?

	return
}

func (p *PSQL) StoreCutUrl(cutUrl string, fullUrl string) (err error) {
	/*add to sql*/

	p.cache.Add(cutUrl, fullUrl)

	return
}