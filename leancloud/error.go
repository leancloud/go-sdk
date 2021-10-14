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
	Message    string `json:"error"`
	StatusCode int    `json:"-"`
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
			StatusCode: http.StatusBadRequest,
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

	if cloudErr.StatusCode == 0 {
		cloudErr.StatusCode = http.StatusBadRequest
	}
	w.WriteHeader(cloudErr.StatusCode)
	w.Write(cloudErrJSON)
}
