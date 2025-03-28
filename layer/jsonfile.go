// Package layer 提供了Docker镜像和镜像层管理的功能
package layer

import (
	"encoding/json"
	"fmt"
	"os"
)

// JsonFile 提供JSON文件的读写操作功能
type JsonFile struct {
	// filePath 是JSON文件的路径
	filePath string
}

// NewJsonFile 创建一个新的JsonFile实例
// fp: JSON文件的路径
func NewJsonFile(fp string) (*JsonFile, error) {
	return &JsonFile{
		filePath: fp,
	}, nil
}

// Load 从JSON文件中加载数据到指定的结构体
// src: 目标结构体指针
// 返回值：加载后的结构体，错误信息
func (j *JsonFile) Load(src any) (any, error) {
	if !Exists(j.filePath) {
		return nil, fmt.Errorf("文件%s不存在", j.filePath)
	}
	data, err := os.ReadFile(j.filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, src)
	if err != nil {
		return nil, fmt.Errorf("json 解析失败")
	}
	return src, nil
}

// Dump 将数据以JSON格式写入文件
// src: 要写入的数据
// 返回值：错误信息
func (j *JsonFile) Dump(src any) error {
	if !Exists(j.filePath) {
		_, err := os.Create(j.filePath)
		if err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(src, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(j.filePath, data, 0755)
}

// Exists 检查指定路径的文件是否存在
// path: 文件路径
// 返回值：文件是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
