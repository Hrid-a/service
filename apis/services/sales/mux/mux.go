package mux

import (
	"os"

	"github.com/Hrid-a/service/apis/services/api/mid"
	"github.com/Hrid-a/service/apis/services/sales/route/sys/checkapi"
	"github.com/Hrid-a/service/foundation/logger"
	"github.com/Hrid-a/service/foundation/web"
)

func WebAPI(log *logger.Logger, ch chan os.Signal) *web.App {

	app := web.NewApp(ch, mid.Logger(log), mid.Error(log), mid.Metrics(), mid.Panic())

	checkapi.Routes(app)
	return app
}
