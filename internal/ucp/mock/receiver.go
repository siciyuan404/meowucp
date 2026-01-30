package mock

import (
	"io"
	"net/http"
)

type Receiver struct {
	FailStatus int
}

func NewReceiver(failStatus int) *Receiver {
	return &Receiver{FailStatus: failStatus}
}

func (r *Receiver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.FailStatus > 0 {
		w.WriteHeader(r.FailStatus)
		return
	}

	_, _ = io.ReadAll(req.Body)
	_ = req.Body.Close()
	writeJSONOK(w)
}

func writeJSONOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
