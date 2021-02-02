package leancloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

// CloudError contains user-defined error
type CloudError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	panic   bool
}

func (err CloudError) Error() string {
	return fmt.Sprintf("CloudError: code: %d, Message: %s\n", err.Code, err.Message)
}

func cloudError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Add("Contetn-Type", "application/json; charset=UTF-8")
	if cloudErr, ok := err.(CloudError); ok {
		cloudErrJSON, err := json.Marshal(cloudErr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%s: %s\n", err.Error(), cloudErr.Error())))
			return
		}
		if cloudErr.panic {
			w.WriteHeader(http.StatusInternalServerError)
			debug.PrintStack()
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write(cloudErrJSON)
		return
	}

	cloudError(w, r, CloudError{
		Code:    1,
		Message: err.Error(),
	})
}

func validationError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(err.Error()))
}

func internalError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Fprintln(os.Stderr, err)
	debug.PrintStack()
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
