package server

import (
	"UrlCut/internal/logic"
	"context"
	"fmt"
	"net/http"
)

var srvCtx context.Context

func SetContext(ctx context.Context) {
	srvCtx = ctx
}

func NewHTTP(logic *logic.Logic, addr string) (h *HTTP, err error) {
	if logic == nil {
		err = fmt.Errorf("logic == nil")
		return
	}

	h = &HTTP{
		logic: logic,
		addr:  addr,
		mux:   http.NewServeMux(),
	}

	h.createRoutes()

	h.srv = &http.Server{
		Addr:    addr,
		Handler: h.mux,
	}

	return
}

func NewTerminal(logic *logic.Logic) (t *Terminal, err error) {
	if logic == nil {
		err = fmt.Errorf("logic == nil")
		return
	}

	t = &Terminal{
		logic: logic,
	}

	return
}
