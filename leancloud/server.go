package leancloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const cloudFunctionTimeout = time.Second * 15
const beforeHookTimeout = time.Second * 10
const generalHookTimeout = time.Second * 3

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
						hookHandler(w, r, uri[3], uri[4])
					} else {
						if functions[realtimeHookmap[uri[3]]] != nil {
							if !hookAuthenticate(r.Header.Get("X-LC-Hook-Key")) {
								errorResponse(w, r, fmt.Errorf("Hook key check failed, request from %s", r.RemoteAddr))
								return
							}
						}
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
	meta, err := generateMetadata()
	if err != nil {
		errorResponse(w, r, err)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Write(meta)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(map[string]string{
		"runtime": "go-1.14",
		"version": "0.1.0",
	})
	if err != nil {
		errorResponse(w, r, err)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}

func functionHandler(w http.ResponseWriter, r *http.Request, name string, rpc bool) {
	if functions[name].defineOption["internal"] == true {
		errorResponse(w, r, fmt.Errorf("Internal cloud function, request from %s", r.RemoteAddr))
		return
	}

	request, err := constructRequest(r, name, rpc)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	ret, err := executeTimeout(request, name, cloudFunctionTimeout)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	resp := functionResponse{
		Result: ret,
	}

	respJSON, err := json.Marshal(resp)
	w.Header().Add("Contetn-Type", "application/json; charset=UTF-8")
	w.Write(respJSON)
}

func hookHandler(w http.ResponseWriter, r *http.Request, class, hook string) {
	if !hookAuthenticate(r.Header.Get("X-LC-Hook-Key")) {
		errorResponse(w, r, fmt.Errorf("Hook key check failed, request from %s", r.RemoteAddr))
		return
	}

	var name string
	if storageHookmap[hook] != "" {
		name = fmt.Sprint(storageHookmap[hook], class)
	} else {
		name = realtimeHookmap[hook]
	}

	request, err := constructRequest(r, name, false)
	if err != nil {
		errorResponse(w, r, err)
		return
	}

	var ret interface{}
	if strings.HasPrefix(hook, "before") {
		ret, err = executeTimeout(request, name, beforeHookTimeout)
	} else {
		ret, err = executeTimeout(request, name, generalHookTimeout)
	}

	if err != nil {
		errorResponse(w, r, err)
		return
	}

	var resp map[string]interface{}
	if hook == "beforeSave" {
		resp = encodeObject(ret, false, false)
	} else if strings.HasPrefix(hook, "onIM") {
		resp = ret.(map[string]interface{})
	} else {
		resp = map[string]interface{}{
			"result": "ok",
		}
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		errorResponse(w, r, err)
		return
	}
	w.Header().Add("Contetn-Type", "application/json; charset=UTF-8")
	w.Write(respJSON)
}

func executeTimeout(r *Request, name string, timeout time.Duration) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var ret interface{}
	var err error
	ch := make(chan bool, 0)
	go func() {
		ret, err = functions[name].call(r)
		ch <- true
	}()

	select {
	case <-ch:
		return ret, err
	case <-ctx.Done():
		return nil, fmt.Errorf("LeanEngine: /1.1/functions/%s : function timeout (15000ms)", name)
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

func constructRequest(r *http.Request, name string, rpc bool) (*Request, error) {
	request := new(Request)
	request.Meta = map[string]string{
		"remoteAddr": r.RemoteAddr,
	}
	sessionToken := r.Header.Get("X-LC-Session")
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

func errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Add("Contetn-Type", "application/json; charset=UTF-8")
	switch err.(type) {
	case *functionError:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
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
