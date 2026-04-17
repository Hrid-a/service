package mid

import (
	"context"
	"net/http"

	"github.com/Hrid-a/service/app/api/errs"
	"github.com/Hrid-a/service/app/api/mid"
	"github.com/Hrid-a/service/foundation/logger"
	"github.com/Hrid-a/service/foundation/web"
)

var codeStatus [17]int

// init maps out the error codes to http status codes.
func init() {
	codeStatus[errs.OK.Value()] = http.StatusOK
	codeStatus[errs.Canceled.Value()] = http.StatusGatewayTimeout
	codeStatus[errs.Unknown.Value()] = http.StatusInternalServerError
	codeStatus[errs.InvalidArgument.Value()] = http.StatusBadRequest
	codeStatus[errs.DeadlineExceeded.Value()] = http.StatusGatewayTimeout
	codeStatus[errs.NotFound.Value()] = http.StatusNotFound
	codeStatus[errs.AlreadyExists.Value()] = http.StatusConflict
	codeStatus[errs.PermissionDenied.Value()] = http.StatusForbidden
	codeStatus[errs.ResourceExhausted.Value()] = http.StatusTooManyRequests
	codeStatus[errs.FailedPrecondition.Value()] = http.StatusBadRequest
	codeStatus[errs.Aborted.Value()] = http.StatusConflict
	codeStatus[errs.OutOfRange.Value()] = http.StatusBadRequest
	codeStatus[errs.Unimplemented.Value()] = http.StatusNotImplemented
	codeStatus[errs.Internal.Value()] = http.StatusInternalServerError
	codeStatus[errs.Unavailable.Value()] = http.StatusServiceUnavailable
	codeStatus[errs.DataLoss.Value()] = http.StatusInternalServerError
	codeStatus[errs.Unauthenticated.Value()] = http.StatusUnauthorized
}

func Error(log *logger.Logger) web.MidHandler {

	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			if err := mid.Error(ctx, log, hdl); err != nil {

				errs := err.(errs.Error)

				if err := web.Response(ctx, w, errs, codeStatus[errs.Code.Value()]); err != nil {
					return err
				}

				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}
	}
}
