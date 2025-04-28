package logic


import(
	"UrlCut/internal/interfaces"
	"UrlCut/internal/cutter"
)

type Logic struct {
	storage 	interfaces.Storage
	cutter		*cutter.Cutter
}

func NewLogic(s interfaces.Storage, c *cutter.Cutter) (l *Logic) {
	l = &Logic{
		storage: s,
		cutter: c,
	}

	return
}

func (l *Logic) CutUrl(fullUrl string) (cutUrl string, err error) {
	cutUrl = l.cutter.Cut(fullUrl)

	err = l.storage.StoreCutUrl(cutUrl, fullUrl)

	return 
}

func (l *Logic) Redirect(cutUrl string) (err error) {
	var fullUrl string
	fullUrl, err = l.storage.GetFullUrl(cutUrl)

	_ = fullUrl
	/*redirect in browser*/

	return
}