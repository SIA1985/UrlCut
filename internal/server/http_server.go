package server

import (
	"UrlCut/internal/logic"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type HTTP struct {
	logic *logic.Logic
	mutex sync.Mutex

	addr string
	mux  *http.ServeMux
	srv  *http.Server
}

func (h *HTTP) Listen() {
	h.srv.ListenAndServe()
}

func (h *HTTP) createRoutes() {
	h.mux.HandleFunc("/cut/{fullUrl...}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		select {
		case <-srvCtx.Done():
			w.WriteHeader(http.StatusNotExtended)
			return
		case <-r.Context().Done():
			return
		default:
		}

		fullUrl := r.PathValue("fullUrl")

		var cutUrl string

		h.mutex.Lock()
		cutUrl, err = h.logic.CutUrl(srvCtx, fullUrl)
		h.mutex.Unlock()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		fmt.Fprint(w, h.addr+"/"+cutUrl)
	})

	h.mux.HandleFunc("/{cutUrl}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		select {
		case <-srvCtx.Done():
			w.WriteHeader(http.StatusNotExtended)
			return
		case <-r.Context().Done():
			return
		default:
		}

		cutUrl := r.PathValue("cutUrl")
		if len(cutUrl) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var fullUrl string

		h.mutex.Lock()
		fullUrl, err = h.logic.GetFullUrl(srvCtx, cutUrl)
		h.mutex.Unlock()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		http.Redirect(w, r, fullUrl, http.StatusFound)
	})
}
