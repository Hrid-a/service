package checkapi

import (
	"github.com/Hrid-a/service/foundation/web"
)

func Routes(mux *web.App) {

	mux.HandleFuncNoMiddleware("GET /liveness", liveness)
	mux.HandleFuncNoMiddleware("GET /readiness", readiness)
	mux.HandleFunc("GET /testerr", testerr)
	mux.HandleFunc("GET /testpanic", testpanic)

}
