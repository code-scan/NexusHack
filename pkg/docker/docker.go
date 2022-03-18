package docker

import (
	"fmt"
	"github.com/code-scan/Goal/Gconvert"
	"log"
	"net/url"
	"nexus-hacker/pkg"
	"nexus-hacker/pkg/model"
	"os"
	"strings"
	"sync"
	"time"
)

type Docker struct {
	pkg.Nexus
	thread  int
	latest  bool
	keyword string
	blobs   map[string]int
	images  map[string][]string
}

func NewDocker(host, keyword string, thread int, latest bool, registry ...string) *Docker {
	var docker Docker
	docker.Host = host
	docker.keyword = keyword
	docker.thread = thread
	docker.latest = latest
	docker.Registry = registry
	docker.blobs = make(map[string]int)
	docker.images = make(map[string][]string)
	return &docker
}

func (d *Docker) GetImages() {
	for _, r := range d.Registry {
		d.getImages(r, "v2")
	}
}
func (d *Docker) getImages(registry, node string) {

	ret, err := d.CoreUiBrowse(registry, node)
	if err != nil {
		log.Println(err)
		return
	}
	for _, image := range ret.Result.Data {
		log.Printf("[*] Folder: %s ", image.Id)
		if image.Text == "blobs" {
			continue
		}
		if image.Text == "tags" || image.Text == "manifests" {
			d.GetTags(registry, node)
			break
		}
		d.getImages(registry, image.Id)
	}
}
func (d *Docker) GetTags(registry string, node string) {
	log.Printf("[*] Get Tags: %s ", node)
	if d.keyword != "" {
		if strings.Contains(node, d.keyword) == false {
			log.Printf("[*] not Match ,Skip")
			return
		}
	}
	ret, err := d.CoreUiBrowse(registry, node+"/tags")
	if err != nil {
		log.Println(err)
		return
	}
	for _, tag := range ret.Result.Data {
		log.Printf("	[#] Tags: %s ", tag.Id)
		d.GetManifests(registry, node, tag.Text)
		if d.latest {
			break
		}
	}
}

func (d *Docker) GetManifests(registry, node, tag string) {
	uri := fmt.Sprintf("/repository/%s/%s/manifests/%s", registry, node, tag)
	ret, err := d.Get(uri)
	if err != nil {
		log.Println(err)
		return
	}
	manifests, err := model.UnmarshalDockerManifests(ret)
	if err != nil {
		log.Println(err)
		return
	}
	name := fmt.Sprintf("%s/%s", node, tag)
	for _, fs := range manifests.FsLayers {
		d.blobs[fs.BlobSum] = 1
		d.images[name] = append(d.images[name], fs.BlobSum)
		log.Println("		[$] FsLayer: ", fs.BlobSum)
	}
}

// GetAllBlobs 下载所有的blobs包括没有被任何镜像使用的
func (d *Docker) GetAllBlobs(registry string) {
	ret, err := d.CoreUiBrowse(registry, "v2/blobs")
	if err != nil {
		log.Println(err)
	}
	wg := sync.WaitGroup{}
	hostname, _ := url.Parse(d.Host)
	os.MkdirAll(fmt.Sprintf("out/%s/%s/blobs/", hostname.Host, registry), 0777)
	tasks := make(chan string, 100)
	for i := 0; i < d.thread; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case t := <-tasks:
					uri := fmt.Sprintf("/repository/%s/v2/-/blobs/%s", registry, t)
					f := strings.Split(Gconvert.UrlDecode(t), ":")
					var filename string
					if len(f) == 2 {
						filename = fmt.Sprintf("out/%s/%s/blobs/%s", hostname.Host, registry, f[1])
						log.Println("[*] Start Download: ", filename)
						if d.CheckFileExist(filename) {
							log.Println("[*] File Exist: ", filename)
							continue
						}
						d.Download(uri, filename)
						log.Println("[*] Over Download: ", filename)
					}
				case <-time.After(10):
					break
				}
			}
			wg.Done()
		}()
	}
	for _, blob := range ret.Result.Data {
		log.Println("[*] Blob: ", blob.Text)
		tasks <- blob.Id
	}
	wg.Wait()
}

// GetBlobs 下载所有镜像的fslayer
func (d *Docker) GetBlobs(registry string) {
	wg := sync.WaitGroup{}
	hostname, _ := url.Parse(d.Host)
	os.MkdirAll(fmt.Sprintf("out/%s/%s/blobs/", hostname.Host, registry), 0777)
	tasks := make(chan string, 100)
	for i := 0; i < d.thread; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			for {
				select {
				case t := <-tasks:
					uri := fmt.Sprintf("/repository/%s/v2/-/blobs/%s", registry, t)
					f := strings.Split(Gconvert.UrlDecode(t), ":")
					var filename string
					if len(f) == 2 {
						filename = fmt.Sprintf("out/%s/%s/blobs/%s", hostname.Host, registry, f[1])
						log.Println("[*] Start Download: ", filename)
						if d.CheckFileExist(filename) {
							log.Println("[*] File Exist: ", filename)
							continue
						}
						d.Download(uri, filename)
						log.Println("[*] Over Download: ", filename)
					}
				case <-time.After(5):
					return
				}
			}
		}()
	}
	log.Println("[*] Total Blob: ", len(d.blobs))
	for blob, _ := range d.blobs {
		log.Println("[*] Blob: ", blob)
		tasks <- blob
	}
	wg.Wait()
}

// ExtractFsLayer 从fslayer中解压并生成文件夹
func (d *Docker) ExtractFsLayer(registry string) {
	hostname, _ := url.Parse(d.Host)
	for image, fs := range d.images {
		dir := fmt.Sprintf("out/%s/%s/%s/", hostname.Host, registry, image)
		os.MkdirAll(dir, 0777)
		for _, f := range fs {
			f = strings.ReplaceAll(f, "sha256:", "")
			blob := fmt.Sprintf("out/%s/%s/blobs/%s", hostname.Host, registry, f)
			r, err := os.Open(blob)
			if err != nil {
				fmt.Println("error")
			}
			log.Println("[*] ExtractTarGz : ", blob)
			ExtractTarGz(r, dir)
		}
	}
}
func (d *Docker) CheckFileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	} else {
		return true
	}
	return true
}
