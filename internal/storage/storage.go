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
type PostgresConfig struct {
	Host		string		`json:"host"`
	Port		int 		`json:"port"`
	user		string
	password	string 
	dbname		string	
}

func (c *PostgresConfig) ParseAndInit(postgresCngPath string) (err error) {
	var barr []byte
	
	if barr, err = os.ReadFile(postgresCngPath); err == nil {
        if err = json.Unmarshal(barr, c); err != nil {
            log.Println("Ошибка в парсинге конфигурационного файла " + postgresCngPath, err.Error())
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

type PSQLOption func(*PSQL)

func WithPostgresCngPath(postgresCngPath string) (PSQLOption) {
	return func (p *PSQL) {
		p.postgresCngPath = postgresCngPath
	}
}

func WithCacheSize(cacheSize int) (PSQLOption) {
	return func (p *PSQL) {
		p.cacheSize = cacheSize
	}
}

func NewPSQL(opts... PSQLOption) (p *PSQL, err error) {
	p = &PSQL{
		postgresCngPath: 	"/home/ilia/Desktop/темки/гоня/UrlCut/configs/postgres.json",
		cacheSize:  		1000,
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

//Хранить данные конфигурации удобно для отладки или дампа структуры
type PSQL struct {
	postgresCngPath		string
	cacheSize			int

	cache 		*Cache
	db			*sql.DB
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