package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type metadataResponse struct {
	Result []string `json:"result"`
}

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := strings.Split(r.RequestURI, "/")
		if strings.HasPrefix(r.RequestURI, "/1.1/functions/") {
			if strings.Compare(r.RequestURI, "/1.1/functions/_ops/metadata") == 0 {
				metadataHandler(w, r)
			} else {
				if functions[uri[3]] != nil {
					functionHandler(w, r, uri[3])
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		} else if strings.HasPrefix(r.RequestURI, "/1.1/call/") {
			if functions[uri[3]] != nil {

			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func metadataHandler(w http.ResponseWriter, r *http.Request) {
	meta, err := generateMetadata()
	if err != nil {
		errorResponse(w, err)
	}

	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintln(w, meta)
	w.WriteHeader(http.StatusOK)
}

func functionHandler(w http.ResponseWriter, r *http.Request, name string) {
	request, err := constructRequest(r, name)
	if err != nil {
		errorResponse(w, err)
	}

	resp, err := functions[name].call(request)
	if err != nil {
		errorResponse(w, err)
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		errorResponse(w, err)
	}

	w.Header().Add("Contetn-Type", "application/json; charset=UTF-8")
	fmt.Fprintln(w, respJSON)
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {

}

func unmarshalBody(r *http.Request) (interface{}, error) {
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	body := new(map[string]interface{})
	if err := json.Unmarshal(bodyJSON, body); err != nil {
		return nil, err
	}

	return bodyJSON, nil
}

func constructRequest(r *http.Request, name string) (*Request, error) {
	request := new(Request)
	request.Meta = map[string]string{
		"remoteAddr": r.RemoteAddr,
	}
	sessionToken := r.Header.Get("X-LC-Session")
	if !functions[name].NotFetchUser && sessionToken != "" {
		user, err := client.Users.Become(sessionToken)
		if err != nil {
			return nil, err
		}
		request.CurrentUser = user
		request.SessionToken = sessionToken
	}

	params, err := unmarshalBody(r)
	if err != nil {
		return nil, err
	}
	request.Params = params

	return request, nil
}

func errorResponse(w http.ResponseWriter, err error) {
	w.Header().Add("Contetn-Type", "application/json; charset=UTF-8")
	switch err.(type) {
	case *functionError:
		fmt.Fprintln(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
	default:
		fmt.Fprintln(w, Error(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func generateMetadata() ([]byte, error) {
	meta := new(metadataResponse)
	for k := range functions {
		meta.Result = append(meta.Result, k)
	}
	return json.Marshal(meta)
}
