package logic

import (
	"fmt"
	"log/slog"
	"net/http"
)

func ReturnException(w http.ResponseWriter, err error, functionName string) {
	slog.Error(fmt.Sprintf("%s: %v", functionName, err))
	w.WriteHeader(http.StatusInternalServerError)
}
