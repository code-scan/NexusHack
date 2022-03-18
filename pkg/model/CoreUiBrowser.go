package model

import "encoding/json"

type CoreUiBrowser struct {
	Tid    int    `json:"tid"`
	Action string `json:"action"`
	Method string `json:"method"`
	Result struct {
		Success bool `json:"success"`
		Data    []struct {
			Id          string      `json:"id"`
			Text        string      `json:"text"`
			Type        string      `json:"type"`
			Leaf        bool        `json:"leaf"`
			ComponentId interface{} `json:"componentId"`
			AssetId     interface{} `json:"assetId"`
			PackageUrl  interface{} `json:"packageUrl"`
		} `json:"data"`
	} `json:"result"`
	Type string `json:"type"`
}

func UnmarshalCoreUiBrowser(data []byte) (CoreUiBrowser, error) {
	var core CoreUiBrowser
	err := json.Unmarshal(data, &core)
	return core, err
}
