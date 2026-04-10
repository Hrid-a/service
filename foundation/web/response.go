package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func Response(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	bytData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")

	if _, err := w.Write(bytData); err != nil {
		return err
	}

	return nil
}
