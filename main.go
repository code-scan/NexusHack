package main

import (
	"flag"
	"nexus-hacker/pkg/docker"
)

func main() {
	var host string
	var registry string
	var thread int
	var latest bool
	flag.StringVar(&host, "host", "", "1.1.1.1:8080")
	flag.StringVar(&registry, "registry", "", "xx-registry")
	flag.IntVar(&thread, "thread", 20, "20")
	flag.BoolVar(&latest, "latest", false, "only download latest")
	flag.Parse()
	if host != "" && registry != "" {
		d := docker.NewDocker(host, thread, latest, registry)
		d.GetImages()
		d.GetBlobs(registry)
		d.ExtractFsLayer(registry)
	}
}
