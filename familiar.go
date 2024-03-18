package familiar

import (
	"path/filepath"
)

// Auth: var_rain.
// Date: 2024/3/18.
// Desc: 相似度相关

// ------------------------------------------------------------
// ____________________// 大致流程与方式 //______________________
// 1. 将图像加载到内存.
// 2. 缩放图像到一个指定边长的矩形内.
// 	  - 目的是将图像上无关紧要的信息去除掉, 只保留关键信息, 缩放的大小
//	    决定最终得到的结果的精度, 也会影响计算的速度, 缩放越小, 得到
//	    的精度越低, 速度越快, 缩放越大, 精度越高, 速度也就越慢.
// 3. 对图像进行灰度处理.
//    - 处理后得到更精准的图像信息, 既可以区别像素信息又可以去掉其他干
//      扰, 减少计算.
// 4. 计算整幅图像的灰度均值.
// 5. 将所有像素的灰度与均值比较, 结果大等于均值的结果记为1, 小于记为0.
// 	  - 计算后得到图像的指纹信息.
// 6. 将两幅图像的指纹信息的所有位逐一对比, 记录结果相同的位数.
// 7. 以相同的位数除去指纹信息的总长度得到最终的相似度.
//    - 相似度 = 相同位数 / 结果长度 * 100.0
// ------------------------------------------------------------

// GetImageHash 获取图像文件的Hash值, 通过算法处理并计算图像内容, 最终
// 得到一组由'0'和'1'组成的长度为64bit的值. 当图像文件不存在/无法读取/格
// 式不正确时可能会导致错误
func GetImageHash(file string) (*ImageHashCode, error) {
	// 加载图像
	image, err := loadFileToImage(file)
	if err != nil {
		return nil, err
	}

	// 缩放
	image = scaleImageToRect(image, 8)

	// 处理结果
	ret := processToDim1x(image)

	// 抓换并返回结果
	return &ImageHashCode{
		Code:   u8ListToU64(ret),
		Binary: u8ListToString(ret),
	}, nil
}

// GetImageSimilar 对比两张图像的相似度, 结果范围 0.0-1.0 之间.
func GetImageSimilar(src string, dst string) (float32, error) {
	// 加载图像
	srcImage, err := loadFileToImage(src)
	if err != nil {
		return 0, err
	}
	dstImage, err := loadFileToImage(dst)
	if err != nil {
		return 0, err
	}

	// 缩放
	srcImage = scaleImageToRect(srcImage, 8)
	dstImage = scaleImageToRect(dstImage, 8)

	// 计算结果
	srcRet := processToDim1x(srcImage)
	dstRet := processToDim1x(dstImage)

	// 计算相似度并返回
	return computeSimilar(srcRet, dstRet), nil
}

// ScaleImageAndSave 缩放图像之后保存. `in`为输入的文件路径, `out`为
// 保存的文件, `side`为缩放后的图像边长, 文件的保存格式根据`out`的后缀
// 判断, 当前仅支持以`.jpg`, `.jpeg`, `.png` 为后缀的格式, 如果使用此
// 外的格式, 将自动更改为`.jpg`.
func ScaleImageAndSave(in string, out string, side int) error {
	image, err := loadFileToImage(in)
	if err != nil {
		return err
	}

	// 缩放图像
	image = scaleImageToRect(image, side)

	// 检查输出文件后缀, 如果后缀格式不认识直接默认使用`.jpg`后缀
	suffix := filepath.Ext(out)
	switch suffix {
	case ".jpg", ".JPG", ".jpeg", ".JPEG":
		return saveJpegImage(image, out)
	case ".png", ".PNG":
		return savePngImage(image, out)
	default:
		return saveJpegImage(image, out+".jpg")
	}
}

// GrayImageAndSave 将图像灰度处理之后保存. `in`为输入的文件路径,
// `out`为保存的文件, 文件的保存格式根据`out`的后缀判断, 当前仅支持
// 以`.jpg`, `.jpeg`, `.png` 为后缀的格式, 如果使用此外的格式,
// 将自动更改为`.jpg`.
func GrayImageAndSave(in string, out string) error {
	image, err := loadFileToImage(in)
	if err != nil {
		return err
	}

	// 灰度处理
	image = convertToGrayImage(image)

	// 检查输出文件后缀, 如果后缀格式不认识直接默认使用`.jpg`后缀
	suffix := filepath.Ext(out)
	switch suffix {
	case ".jpg", ".JPG", ".jpeg", ".JPEG":
		return saveJpegImage(image, out)
	case ".png", ".PNG":
		return savePngImage(image, out)
	default:
		return saveJpegImage(image, out+".jpg")
	}
}
