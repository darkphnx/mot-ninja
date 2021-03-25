package api

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}
