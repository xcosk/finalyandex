package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"finalyandex/pkg/api"
)

func Run(password string) error {
	port := 7540
	if env := os.Getenv("TODO_PORT"); env != "" {
		if p, err := strconv.Atoi(env); err == nil && p > 0 {
			port = p
		}
	}

	api.Init(password)
	http.Handle("/", http.FileServer(http.Dir("web")))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
