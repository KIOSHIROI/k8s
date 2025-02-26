package main

import (
	"fmt"
	"layer-scheduler/layer"
)

func main() {
	rr, e := layer.NewRegistry(
		"http://127.0.0.1:5000",
		"admin",
		"admin",
	)

	r, e := rr.GetRemoteImageLayers()
	fmt.Println(e)
	a, b := r.Search(layer.DockerImageName("nginx:la2test"))
	fmt.Println(b)
	fmt.Println(a)
	// fmt.Println(r)
	// // // fmt.Println(ret)
	// // i, e := layer.NewDockerImageLocal()
	// // // fmt.Println(e)
	// // r := i.ListAllLocalImagesInRepo("127.0.0.1")
	// // // fmt.Println(r)
	// // ret, e := rr.GetLocakImageLayers(r)
	// // fmt.Println(ret.GetAllKnownLayers())
	// fmt.Println(e)
	// err := r.Dump("/mnt/z/Code/layer-scheduler/t.json")
	// fmt.Println(err)
	// a, b := layer.NewImageMetadataListFromCache("/tmp/t.json")
	// fmt.Println(b)

	// fmt.Println(a)
}
