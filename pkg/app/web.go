package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/pior/runnable"
)

type Web struct {
	cfg    *config.Config
	router *mux.Router
}

func NewWeb(cfg *config.Config) runnable.Runnable {
	router := mux.NewRouter().StrictSlash(true)

	web := Web{cfg, router}

	router.HandleFunc("/ping", web.Ping).Methods(http.MethodGet)

	return &web
}

func (web *Web) Run(ctx context.Context) error {
	hostname := fmt.Sprintf(":%d", web.cfg.HTTPPort())
	web.cfg.Logger().Info().Int("port", web.cfg.HTTPPort()).Msg("Listening and Serving")
	server := http.Server{Addr: hostname, Handler: web.router}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			web.cfg.Logger().Err(err)
		}
	}()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (web *Web) Ping(w http.ResponseWriter, req *http.Request) {
	web.cfg.Logger().Info().Str("uri", req.RequestURI).Str("remote", req.RemoteAddr).Msg("Request received")
	fmt.Fprintf(w, "pong")
}
