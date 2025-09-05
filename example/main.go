package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/lemoba/chix"
	"github.com/lemoba/chix/example/internal/handler/user"

	xmiddleware "github.com/lemoba/chix/pkg/middleware"
)

func main() {
	r := chix.NewRouter()

	r.Use(middleware.Logger, xmiddleware.Recovery)

	r.Get("/health", user.SayHello)

	fmt.Println("Server started at http://127.0.0.1:3000")
	http.ListenAndServe(":3000", r)
}
