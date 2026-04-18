package checkapi

import (
	"github.com/Hrid-a/service/business/api/auth"
	"github.com/Hrid-a/service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, a *auth.Auth) {
	app.HandleFuncNoMiddleware("GET /liveness", liveness)
	app.HandleFuncNoMiddleware("GET /readiness", readiness)
}
