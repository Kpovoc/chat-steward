package web

import "net/http"

func Start(webDir string) {
	http.Handle(
		"/css/",
		http.StripPrefix(
			"/css/", http.FileServer(http.Dir(webDir + "/css"))))
	http.ListenAndServe(":8080", nil)
}