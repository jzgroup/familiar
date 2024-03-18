package main

import "fmt"

// Auth: var_rain.
// Date: 2024/3/18.
// Desc: 命令行参数接收类型定义

// ImageArray 定义用于从命令行接收多个图像参数的数组类型值
type ImageArray []string

func (i *ImageArray) String() string {
	return fmt.Sprint(*i)
}

func (i *ImageArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}
