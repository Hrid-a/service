package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
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

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
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
			if validateError(err) {
				a.SignalShutdown()
			}
			return
		}
		// do something before
	}

	a.ServeMux.HandleFunc(pattern, h)

}

// HandleFuncNoMiddleware set a handler function
func (a *App) HandleFuncNoMiddleware(pattern string, handler Handler, mw ...MidHandler) {

	h := func(w http.ResponseWriter, r *http.Request) {
		// do something before
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}

		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			if validateError(err) {
				a.SignalShutdown()
			}
			return
		}
		// do something before
	}

	a.ServeMux.HandleFunc(pattern, h)

}

// validateError validates the error for special conditions that do not
// warrant an actual shutdown by the system.
func validateError(err error) bool {

	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return false
	}

	return true
}
