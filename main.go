package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/gift"
	"github.com/urfave/cli"
)

// Declarations
const imgIcon = "icon.png"
const imgSplash = "splash.png"

func main() {
	//Defining Version
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "prints utility version",
	}

	// Magic happens here
	app := cli.NewApp()
	app.Name = "generate, g"
	app.Version = "1.0.0"
	app.Usage = "generate ionic icons and splash screens!"
	app.Action = func(c *cli.Context) error {
		for _, platform := range Dirs(".") {
			GenerateResources(GetImages(platform))
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Dirs(path string) []string {
	var platforms []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if IsDir(file.Name()) {
			platforms = append(platforms, file.Name())
		}
	}

	return platforms
}

func IsDir(path string) bool {
	fileInfo, _ := os.Stat(path)
	return fileInfo.IsDir()
}

func GenerateResources(pngs []string) {
	for _, p := range pngs {
		width, height := GetImageDimension(p)
		src := LoadImage(imgIcon)
		if strings.Contains(p, "\\splash\\") {
			src = LoadImage(imgSplash)
		}
		filters := map[string]gift.Filter{
			"crop_to_size": gift.ResizeToFill(width, height, gift.LanczosResampling, gift.CenterAnchor),
		}
		for _, filter := range filters {
			g := gift.New(filter)
			dst := image.NewNRGBA(g.Bounds(src.Bounds()))
			g.Draw(dst, src)
			SaveImage(p, dst)
		}
		fmt.Printf("✓ Updated — %s (%dx%d)\n", p, width, height)
	}
}

func GetImages(root string) []string {
	var pngs []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".png" {
			pngs = append(pngs, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return pngs
}

func GetImageDimension(path string) (int, int) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
	}
	return image.Width, image.Height
}

func LoadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("os.Open failed: %v", err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("image.Decode failed: %v", err)
	}
	return img
}

func SaveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("os.Create failed: %v", err)
	}
	err = png.Encode(f, img)
	if err != nil {
		log.Fatalf("png.Encode failed: %v", err)
	}
}
