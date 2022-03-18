package model

import "encoding/json"

type DockerManifests struct {
	SchemaVersion int    `json:"schemaVersion"`
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	Architecture  string `json:"architecture"`
	FsLayers      []struct {
		BlobSum string `json:"blobSum"`
	} `json:"fsLayers"`
	History []struct {
		V1Compatibility string `json:"v1Compatibility"`
	} `json:"history"`
}

func UnmarshalDockerManifests(data []byte) (DockerManifests, error) {
	var ret DockerManifests
	err := json.Unmarshal(data, &ret)
	return ret, err
}
