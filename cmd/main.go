package main

import(
	"UrlCut/internal/storage"
	"UrlCut/internal/interfaces"

	"fmt"
)

func tryStorage(s interfaces.Storage) {
	s.Init()
	url, _ := s.GetFullUrl("aaaaaa")
	url2, _ := s.GetFullUrl("aaaaaa")

	if url == url2 {
		fmt.Print("Sotrage ok")
		return
	}
}

func main() {
	tryStorage(&storage.PSQL{})
}