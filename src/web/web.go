package web

import (
	"net/http"
)

func Start(webDir string, websitePort string) {
	if websitePort == "" || len(websitePort) > 5 {
		websitePort = "8080"
	}

	http.Handle(
		"/css/",
		http.StripPrefix(
			"/css/", http.FileServer(http.Dir(webDir + "/css"))))
	http.ListenAndServe(":" + websitePort, nil)
}