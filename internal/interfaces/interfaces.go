package interfaces

import "context"

/*Запись и чтение из хранилища*/
type Storage interface {
	GetFullUrl(ctx context.Context, cutUrl string) (fullUrl string, err error)
	StoreCutUrl(ctx context.Context, cutUrl string, fullUrl string) (err error)
	Close()
}

/*Обработка событий пользователя*/
type Server interface {
	Listen()
}
