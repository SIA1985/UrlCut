package main

import (
	"UrlCut/internal/cutter"
	"UrlCut/internal/interfaces"
	"UrlCut/internal/logic"
	"UrlCut/internal/server"
	"UrlCut/internal/storage"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var err error

	mode := flag.String("mode", "http", "Тип взаимодействия: http, terminal")
	storageType := flag.String("storage_type", "postgres", "Тип хранилища: postgres")
	storageConfigPath := flag.String("storage_config_path", "/home/ilia/Desktop/temki/gonya/UrlCut/configs/postgres.json", "Путь кофигурационному файлу для хранилища")
	cacheSize := flag.Int("cache_size", 100, "Размер кэша недавних ссылок")
	cutUrlLen := flag.Int("cut_url_len", 6, "")

	httpAddr := flag.String("http_addr", "localhost:8090", "Адрес http-сервера")

	flag.Parse()

	var p *storage.PSQL
	switch *storageType {
	case "postgres":
		p, err = storage.NewPSQL(storage.WithCacheSize(*cacheSize),
			storage.WithPostgresCngPath(*storageConfigPath))

	}
	if err != nil {
		log.Fatalf("Ошибка в создании объекта storage: %v", err)
	}

	var c *cutter.Cutter
	c, err = cutter.NewCutter(*cutUrlLen)
	if err != nil {
		log.Fatalf("Ошибка в создании объекта cutter: %v", err)
		return
	}

	var l *logic.Logic
	l, err = logic.NewLogic(logic.WithStorage(p),
		logic.WithCutter(c))
	if err != nil {
		log.Fatalf("Ошибка в создании объекта logic: %v", err)
		return
	}

	var s interfaces.Server
	switch *mode {
	case "http":
		s, err = server.NewHTTP(l, *httpAddr)

	case "terminal":
		s, err = server.NewTerminal(l)

	}
	if err != nil {
		log.Fatalf("Ошибка в создании объекта server: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	server.SetContext(ctx)

	go s.Listen()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	cancel()
}
