package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"collectd.org/api"
	"collectd.org/config"
	"collectd.org/plugin"
	"go.uber.org/multierr"
)

const pluginName = "restapi"

type restapi struct {
	srv *http.Server
}

func init() {
	ra := &restapi{}

	plugin.RegisterConfig(pluginName, ra)
	plugin.RegisterShutdown(pluginName, ra)
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

type srvConfig struct {
	Args     string // unused
	Addr     string
	Port     string
	CertFile string
	KeyFile  string
}

func (ra *restapi) Configure(_ context.Context, rawConfig config.Block) error {
	cfg := srvConfig{}
	if err := rawConfig.Unmarshal(&cfg); err != nil {
		return err
	}

	if (cfg.CertFile == "") != (cfg.KeyFile == "") {
		return fmt.Errorf("CertFile=%q, KeyFile=%q; need both for TLS setup",
			cfg.CertFile, cfg.KeyFile)
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
		if cfg.CertFile != "" {
			cfg.Port = "8443"
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/valueList", valueListHandler)

	ra.srv = &http.Server{
		Addr:    net.JoinHostPort(cfg.Addr, cfg.Port),
		Handler: mux,
	}

	go func() {
		var err error
		if cfg.CertFile != "" {
			err = ra.srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = ra.srv.ListenAndServe()
		}
		if !errors.Is(err, http.ErrServerClosed) {
			plugin.Errorf("%s plugin: ListenAndServe(): %v", pluginName, err)
		}
	}()

	return nil
}

func (ra *restapi) Shutdown(ctx context.Context) error {
	if ra == nil || ra.srv == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return ra.srv.Shutdown(ctx)
}

func main() {} // ignored
