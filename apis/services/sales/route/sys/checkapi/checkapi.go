package checkapi

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/Hrid-a/service/app/api/errs"
	"github.com/Hrid-a/service/foundation/web"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	status := struct {
		Status string
	}{
		Status: "Ok liveness",
	}

	return web.Response(ctx, w, status, http.StatusOK)
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	status := struct {
		Status string
	}{
		Status: "Ok readiness",
	}

	return web.Response(ctx, w, status, http.StatusOK)
}

func testerr(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100); n%2 == 0 {
		return errs.Newf(errs.FailedPrecondition, "this message is trused")
	}

	status := struct {
		Status string
	}{
		Status: "Ok readiness",
	}

	return web.Response(ctx, w, status, http.StatusOK)
}
func testpanic(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	panic("We're cooked!!!!")
}
