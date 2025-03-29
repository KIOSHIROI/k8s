// Package layer 提供了Docker镜像和镜像层管理的功能
package layer

import (
	"bytes"
	"encoding/json"
	"os"

	"k8s.io/klog/v2"
)

// LayerMetadata 表示Docker镜像层的元数据信息
type LayerMetadata struct {
	// Size 是镜像层的大小（字节）
	Size int64 `json:"size"`
	// Layer 是镜像层的标识符
	Layer string `json:"layer"`
}

// ImageMetadata 表示Docker镜像的元数据信息
type ImageMetadata struct {
	// Id 是镜像的唯一标识符
	Id string `json:"id"`
	// Name 是镜像的名称
	Name string `json:"name"`
	// NameWithoutRepo 是不包含仓库地址的镜像名称
	NameWithoutRepo string `json:"name_without_repo"`
	// Tag 是镜像的标签
	Tag string `json:"tag"`
	// TotalSize 是镜像的总大小（字节）
	TotalSize int64 `json:"total_size"`
	// LayerMetadata 是镜像的所有层信息
	LayerMetadata []LayerMetadata `json:"layer_metadata"`
}

// ImageMetadataLists 管理镜像元数据的集合
type ImageMetadataLists struct {
	// CatchFile 是缓存文件的路径
	CatchFile string
	// Lists 存储镜像名称到镜像元数据的映射
	Lists map[string]ImageMetadata
}

// NewImageMetadataListFromCache 从缓存文件中加载镜像元数据列表
func NewImageMetadataListFromCache(filePath string) (*ImageMetadataLists, error) {
	klog.Infoln("Load image metadata from cache file", filePath)
	jf, err := NewJsonFile(filePath)
	if err != nil {
		return &ImageMetadataLists{}, err
	}
	res := &ImageMetadataLists{}
	_, err = jf.Load(res)
	res.CatchFile = filePath
	return res, err
}

// GetAllKnownLayers 获取所有已知的镜像层信息
func (re *ImageMetadataLists) GetAllKnownLayers() []LayerMetadata {
	res := []LayerMetadata{}
	for _, mt := range re.Lists {
		for _, layerStr := range mt.LayerMetadata {
			res = append(res, layerStr)
		}
	}
	return res
}

// Dump 将镜像元数据列表保存到指定文件
func (re *ImageMetadataLists) Dump(filePath string) error {
	jf, err := NewJsonFile(filePath)
	if err != nil {
		return err
	}
	return jf.Dump(&re)
}

// Fromat 将镜像元数据格式化为美观的JSON字符串
func (re *ImageMetadataLists) Fromat() (bytes.Buffer, error) {
	var str bytes.Buffer
	b, err := json.Marshal(re)

	if err != nil {
		return str, err
	}
	_ = json.Indent(&str, b, "", "     ")
	return str, nil
}

// Search 根据镜像名称搜索镜像元数据
func (re *ImageMetadataLists) Search(image DockerImageName) (ImageMetadata, error) {
	res, ok := re.Lists[image.NameWithoutRepoAddr()]
	if ok {
		return res, nil
	}
	return res, os.ErrNotExist
}

// SearchLayer 根据层标识符搜索层大小
func (re *ImageMetadataLists) SearchLayer(layer string) int64 {
	allLayer := re.GetAllKnownLayers()
	for _, l := range allLayer {
		if l.Layer == layer {
			return l.Size
		}
	}
	return 0
}

// ComputeLayerSize 计算多个镜像层的总大小
func ComputeLayerSize(metadata []LayerMetadata) int64 {
	res := int64(0)
	for _, data := range metadata {
		res += data.Size
	}
	return res
}
