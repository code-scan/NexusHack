package tests

import (
	"nexus-hacker/pkg/docker"
	"testing"
)

func TestDocker(t *testing.T) {
	n := docker.NewDocker("http://78.27.198.60:8081/", "", 20, false, "rz-registry")
	n.GetImages()
}
func TestManifests(t *testing.T) {
	n := docker.NewDocker("http://78.27.198.60:8081/", "", 20, false, "rz-registry")
	n.GetManifests("rz-registry", "v2/adminer", "latest")
}

func TestBlob(t *testing.T) {
	n := docker.NewDocker("http://78.27.198.60:8081/", "", 20, false, "rz-registry")
	n.GetAllBlobs("rz-registry")
}
