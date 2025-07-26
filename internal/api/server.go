package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/risingwavelabs/eris"

	v1 "template/internal/api/v1"
	"template/internal/config"
)

type Server struct {
	server *http.Server
}

func (Server) Name() string { return "API Server" }

func (svr *Server) Init(ctx context.Context) error {
	v1Swagger, err := v1.GetSwagger()
	if err != nil {
		return eris.Wrap(err, "failed to load v1 swagger")
	}

	v1BasePath, err := v1Swagger.Servers.BasePath()
	if err != nil {
		return eris.New("failed to load base path from v1 swagger")
	}

	v1Router := http.NewServeMux()
	v1Router.Handle(v1BasePath+"/", http.StripPrefix(
		v1BasePath,
		v1.HandlerFromMux(v1.Controller{}, http.NewServeMux()),
	))

	svr.server = &http.Server{
		Addr:    ":" + strconv.Itoa(int(config.C.APIPort)),
		Handler: v1Router,
	}

	return nil
}

func (svr *Server) Run(ctx context.Context) error {
	fmt.Printf("%s is listening on port %s\n", svr.Name(), svr.server.Addr)

	svr.server.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}

	err := svr.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return eris.Wrapf(err, "%s stopped", svr.Name())
	}

	return nil
}

func (svr *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := svr.server.Shutdown(ctx)
	if err != nil {
		return eris.Wrapf(err, "failed to shut down %s", svr.Name())
	}

	return nil
}
