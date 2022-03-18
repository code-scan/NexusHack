package main

import (
	"flag"
	"nexus-hacker/pkg/docker"
)

func main() {
	var host string
	var registry string
	flag.StringVar(&host, "host", "", "1.1.1.1:8080")
	flag.StringVar(&registry, "registry", "", "xx-registry")
	flag.Parse()
	if host != "" && registry != "" {
		d := docker.NewDocker(host, registry)
		d.GetImages()
		//d.GetBlobs(registry)
		d.ExtractFsLayer(registry)
	}
}
