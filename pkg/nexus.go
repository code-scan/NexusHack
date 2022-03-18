package pkg

import (
	"github.com/code-scan/Goal/Ghttp"
	"nexus-hacker/pkg/model"
	"strings"
)

type Nexus struct {
	Host     string
	Registry []string
}

func (n *Nexus) SendRequest(action string, data map[string]interface{}) ([]byte, error) {
	body := make(map[string]interface{})
	body["action"] = action
	body["method"] = "read"
	body["tid"] = 1
	body["type"] = "rcp"
	body["data"] = []interface{}{data}
	var ghttp Ghttp.Http
	uri := strings.TrimSuffix(strings.TrimSpace(n.Host), "/") + "/service/extdirect"
	ghttp.New("POST", uri)
	ghttp.SetPostJson(body)
	ghttp.SetContentType("application/json")
	ghttp.Execute()
	return ghttp.Byte()
}
func (n *Nexus) Get(uri string) ([]byte, error) {
	var ghttp Ghttp.Http
	uri = strings.TrimSuffix(strings.TrimSpace(n.Host), "/") + uri
	ghttp.New("GET", uri)
	ghttp.Execute()
	return ghttp.Byte()
}
func (n *Nexus) CoreUiBrowse(repositoryName string, node string) (model.CoreUiBrowser, error) {
	data := make(map[string]interface{})
	data["repositoryName"] = repositoryName
	data["node"] = node
	if resp, err := n.SendRequest("coreui_Browse", data); err == nil {
		return model.UnmarshalCoreUiBrowser(resp)
	} else {
		return model.CoreUiBrowser{}, err
	}
}

func (n *Nexus) Download(uri, file string) {
	var ghttp Ghttp.Http
	uri = strings.TrimSuffix(strings.TrimSpace(n.Host), "/") + uri
	ghttp.New("GET", uri)
	ghttp.Execute()
	ghttp.SaveToFile(file)
}
