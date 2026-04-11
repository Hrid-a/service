package mid

import (
	"context"
	"net/http"

	"github.com/Hrid-a/service/app/api/mid"
	"github.com/Hrid-a/service/foundation/web"
)

func Panic() web.MidHandler {

	return func(h web.Handler) web.Handler {

		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hl := func(ctx context.Context) error {
				return h(ctx, w, r)
			}

			return mid.Panic(ctx, hl)
		}
	}
}
