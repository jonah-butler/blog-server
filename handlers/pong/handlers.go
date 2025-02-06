package handlers

import "net/http"

func RegisterPongHanlder(prefix string, server *http.ServeMux) {

	server.HandleFunc(prefix, pong)

}

func pong(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong"))
}
