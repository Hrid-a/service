package mid

import (
	"context"
	"net/http"

	"github.com/Hrid-a/service/app/api/authclient"
	"github.com/Hrid-a/service/app/api/mid"
	"github.com/Hrid-a/service/foundation/logger"
	"github.com/Hrid-a/service/foundation/web"
)

// AuthorizeService executes the authorize middleware functionality.
func AuthorizeService(log *logger.Logger, client *authclient.Client, rule string) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.AuthorizeService(ctx, log, client, rule, hdl)
		}

		return h
	}

	return m
}
