package httputils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rest/internal/pkg/appLogger"
	"rest/internal/pkg/errorspkg"
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

func WriteError(ctx context.Context, w http.ResponseWriter, logger appLogger.IAppLogger, err error) {
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
	statusCode := determineStatusCode(err)
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("writeError error: %w", err))
		return
	}
}

func determineStatusCode(err error) int {
	var appErr *errorspkg.AppError
	switch {
	case errors.Is(err, errorspkg.ErrWrongOperationType), errors.Is(err, errorspkg.ErrWrongAmount), errors.Is(err, errorspkg.ErrWalletUUIDIsMissed):
		return http.StatusBadRequest
	case errors.As(err, &appErr):
		return appErr.Status
	default:
		return http.StatusInternalServerError
	}
}
