package authapi

import (
	"github.com/Hrid-a/service/apis/services/api/mid"
	"github.com/Hrid-a/service/business/api/auth"
	"github.com/Hrid-a/service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, a *auth.Auth) {
	authen := mid.AuthenticateLocal(a)

	api := newAPI(a)
	app.HandleFunc("GET /auth/token/{kid}", api.token, authen)
	app.HandleFunc("GET /auth/authenticate", api.authenticate, authen)
	app.HandleFunc("POST /auth/authorize", api.authorize)
}
