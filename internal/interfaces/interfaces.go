package interfaces

/*Запись и чтение из хранилища*/
type Storage interface {
	GetFullUrl(cutUrl string) (fullUrl string, err error)
	StoreCutUrl(cutUrl string, fullUrl string) (err error)
	Close()
}

/*Обработка событий пользователя*/
type UI interface {
	Exec()
}
