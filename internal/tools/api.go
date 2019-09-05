package tools

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/energieip/common-components-go/pkg/dswitch"
	"github.com/romana/rlog"
)

var (
	ServerURL  = "127.0.0.1"
	ServerPORT = "8888"
)

func GetSwitchConsumption() (dswitch.SwitchConsumptions, error) {
	var m dswitch.SwitchConsumptions
	url := "https://" + ServerURL + ":" + ServerPORT + "/v1.0/status/consumptions"

	req, _ := http.NewRequest("GET", url, nil)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	req.Close = true
	client := &http.Client{Transport: transCfg}
	resp, err := client.Do(req)

	if err != nil {
		rlog.Error(err.Error())
		return m, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	rlog.Info(string(body))

	err = json.Unmarshal(body, &m)
	return m, err
}
