package familiar

import (
	"image"
)

// Auth: var_rain.
// Date: 2024/3/18.
// Desc: 基本实现相关

// -------------------------------------------------------
// 大致逻辑如下:
// 1. 将图像转为灰度
// 2. 计算灰度均值
// 3. 将灰度大等于均值的像素点标记为1, 小于均值的标记为0
// 4. 最终对比相同位数
// 5. 相似度 = 相同位数 / 结果总数 * 100.0
// -------------------------------------------------------

// processToDim1x 处理图像并将结果写入到一维数组,
// 数组元素长度为`img`的`width`x`height`, 其中单个元素
// 表示当前像素的灰度值是否大等于均值. 此操作类似计算图像的指纹信息.
func processToDim1x(img image.Image) []uint8 {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	count := width * height
	pixels := make([]uint8, count)

	var total uint32
	var x, y int
	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// 计算灰度
			r = uint32(float32(r) * 0.299)
			g = uint32(float32(g) * 0.587)
			b = uint32(float32(b) * 0.114)
			u := uint8((r + g + b) / 256)

			// 记录数据
			total += uint32(u)
			pixels[y*height+x] = u
		}
	}

	// 计算灰度均值
	avg := uint8(total / uint32(len(pixels)))

	// 逐像素对比灰度均值大小
	ret := make([]uint8, len(pixels))
	for n := 0; n < count; n++ {
		if pixels[n] >= avg {
			ret[n] = 1
		} else {
			ret[n] = 0
		}
	}

	return ret
}

// computeSimilar 计算两组图像信息的相似度, 两组数组元素长度必须相等,
// 如果两组数据长度不相等, 则该方法不会工作, 将直接以'0'为结果返回.
// 正常的返回结果范围在 0.0-1.0 之间.
func computeSimilar(src []uint8, dst []uint8) float32 {
	if len(src) != len(dst) {
		return 0
	}
	// 对比得出相同的位数
	var size = len(src)
	var counter = 0
	var n int
	for n = 0; n < size; n++ {
		if src[n] == dst[n] {
			counter += 1
		}
	}
	// 得到相似度(0.0-1.0)
	return float32(counter) / float32(size)
}

// u8ListToU64 将一个长度为64且元素只包含`0`和`1`的`uint8`数组转为一个`uint64`值. 数组长度
// 必须为64, 否则结果将直接返回`0`.
func u8ListToU64(src []uint8) uint64 {
	if len(src) != 64 {
		return 0
	}

	var ret uint64 = 0
	for n, u := range src {
		if u == 1 {
			ret |= 1 << (63 - n)
		}
	}
	return ret
}

// u8ListToString 将一个元素只包含`0`和`1`的`uint8`数组转为一个`string`值.
func u8ListToString(src []uint8) string {
	var ret = ""
	for _, u := range src {
		if u == 1 {
			ret += "1"
		} else {
			ret += "0"
		}
	}
	return ret
}
