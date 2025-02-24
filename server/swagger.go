package server

import (
	"fmt"
	"net/http"

	"github.com/adminoid/cosmos-sdk/client"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	_ "github.com/adminoid/cosmos-sdk/client/docs/statik" // we use this to help with sderving the docs
)

// RegisterSwaggerAPI provides a common function which registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router, swaggerEnabled bool) error {
	if !swaggerEnabled {
		return nil
	}

	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("failed to create filesystem: %w", err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))

	return nil
}
