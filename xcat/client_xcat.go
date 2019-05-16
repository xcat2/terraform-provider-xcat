package xcat

import (
	"net/http"
	"time"
	"fmt"
	"strconv"
)

func FormatResponse(resp interface{}, err error) (interface{}, int, string) {
	if err != nil {
		errorcode_str := fmt.Sprintf("%s", err)
		errorcode, _ := strconv.Atoi(errorcode_str)
		errormsg := resp.(string)
		return "", errorcode, errormsg
	} else {
		return resp, 0, ""
	}
}

func Login(baseUrl string, username string, password string) (string, int, string) {
	url := baseUrl + "/auth/login"
	httpClient := http.Client{Timeout: time.Second * 30}
        client := &HttpClient{Client: &httpClient, Headers: http.Header{}}
	data := make(map[string]interface{})
        data["username"] = username
	data["password"] = password
	ret, errcode, errmsg := FormatResponse(client.Post(url, nil, nil, data, false))
	if errcode != 0 {
		return "", errcode, errmsg
	}	
	token, isok := ret.(map[string]interface {})["token"].(map[string]interface {})["id"].(string)
	if !isok {
		return "", 1, "Failed to get token from response"
	}
	return token, errcode, errmsg
}

func ListNodeAttr(node string, baseUrl string, token interface{}, attr string) (string, int, string) {
	var endpoint string
	var timeout time.Duration
	switch attr {
	case "status":
		timeout = 5
		endpoint = "/_status"
	case "inventory":
		timeout = 15
                endpoint = "/_inventory"
	case "detail":
		timeout = 15
                endpoint = "/_detail"
	default:
		
		return "", 1, "Invalid node attr: " + attr
	}
	url :=  baseUrl + "/system/nodes/" + node + endpoint
	httpClient := http.Client{Timeout: time.Second * timeout}
        client := &HttpClient{Client: &httpClient, Headers: http.Header{}}
        ret, errcode, errmsg := FormatResponse(client.Get(url, nil, token, nil, true))
        return ret.(string), errcode, errmsg
}

func ApplyNodes(baseUrl string, token interface{}, nodeattrs interface{}) (string, int, string) {
	url := baseUrl + "/manager/resmgr"
	httpClient := http.Client{Timeout: time.Second * 15}
        client := &HttpClient{Client: &httpClient, Headers: http.Header{}}
	tmp_data := nodeattrs
	data := make(map[string]interface{})
	data["criteria_spec"] = tmp_data
	data["capacity"] = 1
	ret, errcode, errmsg := FormatResponse(client.Post(url, nil, token, data, false))
	if errcode == 0 {
		for _, value := range ret.(map[string]interface {}) {
        	        return value.(string), errcode, errmsg
	        }
	}
	return "", errcode, errmsg
}

func ListFreeNodes(baseUrl string, token interface{}) (string, int, string) {
	url := baseUrl + "/manager/resmgr"
	httpClient := http.Client{Timeout: time.Second * 15}
        client := &HttpClient{Client: &httpClient, Headers: http.Header{}}
	ret, errcode, errmsg := FormatResponse(client.Get(url, nil, token, nil, true))
        return ret.(string), errcode, errmsg
}

func ReleaseNode(node string, baseUrl string, token interface{}) (string, int, string) {
	url := baseUrl + "/manager/resmgr"
	httpClient := http.Client{Timeout: time.Second * 10}
	client := &HttpClient{Client: &httpClient, Headers: http.Header{}}
	data := make(map[string]string)
	data["name"] = node
	ret, errcode, errmsg := FormatResponse(client.Delete(url, nil, token, data, true))
	return ret.(string), errcode, errmsg
}

func ProvisionNode(node string, baseUrl string, token interface{}, osimage string) (string, int, string) {
	url := baseUrl + "/system/nodes/" + node + "/_operation"
	httpClient := http.Client{Timeout: time.Second * 30}
        client := &HttpClient{Client: &httpClient, Headers: http.Header{}}
	tmp_data := make(map[string]string)
	tmp_data["osimage"] = osimage
	data := make(map[string]interface{})
	data["action_spec"] = tmp_data
	data["action"] = "rinstall"
	ret, errcode, errmsg := FormatResponse(client.Post(url, nil, token, data, true))
	return ret.(string), errcode, errmsg
}

