package logic

import (
	"fmt"
	"net/http"
)

func ReturnException(w http.ResponseWriter, err error, functionName string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("%s: %v", functionName, err)))
}
