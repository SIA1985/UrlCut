package server

import (
	"UrlCut/internal/logic"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type HTTP struct {
	logic *logic.Logic
	mutex sync.Mutex

	addr string
	mux  *http.ServeMux
}

func (h *HTTP) Listen() {
	//todo: Context
	http.ListenAndServe(h.addr, h.mux)
}

func (h *HTTP) createRoutes() {
	h.mux.HandleFunc("/cut/{fullUrl...}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		fullUrl := r.PathValue("fullUrl")

		_, err = url.ParseRequestURI(fullUrl)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		var cutUrl string

		h.mutex.Lock()
		cutUrl, err = h.logic.CutUrl(fullUrl)
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

		cutUrl := r.PathValue("cutUrl")
		if len(cutUrl) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var fullUrl string

		h.mutex.Lock()
		fullUrl, err = h.logic.GetFullUrl(cutUrl)
		h.mutex.Unlock()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		http.Redirect(w, r, fullUrl, http.StatusSeeOther)
	})
}
