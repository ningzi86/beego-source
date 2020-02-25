package core

import (
	"inaction/gitchat/controller"
	"net/http"
	"os"
)

func Handler01() {

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		dir, _ := os.Getwd()
		fs := []string{
			dir + "/templates/layout.html",
			dir + "/templates/layout/body.html",
			dir + "/templates/layout/css.html",
			dir + "/templates/layout/scripts.html",
			dir + "/templates/index.html",
		}

		data := make(map[string]interface{})
		data["fs"] = fs
		data["Title"] = "hello"

		controller.ExecuteViewPathTemplate(w, "index.html", `index.html`, data)

	})
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
