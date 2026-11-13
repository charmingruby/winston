package winstom

import (
	"net/http"
)

func (w *Winstom) Handle(rw http.ResponseWriter, r *http.Request) {
	w.HandleHTTP(rw, r)
}
