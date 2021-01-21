package leancloud

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/levigross/grequests"
)

type requestMethod string

var defaultServerURL = map[string]string{
	"gzGzoHsz": ".lc-cn-n1-shared.com",
	"9Nh9j0Va": ".lc-cn-e1-shared.com",
	"MdYXbMMI": ".api.lncldglobal.com",
}
var requestCount int

const (
	methodGet    requestMethod = "GET"
	methodPost   requestMethod = "POST"
	methodPut    requestMethod = "PUT"
	methodPatch  requestMethod = "PATCH"
	methodDelete requestMethod = "DELETE"
)

type objectResponse map[string]interface{}

type objectsResponse struct {
	Results []objectResponse `json:"results"`
}

type createObjectResponse struct {
	ObjectID  string    `json:"objectId"`
	CreatedAt time.Time `json:"createdAt"`
}

type cqlResponse struct {
	objectsResponse

	ClassName string `json:"className"`
}

type ParseResponseError struct {
	ParseError     error
	ResponseHeader http.Header
	ResponseText   string
	StatusCode     int
	URL            string
}

type ServerResponseError struct {
	Code       int    `json:"code"`
	Err        string `json:"error"`
	StatusCode int
	URL        string
}

func (err *ParseResponseError) Error() string {
	return fmt.Sprintf("parse response failed(%d): %s [%s (%d)]", err.StatusCode, err.ResponseText, err.URL, err.StatusCode)
}

func (err *ServerResponseError) Error() string {
	return fmt.Sprintf("%d %s [%s (%d)]", err.Code, err.Err, err.URL, err.StatusCode)
}

func (client *Client) getServerURL() string {
	if client.serverURL != "" {
		return client.serverURL
	}

	urlSlice := strings.Split(client.appID, "-")
	return fmt.Sprint("https://", (urlSlice[0])[:7], defaultServerURL[urlSlice[1]])
}

func (client *Client) getRequestOptions() *grequests.RequestOptions {
	return &grequests.RequestOptions{
		UserAgent: getUserAgent(),
		Headers: map[string]string{
			"X-LC-Id":  client.appID,
			"X-LC-Key": client.appKey,
		},
	}
}

func (client *Client) request(method requestMethod, path string, options *grequests.RequestOptions, authOptions ...AuthOption) (*grequests.Response, error) {
	if options == nil {
		options = client.getRequestOptions()
	}

	for _, authOption := range authOptions {
		authOption.apply(client, options)
	}

	URL := fmt.Sprint(client.getServerURL(), path)

	requestID := requestCount
	requestCount++

	if client.requestLogger != nil {
		client.requestLogger.Printf("[REQUEST] request(%d) %s %s %#v\n", requestID, method, URL, options)
	}

	resp, err := getRequestAgentByMethod(method)(URL, options)

	if err != nil {
		return resp, err
	}

	if !resp.Ok {
		errResp := &ServerResponseError{}
		err = resp.JSON(errResp)

		if err != nil {
			return resp, &ParseResponseError{
				ParseError:     err,
				ResponseHeader: resp.Header,
				ResponseText:   string(resp.Bytes()),
				StatusCode:     resp.StatusCode,
				URL:            URL,
			}
		}

		errResp.StatusCode = resp.StatusCode
		errResp.URL = URL

		return resp, errResp
	}

	if client.requestLogger != nil {
		client.requestLogger.Printf("[REQUEST] response(%d) %d %s\n", requestID, resp.StatusCode, string(resp.Bytes()))
	}

	return resp, err
}

func getRequestAgentByMethod(method requestMethod) func(string, *grequests.RequestOptions) (*grequests.Response, error) {
	switch method {
	case methodGet:
		return grequests.Get
	case methodPost:
		return grequests.Post
	case methodPut:
		return grequests.Put
	case methodPatch:
		return grequests.Patch
	case methodDelete:
		return grequests.Delete
	default:
		panic(fmt.Sprint("invalid method: ", method))
	}
}

func getUserAgent() string {
	return fmt.Sprint("LeanCloud-Golang-SDK/", Version, " ", runtime.GOOS, "/"+runtime.Version())
}
