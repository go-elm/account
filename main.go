package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/go-elm/account/auth"
	authsvc "github.com/go-elm/account/auth/service"
	"github.com/go-elm/account/user"
	"github.com/go-elm/account/user/datastore/inmem"
	kitlog "github.com/go-kit/kit/log"
)

func main() {
	ctx := context.Background()
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	logger.Log("msg", "hello")
	defer logger.Log("msg", "goodbye")

	userStore := inmem.New()
	createFakeAdmin(userStore)

	var session auth.Session
	{
		session = authsvc.New(userStore)
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

func createFakeAdmin(ds user.Datastore) {
	user := user.User{
		Username: "groob",
		Password: []byte("secret"),
	}
	_, err := ds.Create(user)
	if err != nil {
		panic(err)
	}
}
