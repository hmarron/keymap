package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/fogleman/gg"
	"github.com/hmarron/keymap/internal"
	"github.com/spf13/viper"
)

func main() {
	// Load config
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Load canvas size.
	W := viper.GetInt("width")
	H := viper.GetInt("height")

	// Load frame style.
	frameStyle := internal.FrameStyle{
		FrameStrokeWidth: viper.GetFloat64("frame_style.stroke"),
		TitleFontSize:    viper.GetFloat64("frame_style.title_font-size"),
		ContentFontSize:  viper.GetFloat64("frame_style.content_font-size"),
		TitlePadding:     viper.GetFloat64("frame_style.title_padding"),
		ContentPadding:   viper.GetFloat64("frame_style.content_padding"),
		FontPath:         viper.GetString("font_path"),
	}

	dc := gg.NewContext(W, H)

	// Set up background gradient
	// TODO need to figure out how to move this to config
	//      Maybe allow solid colors and image backgrounds?
	grad := gg.NewLinearGradient(0, 200, 0, float64(H))
	grad.AddColorStop(0, color.RGBA{0, 0, 50, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 0, 255})
	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, float64(W), float64(H))
	dc.Fill()

	// Load frames
	framePaddingX := viper.GetFloat64("layout.frame_padding_x")
	framePaddingY := viper.GetFloat64("layout.frame_padding_x")
	frameCols := viper.GetInt("layout.cols")
	frameX := framePaddingX
	frameY := framePaddingY

	// TODO need to figure out the total size including padding
	// and then find some algo to better arrange these
	// maybe  something like a grid layout
	// that goes a row at a time
	// each cell uses the above cell to figure out it's Y position
	// each cell uses the left cell to figure out it's X position
	// if the new cell is going to go off the page, then start a new row

	frameConfigs := viper.GetStringMapStringSlice("frames")
	count := 1
	for title, config := range frameConfigs {
		frame := internal.Frame{
			PositionX: frameX,
			PositionY: frameY,
			Title:     title,
			Style:     frameStyle,
			Contents:  config,
		}
		fW, fH := frame.GetCalculatedDimensions(dc)

		// TODO maybe have some auto mode where it only moves to next line if
		//  it's going to go off the canvas-framepadding
		//  might need to sort them or have some way of getting the next frame for that...
		if count%frameCols == 0 {
			frameX = framePaddingX
			// TODO this should use the biggest fH of the current row
			frameY += fH + framePaddingY
		} else {
			frameX += fW + framePaddingX
		}
		count++
		frame.Draw(dc)
	}

	// Save the file
	filename := fmt.Sprintf("%s.png", time.Now().Format("2006-01-02_15-04-05"))
	dc.SavePNG(filename)
	fmt.Println(filename)
}
