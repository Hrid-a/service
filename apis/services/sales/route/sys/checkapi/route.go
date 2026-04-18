package checkapi

import (
	"github.com/Hrid-a/service/apis/services/api/mid"
	"github.com/Hrid-a/service/app/api/authclient"
	"github.com/Hrid-a/service/business/api/auth"
	"github.com/Hrid-a/service/foundation/logger"
	"github.com/Hrid-a/service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger, authClient *authclient.Client) {
	authen := mid.AuthenticateService(log, authClient)
	athAdminOnly := mid.AuthorizeService(log, authClient, auth.RuleAdminOnly)

	app.HandleFuncNoMiddleware("GET /liveness", liveness)
	app.HandleFuncNoMiddleware("GET /readiness", readiness)
	// app.HandleFunc("GET /testerror", testError)
	// app.HandleFunc("GET /testpanic", testPanic)
	app.HandleFunc("GET /testauth", liveness, authen, athAdminOnly)
}
