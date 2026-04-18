// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"os"

	"github.com/Hrid-a/service/apis/services/api/mid"
	"github.com/Hrid-a/service/apis/services/auth/route/authapi"
	"github.com/Hrid-a/service/apis/services/auth/route/checkapi"
	"github.com/Hrid-a/service/business/api/auth"
	"github.com/Hrid-a/service/foundation/logger"
	"github.com/Hrid-a/service/foundation/web"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(log *logger.Logger, auth *auth.Auth, shutdown chan os.Signal) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Error(log), mid.Metrics(), mid.Panic())

	checkapi.Routes(app, auth)
	authapi.Routes(app, auth)

	return app
}
