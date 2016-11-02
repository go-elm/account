package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/go-elm/account/auth"
	authsvc "github.com/go-elm/account/auth/service"
	kitlog "github.com/go-kit/kit/log"
)

func main() {
	ctx := context.Background()
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	logger.Log("msg", "hello")
	defer logger.Log("msg", "goodbye")

	var session auth.Session
	{
		session = authsvc.New(nil)
	}

	var authEndpoints authsvc.Endpoints
	{
		authEndpoints.Login = authsvc.MakeLoginEndpoint(session)
	}
	loginHandler := authsvc.MakeHTTPHandler(ctx, authEndpoints, logger)
	mux := http.NewServeMux()
	mux.Handle("/login", loginHandler)
	mux.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":8000", mux))
}
