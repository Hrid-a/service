package mid

import (
	"context"

	"github.com/Hrid-a/service/app/api/errs"
	"github.com/Hrid-a/service/foundation/logger"
)

func Error(ctx context.Context, log *logger.Logger, handler Handler) error {

	err := handler(ctx)

	if err == nil {
		return nil
	}

	log.Error(ctx, "Message", err.Error())

	if errs.IsError(err) {
		return errs.GetError(err)
	}

	return errs.Newf(errs.Unknown, errs.Unknown.String())
}
