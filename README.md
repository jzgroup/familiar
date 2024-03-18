## familiar

`familiar`是`fast-similar`的缩减拼写词,原单词意为`面熟`的意思,在此项目内用于表示`快速相似度`.
此项目为一种`无关图像大小`且`简单粗暴的图像相似度计算方法`,代码使用`Golang`实现,可简单方便的作为动态库/依赖库应用于其他语言或程序中.

### 原理浅析

示例一:

<img src='images/t_image_1.jpg' width="200px" height="150px" alt='原图'/>
缩放
<img src='doc/scale_side_64.jpg' width="75px" height="75px" alt='64x64'/>
灰度
<img src='doc/gray_side_64.jpg' width="75px" height="75px" alt='64x64'/>

 <br>

示例二:

<img src='images/t_image_1.jpg' width="200px" height="150px" alt='原图'/>
缩放
<img src='doc/scale_side_32.jpg' width="75px" height="75px" alt='32x32'/>
灰度
<img src='doc/gray_side_32.jpg' width="75px" height="75px" alt='32x32'/>

计算灰度平均值并逐像素, 大等于平均值的像素记为`1`, 小于平均值的像素记为`0`, 最终得到结果:

`0000110010000000000110001111110011111100111110001111110011011100`

要获得两张图片相似度, 对比得到的结果相同位数即可.

<img src='images/t_image_1.jpg' width="200px" height="150px" alt='图像一'/>

↓↓↓↓↓↓↓↓↓↓↓ (图像仅做学术研究, 不含其他任何目的)

`0000110010000000000110001111110011111100111110001111110011011100`

<img src='images/t_image_2.jpg' width="200px" height="150px" alt='图像二'/>

↓↓↓↓↓↓↓↓↓↓↓ (图像仅做学术研究, 不含其他任何目的)

`0000110010000000000110001111100011111100001111001111110011111110`

如上, 得到的结果有`58`位是相同的, 总共有`64`位, 那么相似度就可以这么计算: 

`58 / 64 = 0.9062`

转为百分比:

`0.9062 * 100.0 = 90.62%`

即两副图像的相似度为`90.62%`

### 使用方式

在源码目录下执行编译
```shell
go build -ldflags "-s -w" ./cmd/cli/
```

计算图像的Hash:
```shell
cli -mode hash -i images/t_image_1.jpg

# 可输入多张图像
cli -mode hash -i images/t_image_1.jpg -i images/t_image_2.jpg -i images/t_image_3.jpg
```

对比图像的相似度:
```shell
cli -mode diff -i images/t_image_1.jpg -i images/t_image_2.jpg
```

获得帮助:
```shell
cli --help
```
