package familiar

import (
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log/slog"
	"os"
)

// Auth: var_rain.
// Date: 2024/3/18.
// Desc: 基本的图像处理

// loadFileToImage 从文件加载图像, 仅支持JPEG/PNG格式的图像加载, 如果
// 需要加载其他图片格式的文件请自行处理解码器. 文件不存在/不可读取/格式不正确时会发生错误.
func loadFileToImage(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer closer(f)

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// scaleImageToRect 将图像缩放为指定边长的矩形, 使用双线性算法缩放图像,
// 缩放后得到边长为`side`, 颜色模式为`RGBA`的正方形图像.
func scaleImageToRect(img image.Image, side int) image.Image {
	// 缩放后的图像大小 (x,y,w,h)
	dstRect := image.Rect(0, 0, side, side)

	// 创建缩放后的目标图像
	dst := image.NewRGBA(dstRect)

	// 使用双线性算法缩放图像
	// 使用最邻近算法或近似双线性可能会导致图像部分细节信息丢失
	// 无关紧要可修改为其他算法以加快处理速度
	draw.BiLinear.Scale(dst, dstRect, img, img.Bounds(), draw.Over, nil)

	return dst
}

// convertToGrayImage 将图像转为灰度图.
func convertToGrayImage(img image.Image) image.Image {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	var gray = image.NewGray(image.Rect(0, 0, width, height))
	var x, y int
	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// 计算灰度
			r = uint32(float32(r) * 0.299)
			g = uint32(float32(g) * 0.587)
			b = uint32(float32(b) * 0.114)

			// 设置像素值
			gray.Set(x, y, color.Gray{Y: uint8((r + g + b) / 256)})
		}
	}

	return gray
}

// saveJpegImage 将图像编码到以`.jpg`或`.jpeg`结尾的文件中. 如果图
// 像已存在, 源文件将会被直接覆盖, 图像不存在, 将会创建一个新的文件, 并
// 写入数据.
func saveJpegImage(img image.Image, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer closer(f)

	return jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
}

// savePngImage 将图像编码到以`.png`结尾的文件中. 如果图像已存在, 源
// 文件将会被直接覆盖, 图像不存在, 将会创建一个新的文件, 并写入数据.
func savePngImage(img image.Image, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer closer(f)

	return png.Encode(f, img)
}

// closer 单纯是因为这么长串方法写在`defer`后面不好看所以单独抽出来了.
func closer(f io.Closer) {
	if err := f.Close(); err != nil {
		slog.Warn("关闭文件时出错: {}", slog.AnyValue(err))
	}
}
