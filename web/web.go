package web

import "net/http"

func Start() {
	http.Handle(
		"/css/",
		http.StripPrefix(
			"/css/", http.FileServer(http.Dir("web/css"))))
	http.ListenAndServe(":8080", nil)
}