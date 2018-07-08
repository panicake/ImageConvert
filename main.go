package main

import (
	"github.com/spf13/pflag"
	"os"
	"log"
	"github.com/disintegration/imaging"
	"io/ioutil"
	"encoding/json"
	"strings"
	"image"
	"fmt"
)

var convertType string
var path string
var config string
func AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&convertType, "type", "icon", "convert icon, splash")
	fs.StringVar(&path, "path", "icon.png", "image path")
	fs.StringVar(&config, "config", "devices.json", "image path")
	fs.Parse(os.Args[1:])
}

type SizeSpec struct {
	Path string `json:"path"`
	Width int   `json:"width"`
	Height int  `json:"height"`
}
type DeviceConfig struct {
	Icon map[string][]SizeSpec `json:"icon"`
	Splash map[string][]SizeSpec `json:"splash"`
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func generateImage(img image.Image, sizeSpec map[string][]SizeSpec)  {
	for k, v := range sizeSpec {
		log.Printf("generate %s icon", k)
		for _, spec := range v {
			index := strings.LastIndex(spec.Path,"/")
			path := spec.Path[:index]
			ok, _ := pathExists(path)
			if !ok {
				err := os.MkdirAll(path, os.ModePerm)
				if err != nil {
					fmt.Printf("mkdir failed![%v]\n", err)
				} else {
					fmt.Printf("mkdir success!\n")
				}
			}
			dist := imaging.Resize(img, spec.Width, spec.Height, imaging.Lanczos)
			imaging.Save(dist, spec.Path)
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	AddFlags(pflag.CommandLine)
	if convertType == "" {
		convertType = "icon"
	}
	if path == "" {
		log.Fatalln("Image path is empty")
		return
	}
	devices, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	var deviceConfig DeviceConfig
	err = json.Unmarshal(devices, &deviceConfig)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	originImg, err := imaging.Open(path)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	if convertType == "icon" {
		generateImage(originImg, deviceConfig.Icon)
	} else if convertType == "splash" {
		generateImage(originImg, deviceConfig.Splash)
	}

	if err != nil {
		log.Fatalln(err.Error())
	}
}


