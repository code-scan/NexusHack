package docker

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
)

func ExtractTarGz(gzipStream io.Reader, dir string) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}
	tarReader := tar.NewReader(uncompressedStream)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
			break
		}
		//log.Println(header.Typeflag)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(dir+"/"+header.Name, 0755); err != nil {
				//log.Printf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(dir + "/" + header.Name)
			if err != nil {
				//log.Printf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				//log.Printf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()
		default:
			//log.Printf(
			//	"ExtractTarGz: uknown type: %s in %s",
			//	header.Typeflag,
			//	header.Name)
		}

	}
}
