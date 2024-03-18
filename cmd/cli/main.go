package main

import (
	"flag"
	"fmt"
	"github.com/jzgroup/familiar"
	"time"
)

// Auth: var_rain.
// Date: 2024/3/17.
// Desc: 简单的命令行工具

func main() {
	// 工作模式
	usage := `指定工作模式, 分别有 "hash", "diff", "scale", "gray" 四种模式:
	"hash"  (默认)模式会逐一计算输入的图像文件, 并得到一组元素长度为64的字符串与一组由64个 "0" 或 "1" 组成的 int64 类型值.
	"diff"  模式会取输入的两张图像文件并通过算法对比得到两张图的相似度百分比.
	"scale" 模式仅将输入图像缩放到指定大小(默认为8)像素边长的正方形.
	"gray"  模式仅将输入图像转为灰度图.`
	mode := flag.String("mode", "hash", usage)

	// 图像输入
	var images ImageArray
	flag.Var(&images, "i", "指定输入图像文件路径, 支持多个输入. 目前仅支持JPEG/PNG格式的图像.")

	// 缩放边长
	side := flag.Int("side", 8, `缩放时图像时矩形的边长. 此参数仅在 "scale" 模式下生效.`)

	// 输出文件名
	output := flag.String("o", "output.jpg", `输出图像文件名. 此参数仅在 "scale", "gray" 模式下生效.`)

	// 解析参数
	flag.Parse()

	// 工作模式检查
	switch *mode {
	case "hash":
		// 检查图像输入
		if len(images) < 1 {
			fmt.Println("没有输入任何文件. 尝试通过 --help 获得更多帮助.")
			return
		}
		for _, img := range images {
			stamp := time.Now().UnixNano()
			hash, err := familiar.GetImageHash(img)
			if err != nil {
				fmt.Printf("图像 %s 解析失败. err: %v \n", img, err)
			} else {
				use := time.Now().UnixNano() - stamp
				fmt.Printf("--------------------------------------\n")
				fmt.Printf("图像: %s\n", img)
				fmt.Printf("HASH: %d\n", hash.Code)
				fmt.Printf("指纹: %s\n", hash.Binary)
				fmt.Printf("耗时: %dns -> %dμs -> %dms\n", use, use/1e3, use/1e6)
				fmt.Printf("--------------------------------------\n")
			}
		}
		fmt.Printf("HASH计算操作完成.\n")
	case "diff":
		// 相似度对比模式必须要有两张以上的图像输入
		if len(images) < 2 {
			fmt.Println("计算相似度需要两张图像才能工作. 尝试通过 --help 获得更多帮助.")
			return
		}
		stamp := time.Now().UnixNano()
		similar, err := familiar.GetImageSimilar(images[0], images[1])
		if err != nil {
			fmt.Printf("图像解析失败. err: %v \n", err)
		} else {
			use := time.Now().UnixNano() - stamp
			fmt.Printf("--------------------------------------\n")
			fmt.Printf("图像: %s, %s\n", images[0], images[1])
			fmt.Printf("相似度: %.4f\n", similar)
			fmt.Printf("百分比: %.4f%%\n", similar*100.0)
			fmt.Printf("耗时: %dns -> %dμs -> %dms\n", use, use/1e3, use/1e6)
			fmt.Printf("--------------------------------------\n")
			fmt.Printf("相似度计算操作完成.\n")
		}
	case "scale":
		// 检查图像输入
		if len(images) < 1 {
			fmt.Println("没有输入任何文件. 尝试通过 --help 获得更多帮助.")
			return
		}
		if err := familiar.ScaleImageAndSave(images[0], *output, *side); err != nil {
			fmt.Printf("图像缩放处理失败. err: %v \n", err)
		} else {
			fmt.Printf("图像缩放操作完成.\n")
		}
	case "gray":
		// 检查图像输入
		if len(images) < 1 {
			fmt.Println("没有输入任何文件. 尝试通过 --help 获得更多帮助.")
			return
		}
		if err := familiar.GrayImageAndSave(images[0], *output); err != nil {
			fmt.Printf("图像转灰度处理失败. err: %v \n", err)
		} else {
			fmt.Printf("图像转灰度操作完成.\n")
		}
	default:
		fmt.Printf("未知的工作模式类型: %s. 请通过 --help 参数获得更多帮助.", mode)
	}
}
