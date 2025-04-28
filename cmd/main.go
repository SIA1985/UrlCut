package main

import(
	"UrlCut/internal/storage"
	"UrlCut/internal/interfaces"
	"UrlCut/internal/cutter"

	"fmt"
)

func tryStorage(s interfaces.Storage) {
	url, _ := s.GetFullUrl("aaaaaa")
	url2, _ := s.GetFullUrl("aaaaaa")

	if url == url2 && len(url) != 0 {
		fmt.Println("Storage ok")
		return
	}

	s.Close()
}

func tryCutter(c *cutter.Cutter) {
	fmt.Println("CutterOk: " + c.Cut("www.yandex.ru"))
}

func main() {
	var err error

	var p *storage.PSQL
	p, err = storage.NewPSQL(	storage.WithCacheSize(5), 
								storage.WithPostgresCngPath("/home/ilia/Desktop/темки/гоня/UrlCut/configs/postgres.json") )
	if err != nil {
		return
	}
	tryStorage(p)


	var c *cutter.Cutter
	c, err = cutter.NewCutter(	cutter.WithCutSize(6) )
	if err != nil {
		return
	}
	tryCutter(c)
}