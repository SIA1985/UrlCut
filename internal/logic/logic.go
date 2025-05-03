package logic

import (
	"UrlCut/internal/cutter"
	"UrlCut/internal/interfaces"
	"UrlCut/internal/storage"
	"fmt"
	"os/exec"
	"runtime"
)

func openBrowser(fullUrl string) (err error) {
	switch os := runtime.GOOS; os {
	case "linux":
		err = exec.Command("x-www-browser", fullUrl).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", fullUrl).Start()
	case "darwin":
		err = exec.Command("open", fullUrl).Start()
	default:
		err = fmt.Errorf("Не поддерживаемая платформа!")
	}

	return
}

type Logic struct {
	storage interfaces.Storage
	cutter  *cutter.Cutter

	redirectFunc func(fullUrl string) (err error)
}

type LogicOption func(*Logic)

func WithStorage(storage interfaces.Storage) LogicOption {
	return func(l *Logic) {
		if l.storage != nil {
			l.storage.Close()
		}
		l.storage = storage
	}
}

func WithCutter(cutter *cutter.Cutter) LogicOption {
	return func(l *Logic) {
		l.cutter = cutter
	}
}

func WithRedirectFunc(redirectFunc func(fullUrl string) (err error)) LogicOption {
	return func(l *Logic) {
		l.redirectFunc = redirectFunc
	}
}

func NewLogic(opts ...LogicOption) (l *Logic, err error) {
	var p *storage.PSQL
	p, err = storage.NewPSQL() //todo: 2 раз открывается подключение
	if err != nil {
		return
	}

	var c *cutter.Cutter
	c, err = cutter.NewCutter(6)
	if err != nil {
		return
	}

	l = &Logic{
		storage: p,
		cutter:  c,

		redirectFunc: openBrowser,
	}

	for _, opt := range opts {
		opt(l)
	}

	return
}

func (l *Logic) CutUrl(fullUrl string) (cutUrl string, err error) {
	cutUrl = l.cutter.Cut(fullUrl)

	err = l.storage.StoreCutUrl(cutUrl, fullUrl)
	if err != nil {
		return
	}

	return
}

func (l *Logic) Redirect(cutUrl string) (err error) {
	var fullUrl string
	fullUrl, err = l.storage.GetFullUrl(cutUrl)
	if err != nil {
		return
	}

	err = l.redirectFunc(fullUrl)
	if err != nil {
		return
	}

	return
}
