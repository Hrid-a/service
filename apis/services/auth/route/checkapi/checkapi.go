// Package checkapi maintains the web based api for system access.
package checkapi

import (
	"context"
	"net/http"

	"github.com/Hrid-a/service/foundation/web"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Response(ctx, w, status, http.StatusOK)
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Response(ctx, w, status, http.StatusOK)
}
