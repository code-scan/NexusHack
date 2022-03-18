package main

import (
	"flag"
	"nexus-hacker/pkg/docker"
)

func main() {
	var host string
	var registry string
	var thread int
	flag.StringVar(&host, "host", "", "1.1.1.1:8080")
	flag.StringVar(&registry, "registry", "", "xx-registry")
	flag.IntVar(&thread, "thread", 20, "20")

	flag.Parse()
	if host != "" && registry != "" {
		d := docker.NewDocker(host, thread, registry)
		d.GetImages()
		//d.GetBlobs(registry)
		d.ExtractFsLayer(registry)
	}
}
