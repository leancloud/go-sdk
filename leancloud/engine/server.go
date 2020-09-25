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
					functionHandler(w, r, uri)
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
		fmt.Fprintln(w, err)
	}
	fmt.Fprintln(w, meta)
}

func functionHandler(w http.ResponseWriter, r *http.Request, uri []string) {
	request, err := generateRequest(r, uri)
	if err != nil {
		fmt.Fprintln(w, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := functions[uri[3]].call(request)
	if err != nil {
		fmt.Fprintln(w, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintln(w, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintln(w, respJSON)
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {

}

func marshalBody(r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	bodyJSON := new(map[string]interface{})
	if err := json.Unmarshal(body, bodyJSON); err != nil {
		return nil, err
	}

	return bodyJSON, nil
}

func generateRequest(r *http.Request, uri []string) (*Request, error) {
	request := new(Request)
	request.Meta = http.Request{
		RemoteAddr: r.RemoteAddr,
	}
	sessionToken := r.Header.Get("X-LC-Session")
	if !functions[uri[3]].NotFetchUser && sessionToken != "" {
		user, err := client.Users.Become(sessionToken)
		if err != nil {
			return nil, err
		}
		request.CurrentUser = user
		request.SessionToken = sessionToken
	}

	params, err := marshalBody(r)
	if err != nil {
		return nil, err
	}
	request.Params = params

	return request, nil
}

func generateMetadata() ([]byte, error) {
	meta := new(metadataResponse)
	for k := range functions {
		meta.Result = append(meta.Result, k)
	}
	return json.Marshal(meta)
}
