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
}

func (h *HTTP) Listen() {
	http.ListenAndServe("localhost:8090", nil)
}

func (h *HTTP) init() {
	http.HandleFunc("/cut/{fullUrl}", func(w http.ResponseWriter, r *http.Request) {
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

		fmt.Fprint(w, cutUrl)
	})

	http.HandleFunc("/redirect/{cutUrl}", func(w http.ResponseWriter, r *http.Request) {
		var err error

		cutUrl := r.PathValue("cutUrl")
		if len(cutUrl) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.mutex.Lock()
		err = h.logic.Redirect(cutUrl)
		h.mutex.Unlock()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
