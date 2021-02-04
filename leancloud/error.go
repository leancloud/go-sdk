package leancloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// CloudError contains user-defined error
type CloudError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	statusCode int
	callStack  []byte
}

func (err CloudError) Error() string {
	return fmt.Sprintf("CloudError: code: %d, Message: %s\n", err.Code, err.Message)
}

func writeCloudError(w http.ResponseWriter, r *http.Request, err error) {
	cloudErr, ok := err.(CloudError)
	if !ok {
		writeCloudError(w, r, CloudError{
			Code:       1,
			Message:    err.Error(),
			statusCode: http.StatusBadRequest,
		})
		return
	}

	cloudErrJSON, err := json.Marshal(cloudErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s: %s\n", err.Error(), cloudErr.Error())))
		return
	}

	if len(cloudErr.callStack) != 0 {
		builder := new(strings.Builder) // for better performance when converting []byte to string. see https://pkg.go.dev/strings#Builder
		builder.Write(cloudErr.callStack)
		fmt.Fprintln(os.Stderr, builder.String())
	}

	w.Header().Add("Contetn-Type", "application/json; charset=UTF-8")
	w.WriteHeader(cloudErr.statusCode)
	w.Write(cloudErrJSON)
}
