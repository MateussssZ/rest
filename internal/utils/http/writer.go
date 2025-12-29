package httputils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"rest/internal/pkg/appLogger"
)

type errorStruct struct {
	Message string `json:"Message"`
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, logger appLogger.IAppLogger, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	jsonBody, err := json.Marshal(body)
	if err != nil {

		logger.Error(ctx, fmt.Errorf("json.Marshal in response error: %w", err))
	}

	_, err = w.Write(jsonBody)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("w.Write in response error: %w", err))
	}
}

func WriteError(ctx context.Context, w http.ResponseWriter, logger appLogger.IAppLogger, statusCode int, err error) {
	var responseErr errorStruct

	if err != nil {
		responseErr.Message = err.Error()
	}

	response, err := json.Marshal(responseErr)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("writeError error: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("writeError error: %w", err))
		return
	}
}
