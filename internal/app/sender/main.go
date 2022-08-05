package sender

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Start() {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/send", Send)
	router.HandlerFunc(http.MethodPost, "/upload", Upload)
	router.HandlerFunc(http.MethodGet, "/upload", ShowUplaodPage)
	router.ServeFiles("/static/*filepath", http.Dir("public"))
	server := &http.Server{
		Addr:    ":4000",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
