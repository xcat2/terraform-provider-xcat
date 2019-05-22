package xcat

import (
	"fmt"
	//"io/ioutil"
	//"net/http"
        "os/exec"
	//"github.com/hashicorp/terraform/helper/pathorcontents"
)

type Config struct {
	Url      string
	Username string
	Password string

	Token string
}

func (c *Config) loadAndValidate() error {
	/*
        response, err := http.Get(c.Url + "/login/?username=" + c.Username + "/?password=" + c.Password)
	if err != nil {
		return fmt.Errorf("Error to apply token: %s", err)
	}

	respdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error to parse response: %s", err)
	}
	c.Token = string(respdata)
        */
        cmd := exec.Command("lsxcatd","-v")
        stdout, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Error to apply token: %s", err)
	}

        fmt.Printf(string(stdout))



	return nil
}
