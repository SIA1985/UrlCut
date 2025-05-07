package main

import (
	"UrlCut/internal/cutter"
	"UrlCut/internal/interfaces"
	"UrlCut/internal/logic"
	"UrlCut/internal/server"
	"UrlCut/internal/storage"
	"log"

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

	var fullUrl string
	fullUrl, err = l.GetFullUrl(cutUrl)
	if err != nil {
		return
	}

	if fullUrl == "www.yandex.ru" {
		fmt.Println("Logic ok")
	}
}

func tryCutter(c *cutter.Cutter) {
	fmt.Println("CutterOk: " + c.Cut("www.yandex.ru"))
}

func main() {
	var err error

	var p *storage.PSQL
	p, err = storage.NewPSQL(storage.WithCacheSize(5),
		storage.WithPostgresCngPath("/home/ilia/Desktop/temki/gonya/UrlCut/configs/postgres.json"))
	if err != nil {
		log.Fatalf("Ошибка в создании объекта storage: %v", err)
		return
	}
	// tryStorage(p)

	var c *cutter.Cutter
	c, err = cutter.NewCutter(6)
	if err != nil {
		log.Fatalf("Ошибка в создании объекта cutter: %v", err)
		return
	}
	// tryCutter(c)

	var l *logic.Logic
	l, err = logic.NewLogic(logic.WithStorage(p),
		logic.WithCutter(c))
	if err != nil {
		log.Fatalf("Ошибка в создании объекта logic: %v", err)
		return
	}
	// tryLogic(l)

	var s interfaces.Server
	s, err = server.NewTerminal(l)
	if err != nil {
		log.Fatalf("Ошибка в создании объекта server: %v", err)
		return
	}

	s.Listen()
}
