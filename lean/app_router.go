package lean

import (
	"fmt"
	"os"
	"strings"

	"github.com/levigross/grequests"
)

type Region string

const (
	RegionCN    Region = "cn"
	RegionCN_N1 Region = "cn-n1"
	RegionCN_E1 Region = "cn-e1"
	RegionUS    Region = "us-w1"
)

type ServiceModule string

const (
	ServiceAPI    ServiceModule = "api_server"
	ServiceEngine ServiceModule = "engine_server"
)

type RouterResponse struct {
	TTL          int    `json:"ttl"`
	APIServer    string `json:"api_server"`
	EngineServer string `json:"engine_server"`
}

var defaultURL = map[Region]string{
	RegionCN:    "https://api.leancloud.cn",
	RegionCN_N1: "https://api.leancloud.cn",
	RegionCN_E1: "https://tab.leancloud.cn",
	RegionUS:    "https://us-api.leancloud.cn",
}

func NewRegionFromString(str string) Region {
	switch strings.ToLower(str) {
	case "", "cn":
		return RegionCN
	case "cn-n1":
		return RegionCN_N1
	case "cn-e1":
		return RegionCN_E1
	case "us", "us-w1":
		return RegionUS
	default:
		panic(fmt.Sprint("invalid region: ", str))
	}
}

func GetServiceURL(region Region, appID string, service ServiceModule) string {
	if region != RegionUS {
		routerInfo, err := queryAppRouter(appID)

		if err != nil {
			fmt.Fprintln(os.Stderr, err) // Ignore app router error
		} else {
			switch service {
			case ServiceAPI:
				return "https://" + routerInfo.APIServer
			case ServiceEngine:
				return "https://" + routerInfo.EngineServer
			}
		}
	}

	return defaultURL[region]
}

// Not applicable for RegionUS
func queryAppRouter(appID string) (result RouterResponse, err error) {
	URL := fmt.Sprint("https://app-router.leancloud.cn/2/route?appId=", appID)

	resp, err := grequests.Get(URL, &grequests.RequestOptions{
		UserAgent: getUserAgent(),
	})

	if err != nil {
		return result, err
	}

	if !resp.Ok {
		return result, fmt.Errorf("query app router failed: %s [%s %d]", string(resp.Bytes()), URL, resp.StatusCode)
	}

	if err = resp.JSON(&result); err != nil {
		return result, err
	}

	return result, nil
}
