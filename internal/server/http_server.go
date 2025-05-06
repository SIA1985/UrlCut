package server

import (
	"UrlCut/internal/logic"
	"fmt"
	"net/http"
	"sync"
)

type HTTP struct {
	logic *logic.Logic
	mutex sync.Mutex

	addr string
}

func (h *HTTP) Listen() {
	http.ListenAndServe(h.addr, nil)
}

func (h *HTTP) createRoutes() {
	http.HandleFunc("/cut/{fullUrl...}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		fullUrl := r.PathValue("fullUrl")
		if len(fullUrl) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var cutUrl string

		h.mutex.Lock()
		cutUrl, err = h.logic.CutUrl(fullUrl)
		h.mutex.Unlock()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, h.addr+"/"+cutUrl)
	})

	http.HandleFunc("/{cutUrl}", func(w http.ResponseWriter, r *http.Request) {
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
			return
		}

		http.Redirect(w, r, fullUrl, http.StatusSeeOther)
	})
}
