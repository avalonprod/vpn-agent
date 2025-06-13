package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	defaultAllowedIPs = "0.0.0.0/0, ::/0"
)

var (
	usedIPs = map[string]bool{}
	nextIP  = 2
)

func init() {
	mustLoad()

	restoreUsedIPs()
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/register-peer", RegisterNewPeerHandler).Methods(http.MethodPost)
	r.HandleFunc("/peer", DeletePeerHandler).Methods(http.MethodDelete)
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PONG"))
	}).Methods(http.MethodGet)

	http.ListenAndServe(fmt.Sprintf(":%s", ServerConfig.Port), r)
}
