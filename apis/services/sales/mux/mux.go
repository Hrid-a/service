// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"os"

	"github.com/Hrid-a/service/apis/services/api/mid"
	"github.com/Hrid-a/service/apis/services/sales/route/sys/checkapi"
	"github.com/Hrid-a/service/app/api/authclient"
	"github.com/Hrid-a/service/foundation/logger"
	"github.com/Hrid-a/service/foundation/web"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(log *logger.Logger, authClient *authclient.Client, shutdown chan os.Signal) *web.App {
	app := web.NewApp(shutdown, mid.Logger(log), mid.Error(log), mid.Metrics(), mid.Panic())

	checkapi.Routes(app, log, authClient)

	return app
}
