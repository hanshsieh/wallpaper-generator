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
	width int
	height int
}

func NewConsumer(dirPath string, width int, height int) *Consumer {
	return &Consumer{
		dirPath: dirPath,
		width: width,
		height: height,
	}
}

func (c *Consumer) PutEntry(entry *entry.Entry) error {
	props := c.createProperties(entry.Image)
	bkImg := c.createBackgroundImg(entry.Image, props)
	fgImg := c.createForegroundImg(entry.Image, props)
	resultImg := imaging.PasteCenter(bkImg, fgImg)
	filePath := path.Join(c.dirPath, entry.Name + ".jpg")
	err := imaging.Save(resultImg, filePath, imaging.JPEGQuality(90))
	if err != nil {
		return errors.Wrapf(err, "failed to save image")
	}
	return nil
}

type properties struct {
	bkCroppedWidth  int
	bkCroppedHeight int
	fgWidth int
	fgHeight int
}

func (c *Consumer) createBackgroundImg(srcImg image.Image, props *properties) image.Image {
	croppedImg := imaging.CropCenter(srcImg, props.bkCroppedWidth, props.bkCroppedHeight)
	resizedImg := imaging.Resize(croppedImg, c.width, c.height, imaging.CatmullRom)
	return imaging.Blur(resizedImg, 20)
}

func (c *Consumer) createForegroundImg(srcImg image.Image, props *properties) image.Image {
	return imaging.Resize(srcImg, props.fgWidth, props.fgHeight, imaging.CatmullRom)
}

func (c *Consumer) createProperties(srcImg image.Image) *properties {
	bounds := srcImg.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()
	srcWidthHeightRatio := float64(srcWidth) / float64(srcHeight)
	props := &properties{}
	props.bkCroppedWidth = srcWidth
	props.bkCroppedHeight = srcHeight
	props.fgWidth = c.width
	props.fgHeight = c.height
	targetWidthHeightRatio := float64(c.width) / float64(c.height)
	if srcWidthHeightRatio > targetWidthHeightRatio {
		props.bkCroppedWidth = int(math.Round(float64(srcHeight) * targetWidthHeightRatio))
		props.fgHeight = int(math.Round(float64(props.fgWidth) / srcWidthHeightRatio))
	} else {
		props.bkCroppedHeight = int(math.Round(float64(srcWidth) / targetWidthHeightRatio))
		props.fgWidth = int(math.Round(float64(props.fgHeight) * srcWidthHeightRatio))
	}
	return props
}
