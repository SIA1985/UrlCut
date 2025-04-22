package interfaces


/*Запись и чтение из хранилища*/
type Storage interface {
	Init()
	GetFullUrl(cutUrl string) (fullUrl string, err error)
	StoreCutUrl(cutUrl string, fullUrl string) (err error)
}