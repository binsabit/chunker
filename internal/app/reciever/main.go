package reciever

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const BufferSize = 10485760

func Start() {
	r := httprouter.New()
	r.HandlerFunc(http.MethodPost, "/recieve", Recieve)
	r.HandlerFunc(http.MethodPost, "/init", InitUpload)
	server := &http.Server{
		Addr:    ":5000",
		Handler: r,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

///LOGIC

//client initiates upload

//api responds with upload_id
//client using this uplaad id does post request to api
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
