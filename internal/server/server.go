package server

import (
	"UrlCut/internal/logic"
	"fmt"
)

func NewHTTP(logic *logic.Logic, addr string) (h *HTTP, err error) {
	if logic == nil {
		err = fmt.Errorf("logic == nil")
		return
	}

	h = &HTTP{
		logic: logic,
		addr:  addr,
	}

	h.createRoutes()

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
