package mux

import (
	"os"

	"github.com/Hrid-a/service/apis/services/sales/route/sys/checkapi"
	"github.com/Hrid-a/service/foundation/web"
)

func WebAPI(ch chan os.Signal) *web.App {

	app := web.NewApp(ch)

	checkapi.Routes(app)
	return app
}
