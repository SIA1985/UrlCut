package main

import (
	"UrlCut/internal/cutter"
	"UrlCut/internal/interfaces"
	"UrlCut/internal/logic"
	"UrlCut/internal/server"
	"UrlCut/internal/storage"

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

func tryLogic(l *logic.Logic) {
	var err error

	var cutUrl string
	cutUrl, err = l.CutUrl("www.yandex.ru")
	if err != nil {
		return
	}

	err = l.Redirect(cutUrl)
	if err != nil {
		return
	}

	fmt.Println("Logic ok")
}

func tryCutter(c *cutter.Cutter) {
	fmt.Println("CutterOk: " + c.Cut("www.yandex.ru"))
}

func main() {
	var err error

	var p *storage.PSQL
	p, err = storage.NewPSQL(storage.WithCacheSize(5),
		storage.WithPostgresCngPath("/home/ilia/Desktop/темки/гоня/UrlCut/configs/postgres.json"))
	if err != nil {
		return
	}
	tryStorage(p)

	var c *cutter.Cutter
	c, err = cutter.NewCutter(6)
	if err != nil {
		return
	}
	tryCutter(c)

	var l *logic.Logic
	l, err = logic.NewLogic(logic.WithStorage(p),
		logic.WithCutter(c))
	if err != nil {
		return
	}
	tryLogic(l)

	var s interfaces.UI
	s, err = server.NewHTTP(l)
	if err != nil {
		return
	}

	s.Exec()
}
