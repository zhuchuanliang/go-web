package main

import (
	"git.woa.com/alanclzhu/go-web/core"
	"git.woa.com/alanclzhu/go-web/router"
	"net/http"
)

func main() {
	core := core.NewCore()
	router.RegisterRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    "localhost:8181",
	}

	server.ListenAndServe()
}
