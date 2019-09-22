package consumer

import (
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/you/hello/pkg/entry"
	"image"
	"math"
	"path"
)

type Consumer struct {
	dirPath string
	widthHeightRatio float64
}

func NewConsumer(dirPath string, widthHeightRatio float64) *Consumer {
	return &Consumer{
		dirPath: dirPath,
		widthHeightRatio: widthHeightRatio,
	}
}

func (c *Consumer) PutEntry(entry *entry.Entry) error {
	bkImg := c.createBackgroundImg(entry.Image)
	resultImg := imaging.PasteCenter(bkImg, entry.Image)
	filePath := path.Join(c.dirPath, entry.Name + ".jpg")
	err := imaging.Save(resultImg, filePath, imaging.JPEGQuality(90))
	if err != nil {
		return errors.Wrapf(err, "failed to save image")
	}
	return nil
}

type backgroundProperties struct {
	croppedWidth int
	croppedHeight int
	resizedWidth int
	resizedHeight int
}

func (c *Consumer) createBackgroundImg(srcImg image.Image) image.Image {
	properties := c.createBackgroundProperties(srcImg)
	croppedImg := imaging.CropCenter(srcImg, properties.croppedWidth, properties.croppedHeight)
	resizedImg := imaging.Resize(croppedImg, properties.resizedWidth, properties.resizedHeight, imaging.CatmullRom)
	return imaging.Blur(resizedImg, 20)
}

func (c *Consumer) createBackgroundProperties(srcImg image.Image) *backgroundProperties {
	bounds := srcImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	widthHeightRatio := float64(width) / float64(height)
	properties := &backgroundProperties{}
	properties.croppedWidth = width
	properties.croppedHeight = height
	properties.resizedWidth = width
	properties.resizedHeight = height
	if widthHeightRatio > c.widthHeightRatio {
		properties.croppedWidth = int(math.Round(float64(height) * c.widthHeightRatio))
		properties.resizedHeight = int(math.Round(float64(width) / c.widthHeightRatio))
	} else {
		properties.croppedHeight = int(math.Round(float64(width) / c.widthHeightRatio))
		properties.resizedWidth = int(math.Round(float64(height) * c.widthHeightRatio))
	}
	return properties
}
