package storage


/*SQL-хранилище*/

type SQL struct {
	cache 		*Cache
}

func (s *SQL) Init() {
	s.cache = NewCache(1000)
}

func (s *SQL) GetFullUrl(cutUrl string) (fullUrl string, err error) {

}

func (s *SQL) StoreCutUrl(cutUrl string, fullUrl string) (err error) {

}