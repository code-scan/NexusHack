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
	var keyword string
	flag.StringVar(&host, "host", "", "1.1.1.1:8080")
	flag.StringVar(&registry, "registry", "", "xx-registry")
	flag.IntVar(&thread, "thread", 20, "20")
	flag.BoolVar(&latest, "latest", false, "only download latest")
	flag.StringVar(&keyword, "keyword", "xx/xx", "only download match keyword ")
	flag.Parse()
	if host != "" && registry != "" {
		d := docker.NewDocker(host, keyword, thread, latest, registry)
		d.GetImages()
		d.GetBlobs(registry)
		d.ExtractFsLayer(registry)
		return
	}
	flag.Usage()
}
