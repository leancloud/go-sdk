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

type mux map[string]http.HandlerFunc

var router mux

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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	bodyMap := new(map[string]interface{})
	if err := json.Unmarshal(body, bodyMap); err != nil {
		fmt.Fprintln(w, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	req := Request{
		Params: bodyMap,
		Meta: http.Request{
			RemoteAddr: r.RemoteAddr,
		},
	}

	if !functions[uri[3]].NotFetchUser {
		sessionToken := r.Header.Get("X-LC-Session")
		user, err := client.Users.Become(sessionToken)
		if err != nil {
			fmt.Fprintln(w, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		req.CurrentUser = user
	}

	if r.Header.Get("X-LC-Session") != "" {
		req.SessionToken = r.Header.Get("X-LC-Session")
	}

	resp, err := functions[uri[3]].call(&req)
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

func generateMetadata() ([]byte, error) {
	meta := new(metadataResponse)
	for k := range functions {
		meta.Result = append(meta.Result, k)
	}
	return json.Marshal(meta)
}
