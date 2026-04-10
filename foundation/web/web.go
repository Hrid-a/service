package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
	mw       []MidHandler
}

func NewApp(chutdown chan os.Signal, mw ...MidHandler) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: chutdown,
		mw:       mw,
	}
}

// HandleFunc set a handler function
func (a *App) HandleFunc(pattern string, handler Handler, mw ...MidHandler) {

	handler = wrapMid(mw, handler)
	handler = wrapMid(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		// do something before
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}

		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			fmt.Println(err)
		}
		// do something before
	}

	a.ServeMux.HandleFunc(pattern, h)

}
