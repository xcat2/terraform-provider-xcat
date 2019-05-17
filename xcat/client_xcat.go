package xcat 

import (
	"github.com/tidwall/gjson"
	"net/http"
	"time"
	"fmt"
	"strconv"
	"strings"
	"crypto/tls"
)

func FormatResponse(resp interface{}, err error) (interface{}, int, string) {
	if err != nil {
		errorcode_str := fmt.Sprintf("%s", err)	
		errorcode, _ := strconv.Atoi(errorcode_str)
		errormsg := ""
 		if strings.Contains(errorcode_str, "timeout") {
			errorcode = 504
 			errormsg = "Timeout, please try again later"
		}
 		if strings.Contains(errorcode_str, "no such host") {
 			errorcode = 1
 			errormsg = "Failed to resolve host, please check"
 		}
 		if strings.Contains(errorcode_str, "Can not read the message form response") {
			errorcode = 1
			errormsg = errorcode_str
 		}
		if strings.Contains(errorcode_str, "connection reset by peer") {
			errorcode = 1
			errormsg = errorcode_str
		}
 		if resp != nil {
 			errormsg = resp.(string)
 		}
		return "", errorcode, errormsg
	} else {
		return resp, 0, ""
	}
}

func GenerateClient(baseUrl string, timeout time.Duration) *HttpClient {
	httpClient := http.Client{Timeout: time.Second * timeout}
	if strings.Contains(baseUrl, "https://") {
		tr := &http.Transport{
        	        TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
        	}
		httpClient = http.Client{Timeout: time.Second * timeout,
                                Transport: tr}
	}
	return &HttpClient{Client: &httpClient, Headers: http.Header{}}
}

func Login(baseUrl string, username string, password string) (string, int, string) {
	url := baseUrl + "/auth/login"
	client := GenerateClient(baseUrl, 20)
	data := make(map[string]interface{})
        data["username"] = username
	data["password"] = password
	ret, errcode, errmsg := FormatResponse(client.Post(url, nil, nil, data, false))
	if errcode != 0 {
		return "", errcode, errmsg
	}	
	token := gjson.Get(ret.(string), "token.id")
	return token.String(), errcode, errmsg
}

func ApplyNodes(baseUrl string, token interface{}, nodeattrs interface{}) (string, int, string) {
	url := baseUrl + "/manager/resmgr"
	client := GenerateClient(baseUrl, 30)
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

func ListNodeStatus(node string, baseUrl string, token interface{}) (string, int, string) {
	url := baseUrl + "/system/nodes/" + node + "/_status"
	client := GenerateClient(baseUrl, 10)
	ret, errcode, errmsg := FormatResponse(client.Get(url, nil, token, nil, false))
	if errcode != 0 {
		return "", errcode, errmsg
	}
	status := gjson.Get(ret.(string), "status.boot.state")
	return status.String(), 0, ""
}

func ListNodeDetail(node string, baseUrl string, token interface{}) (string, int, string) {
        url :=  baseUrl + "/system/nodes/" + node + "/_detail"
	client := GenerateClient(baseUrl, 15)
        ret, errcode, errmsg := FormatResponse(client.Get(url, nil, token, nil, true))
        return ret.(string), errcode, errmsg
}

func ReleaseNode(node string, baseUrl string, token interface{}) (string, int, string) {
	url := baseUrl + "/manager/resmgr" + "?name=" + node
	client := GenerateClient(baseUrl, 10)
	ret, errcode, errmsg := FormatResponse(client.Delete(url, nil, token, nil, true))
	return ret.(string), errcode, errmsg
}

func SetPowerStatus(node string, baseUrl string, token interface{}, status string) (string, int, string) {
	url := baseUrl + "/system/nodes/" + node + "/power" + status
	client := GenerateClient(baseUrl, 10)
	data := make(map[string]interface{})
	data["status"] = status
	ret, errcode, errmsg := FormatResponse(client.Post(url, nil, token, data, true))
	return ret.(string), errcode, errmsg
}

func ProvisionNode(node string, baseUrl string, token interface{}, osimage string) (string, int, string) {
	url := baseUrl + "/system/nodes/" + node + "/_operation"
	client := GenerateClient(baseUrl, 30)
	tmp_data := make(map[string]string)
	tmp_data["osimage"] = osimage
	data := make(map[string]interface{})
	data["action_spec"] = tmp_data
	data["action"] = "rinstall"
	ret, errcode, errmsg := FormatResponse(client.Post(url, nil, token, data, true))
	return ret.(string), errcode, errmsg
}

