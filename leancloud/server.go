package leancloud

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

const cloudFunctionTimeout = time.Second * 15

type metadataResponse struct {
	Result []string `json:"result"`
}

type functionResponse struct {
	Result interface{} `json:"result"`
}

// Handler take all requests related to LeanEngine
func Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := strings.Split(r.RequestURI, "/")
		corsHandler(w, r)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if strings.HasPrefix(r.RequestURI, "/1.1/functions/") || strings.HasPrefix(r.RequestURI, "/1/functions/") {
			if strings.Compare(r.RequestURI, "/1.1/functions/_ops/metadatas") == 0 || strings.Compare(r.RequestURI, "/1/functions/_ops/metadatas") == 0 {
				metadataHandler(w, r)
			} else {
				if uri[3] != "" {
					if len(uri) == 5 {
						classHookHandler(w, r, uri[3], uri[4])
					} else {
						functionHandler(w, r, uri[3], false)
					}
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		} else if strings.HasPrefix(r.RequestURI, "/1.1/call/") || strings.HasPrefix(r.RequestURI, "/1/call/") {
			if functions[uri[3]] != nil {
				functionHandler(w, r, uri[3], true)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else if r.RequestURI == "/__engine/1/ping" || r.RequestURI == "/__engine/1.1/ping" {
			healthCheckHandler(w, r)
		}
	})
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("origin"))

	if r.Method == "OPTIONS" {
		w.Header().Add("Access-Control-Max-Age", "86400")
		w.Header().Add("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", `Content-Type,X-AVOSCloud-Application-Id,X-AVOSCloud-Application-Key,X-AVOSCloud-Application-Production,X-AVOSCloud-Client-Version,X-AVOSCloud-Request-Sign,X-AVOSCloud-Session-Token,X-AVOSCloud-Super-Key,X-LC-Hook-Key,X-LC-Id,X-LC-Key,X-LC-Prod,X-LC-Session,X-LC-Sign,X-LC-UA,X-Requested-With,X-Uluru-Application-Id,X-Uluru-Application-Key,X-Uluru-Application-Production,X-Uluru-Client-Version,X-Uluru-Session-Token`)
	}
}

func metadataHandler(w http.ResponseWriter, r *http.Request) {
	if !validateMasterKey(r) {
		writeCloudError(w, r, CloudError{
			Code:       http.StatusUnauthorized,
			Message:    fmt.Sprintf("Master Key check failed, request from %s", r.RemoteAddr),
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	meta, err := generateMetadata()
	if err != nil {
		writeCloudError(w, r, CloudError{
			Code:       1,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			callStack:  debug.Stack(),
		})
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Write(meta)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(map[string]string{
		"runtime": runtime.Version(),
		"version": Version,
	})
	if err != nil {
		writeCloudError(w, r, CloudError{
			Code:       1,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			callStack:  debug.Stack(),
		})
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}

func functionHandler(w http.ResponseWriter, r *http.Request, name string, rpc bool) {
	if functions[name] == nil {
		writeCloudError(w, r, CloudError{
			Code:       1,
			Message:    fmt.Sprintf("No such cloud function %s", name),
			StatusCode: http.StatusNotFound,
		})
		return
	}

	if functions[name].defineOption["hook"] == true {
		if !validateHookKey(r) {
			writeCloudError(w, r, CloudError{
				Code:       http.StatusUnauthorized,
				Message:    fmt.Sprintf("Hook key check failed, request from %s", r.RemoteAddr),
				StatusCode: http.StatusUnauthorized,
			})
			return
		}
	}

	if functions[name].defineOption["internal"] == true {
		if !validateMasterKey(r) {
			if !validateHookKey(r) {
				master, pass := validateSignature(r)
				if !master || !pass {
					writeCloudError(w, r, CloudError{
						Code:       http.StatusUnauthorized,
						Message:    fmt.Sprintf("Internal cloud function, request from %s", r.RemoteAddr),
						StatusCode: http.StatusUnauthorized,
					})
					return
				}
			}
		}
	}

	if !validateAppKey(r) {
		if !validateMasterKey(r) {
			_, pass := validateSignature(r)
			if !pass {
				writeCloudError(w, r, CloudError{
					Code:       http.StatusUnauthorized,
					Message:    fmt.Sprintf("App key check failed, request from %s", r.RemoteAddr),
					StatusCode: http.StatusUnauthorized,
				})
				return
			}
		}
	}

	request, err := constructRequest(r, name, rpc)
	if err != nil {
		writeCloudError(w, r, CloudError{
			Code:       1,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			callStack:  debug.Stack(),
		})
		return
	}

	ret, err := executeTimeout(request, name, cloudFunctionTimeout)
	if err != nil {
		writeCloudError(w, r, err)
		return
	}
	var resp functionResponse
	if rpc {
		resp.Result = encode(ret, true)
	} else {
		resp.Result = ret
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		writeCloudError(w, r, CloudError{
			Code:       1,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			callStack:  debug.Stack(),
		})
		return
	}

	w.Header().Add("Contetn-Type", "application/json; charset=UTF-8")
	w.Write(respJSON)
}

func classHookHandler(w http.ResponseWriter, r *http.Request, class, hook string) {
	if !validateHookKey(r) {
		writeCloudError(w, r, CloudError{
			Code:       http.StatusUnauthorized,
			Message:    fmt.Sprintf("Hook key check failed, request from %s", r.RemoteAddr),
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	name := fmt.Sprint(classHookmap[hook], class)

	request, err := constructRequest(r, name, false)
	if err != nil {
		writeCloudError(w, r, CloudError{
			Code:       1,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			callStack:  debug.Stack(),
		})
		return
	}

	ret, err := executeTimeout(request, name, cloudFunctionTimeout)

	if err != nil {
		writeCloudError(w, r, err)
		return
	}

	var resp map[string]interface{}
	if hook == "beforeSave" {
		resp = encodeObject(ret, false, false)
	} else {
		resp = map[string]interface{}{
			"result": "ok",
		}
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		writeCloudError(w, r, CloudError{
			Code:       1,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			callStack:  debug.Stack(),
		})
		return
	}

	w.Header().Add("Contetn-Type", "application/json; charset=UTF-8")
	w.Write(respJSON)
}

func executeTimeout(r *FunctionRequest, name string, timeout time.Duration) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var ret interface{}
	var err error
	ch := make(chan bool, 0)
	go func() {
		defer func() {
			if ierr := recover(); ierr != nil {
				err = CloudError{
					Code:       1,
					Message:    fmt.Sprint(ierr),
					StatusCode: http.StatusInternalServerError,
					callStack:  debug.Stack(),
				}
				ch <- true
			}
		}()
		ret, err = functions[name].call(r)
		ch <- true
	}()

	select {
	case <-ch:
		return ret, err
	case <-ctx.Done():
		return nil, CloudError{
			Code:       124,
			Message:    fmt.Sprintf("LeanEngine: /1.1/functions/%s : function timeout (15000ms)", name),
			StatusCode: http.StatusServiceUnavailable,
		}
	}
}

func unmarshalBody(r *http.Request) (interface{}, error) {
	body := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)

	if err == io.EOF {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	return body, nil
}

func constructRequest(r *http.Request, name string, rpc bool) (*FunctionRequest, error) {
	request := new(FunctionRequest)
	request.Meta = map[string]string{
		"remoteAddr": r.RemoteAddr,
	}
	var sessionToken string
	if r.Header.Get("X-LC-Session") != "" {
		sessionToken = r.Header.Get("X-LC-Session")
	} else if r.Header.Get("x-uluru-session-token") != "" {
		sessionToken = r.Header.Get("x-uluru-session-token")
	} else if r.Header.Get("x-avoscloud-session-token") != "" {
		sessionToken = r.Header.Get("x-avoscloud-session-token")
	}

	if functions[name].defineOption["fetchUser"] == true && sessionToken != "" {
		user, err := client.Users.Become(sessionToken)
		if err != nil {
			return nil, err
		}
		request.CurrentUser = user
		request.SessionToken = sessionToken
	}

	if r.Body == nil {
		request.Params = nil
		return request, nil
	}

	params, err := unmarshalBody(r)
	if err != nil {
		return nil, err
	}

	if rpc {
		decodedParams, err := decode(params)
		if err != nil {
			return nil, err
		}

		request.Params = decodedParams
	} else {
		request.Params = params
	}

	return request, nil
}

func generateMetadata() ([]byte, error) {
	meta := metadataResponse{
		Result: []string{},
	}

	for k := range functions {
		meta.Result = append(meta.Result, k)
	}
	return json.Marshal(meta)
}

func validateAppID(r *http.Request) bool {
	if r.Header.Get("X-LC-Id") != "" {
		if os.Getenv("LEANCLOUD_APP_ID") != r.Header.Get("X-LC-Id") {
			return false
		}
	} else if r.Header.Get("x-avoscloud-application-id") != "" {
		if os.Getenv("LEANCLOUD_APP_ID") != r.Header.Get("x-avoscloud-application-id") {
			return false
		}
	} else if r.Header.Get("x-uluru-application-id") != "" {
		if os.Getenv("LEANCLOUD_APP_ID") != r.Header.Get("x-uluru-application-id") {
			return false
		}
	}

	return true
}

func validateAppKey(r *http.Request) bool {
	if !validateAppID(r) {
		return false
	}

	if r.Header.Get("X-LC-Key") != "" {
		if os.Getenv("LEANCLOUD_APP_KEY") != r.Header.Get("X-LC-Key") {
			return false
		}
	} else if r.Header.Get("x-avoscloud-application-key") != "" {
		if os.Getenv("LEANCLOUD_APP_KEY") != r.Header.Get("x-avoscloud-application-key") {
			return false
		}
	} else if r.Header.Get("x-uluru-application-key") != "" {
		if os.Getenv("LEANCLOUD_APP_KEY") != r.Header.Get("x-uluru-application-key") {
			return false
		}
	}
	return true
}

func validateMasterKey(r *http.Request) bool {
	if !validateAppID(r) {
		return false
	}

	if r.Header.Get("X-LC-Key") != "" {
		if strings.TrimSuffix(r.Header.Get("X-LC-Key"), ",master") != os.Getenv("LEANCLOUD_APP_MASTER_KEY") {
			return false
		}
	} else if r.Header.Get("x-avoscloud-master-key") != "" {
		if r.Header.Get("x-avoscloud-master-key") != os.Getenv("LEANCLOUD_APP_MASTER_KEY") {
			return false
		}
	} else if r.Header.Get("x-uluru-master-key") != "" {
		if r.Header.Get("x-uluru-master-key") != os.Getenv("LEANCLOUD_APP_MASTER_KEY") {
			return false
		}
	}
	return true
}

func validateHookKey(r *http.Request) bool {
	if !validateAppID(r) {
		return false
	}

	if os.Getenv("LEANCLOUD_APP_HOOK_KEY") != r.Header.Get("X-LC-Hook-Key") {
		return false
	}
	return true
}

func validateSignature(r *http.Request) (bool, bool) {
	var master, pass bool
	if !validateAppID(r) {
		return master, pass
	}

	var sign string
	if r.Header.Get("X-LC-Sign") != "" {
		sign = r.Header.Get("X-LC-Sign")
	} else if r.Header.Get("x-avoscloud-request-sign") != "" {
		sign = r.Header.Get("x-avoscloud-request-sign")
	}

	if sign == "" {
		return master, pass
	}
	signSlice := strings.Split(sign, ",")
	var hash [16]byte
	if len(signSlice) == 3 && signSlice[2] == "master" {
		hash = md5.Sum([]byte(fmt.Sprint(signSlice[1], os.Getenv("LEANCLOUD_APP_MASTER_KEY"))))
		master = true
	} else {
		hash = md5.Sum([]byte(fmt.Sprint(signSlice[1], os.Getenv("LEANCLOUD_APP_KEY"))))
	}
	if signSlice[0] == fmt.Sprintf("%x", hash) {
		pass = true
	}
	return master, pass
}
