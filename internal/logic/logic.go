package logic

import (
	"UrlCut/internal/cutter"
	"UrlCut/internal/interfaces"
	"UrlCut/internal/storage"
	"context"
	"net/url"
)

type Logic struct {
	storage interfaces.Storage
	cutter  *cutter.Cutter
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
	}

	for _, opt := range opts {
		opt(l)
	}

	return
}

func (l *Logic) CutUrl(ctx context.Context, fullUrl string) (cutUrl string, err error) {
	_, err = url.ParseRequestURI(fullUrl)
	if err != nil {
		return
	}

	cutUrl, err = l.cutter.Cut(fullUrl)
	if err != nil {
		return
	}

	strCtx, _ := context.WithCancel(ctx)
	err = l.storage.StoreCutUrl(strCtx, cutUrl, fullUrl)
	if err != nil {
		return
	}

	return
}

func (l *Logic) GetFullUrl(ctx context.Context, cutUrl string) (fullUrl string, err error) {
	strCtx, _ := context.WithCancel(ctx)
	fullUrl, err = l.storage.GetFullUrl(strCtx, cutUrl)
	return
}
