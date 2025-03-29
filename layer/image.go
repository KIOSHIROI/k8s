// Package layer 提供了Docker镜像和镜像层管理的功能
package layer

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"k8s.io/klog/v2"
)

// DockerImageName 表示Docker镜像的完整名称，包含仓库地址、镜像名和标签
type DockerImageName string

// String 返回镜像的完整名称字符串
func (di DockerImageName) String() string {
	return string(di)
}

// Name 返回不包含标签的镜像名称
func (di DockerImageName) Name() string {
	return strings.Join(strings.Split(di.String(), ":")[:2], ":")
}

// NameWithoutRepoAddr 返回不包含仓库地址的镜像名称
func (di DockerImageName) NameWithoutRepoAddr() string {
	return strings.Join(strings.Split(di.Name(), "/")[1:], "/")
}

// Tag 返回镜像的标签
func (di DockerImageName) Tag() string {
	return strings.Split(di.String(), ":")[2]
}

// DockerImages 提供Docker镜像操作的功能
type DockerImages struct {
	// Cli 是Docker客户端实例
	Cli *client.Client
	// CatchFile 是缓存文件的路径
	CatchFile string
}

// NewDockerImageLocal 创建一个连接到本地Docker守护进程的DockerImages实例
func NewDockerImageLocal() (*DockerImages, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerImages{
		Cli: client,
	}, nil
}

// NewDockerImage 创建一个连接到指定Docker守护进程的DockerImages实例
func NewDockerImage(address string, catchFile string) (*DockerImages, error) {
	client, err := client.NewClientWithOpts(client.WithHost("tcp://"+address+":2375"), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerImages{
		Cli:       client,
		CatchFile: catchFile,
	}, nil
}

// ListAllLocalImagesInRepo 列出本地指定仓库中的所有镜像
func (d *DockerImages) ListAllLocalImagesInRepo(repo string) []DockerImageName {
	klog.Info("LisAllLocalImagesInRepo")
	res := []DockerImageName{}
	r, _ := d.Cli.ImageList(context.TODO(), types.ImageListOptions{})
	for _, v := range r {
		for _, tag := range v.RepoTags {
			if strings.HasPrefix(tag, repo) {
				res = append(res, DockerImageName(tag))
				break
			}
		}
	}
	return res
}

// CheckImageExistOnLocal 检查指定镜像是否存在于本地
func (d *DockerImages) CheckImageExistOnLocal(imageName string) (bool, error) {
	arg := filters.NewArgs(filters.KeyValuePair{
		Key:   "reference",
		Value: imageName,
	})
	images, err := d.Cli.ImageList(context.TODO(), types.ImageListOptions{
		All:     true,
		Filters: arg,
	})
	if err != nil || len(images) == 0 {
		return false, err
	}
	return true, nil
}

// GetImageLayer 获取指定镜像的层信息
func (d *DockerImages) GetImageLayer(imageName string, handler *ImageMetadataLists) (ImageMetadata, error) {
	return handler.Search(DockerImageName(imageName))
}
