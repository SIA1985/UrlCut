package main

import(
	"UrlCut/internal/storage"
	"UrlCut/internal/interfaces"
)

func tryStorage(s interfaces.Storage) {
	s.Init()
	s.GetFullUrl("aaaaaa")
}

func main() {
	tryStorage(&storage.PSQL{})
}