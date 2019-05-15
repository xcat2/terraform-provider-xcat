package xcat

import (
	"net/http"
)

func Login(baseUrl string, username string, password, string) (string, error) {
	url := baseUrl
	httpClient := http.Client{Timeout: time.Second * 30}
        client := &HttpClient{Client: &httpClient, Headers: http.Header{}}
	data := make(map[string]interface{})
        data["username"] = username
	data["password"] = password
	ret, err := client.Post(url, nil, data)
	if err != nil {
		return "", err + " " + ret
	}
	return ret
}

func ListNodeAttr(node string, baseUrl string, token string, attr string) (string, error) {
	switch attr {
	case "status":
		timeout := 5
		endpoint := "/_status"
	case "inventory":
		timeout := 15
                endpoint := "/_inventory"
	default:
		return "", errors.New("Invalid node attr: " + attr)
	}
	url :=  baseurl + "/system/nodes/" + node + endpoint
	httpClient := http.Client{Timeout: time.Second * timeout}
        client := &HttpClient{Client: &httpClient, Headers: http.Header{}}
        return client.Get(baseUrl, nil, nil)
}

func ProvisionNode(node string, baseUrl string, token string, osimage string) (string, error) {
	url := baseurl + "/system/nodes/" + node + "/_operation"
	httpClient := http.Client{Timeout: time.Second * 30}
        client := &HttpClient{Client: &httpClient, Headers: http.Header{}}
	tmp_data := make(map[string]string)
	tmp_data["osimage"] = osimage
	data := make(map[string]interface{})
	data["action_spec"] = tmp_data
	data["action"] = "rinstall"
	return client.Post(url, nil, data)
}

