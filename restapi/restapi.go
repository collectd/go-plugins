package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"collectd.org/api"
	"collectd.org/plugin"
	"go.uber.org/multierr"
)

const pluginName = "restapi"

type restapi struct {
	srv *http.Server
}

func init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/valueList", valueListHandler)

	api := restapi{
		srv: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}

	go func() {
		if err := api.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			plugin.Errorf("%s plugin: ListenAndServe(): %v", pluginName, err)
		}
	}()

	plugin.RegisterShutdown(pluginName, api)
}

func valueListHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintln(w, "Only POST is currently supported.")
		return
	}

	var vls []api.ValueList
	if err := json.NewDecoder(req.Body).Decode(&vls); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "parsing JSON failed:", err)
		return
	}

	var errs error
	for _, vl := range vls {
		errs = multierr.Append(errs,
			plugin.Write(req.Context(), &vl))
	}

	if errs != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "plugin.Write():", errs)
		return
	}
}

func (api restapi) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return api.srv.Shutdown(ctx)
}

func main() {} // ignored
