package main

import (
	"fmt"
	"image/color"

	packer "github.com/InfinityTools/go-binpack2d"
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
		FramePaddingX:    viper.GetFloat64("layout.frame_padding_x"),
		FramePaddingY:    viper.GetFloat64("layout.frame_padding_y"),
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
	frameConfigs := viper.GetStringMapStringSlice("frames")
	frames := make([]internal.Frame, 0, len(frameConfigs))
	p := packer.Create(W, H)
	for title, config := range frameConfigs {
		frame := internal.Frame{
			Title:    title,
			Style:    frameStyle,
			Contents: config,
		}
		frames = append(frames, frame)

		w, h := frame.GetCalculatedDimensions(dc)
		rect, ok := p.Insert(int(w), int(h), packer.RULE_BEST_LONG_SIDE_FIT)
		if !ok {
			panic("Failed to pack")
		}
		frame.Draw(dc, float64(rect.X), float64(rect.Y))
	}

	// Save the file
	// filename := fmt.Sprintf("%s.png", time.Now().Format("2006-01-02_15-04-05"))
	filename := fmt.Sprintf("output.png")
	dc.SavePNG(filename)
	fmt.Println(filename)
}
