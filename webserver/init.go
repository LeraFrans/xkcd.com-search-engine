package webserver

import "net/http"

func Init() {
	http.HandleFunc("GET /login", handleGetLogin)
	http.HandleFunc("POST /login", handlePostLogin)
	http.HandleFunc("GET /comics", handleComics)
}
