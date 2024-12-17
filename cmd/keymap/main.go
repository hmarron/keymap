package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

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

	// Draw background

	// Image
	if viper.GetString("background.image") != "" { // Pattern
		im, err := gg.LoadImage(viper.GetString("background.image"))
		if err != nil {
			panic(err)
		}
		dc.DrawImage(im, 0, 0)
	}

	// Solid color
	if colorStr := viper.GetString("background.solid"); colorStr != "" { // Solid
		bgColor := colorFromString(colorStr)
		dc.SetFillStyle(gg.NewSolidPattern(bgColor))
		dc.DrawRectangle(0, 0, float64(W), float64(H))
		dc.Fill()
	}

	// Gradient
	if viper.GetString("background.gradient.x0") != "" { // Gradient
		grad := gg.NewLinearGradient(
			viper.GetFloat64("background.gradient.x0"),
			viper.GetFloat64("background.gradient.y0"),
			viper.GetFloat64("background.gradient.x1"),
			viper.GetFloat64("background.gradient.y1"),
		)
		// Load start color.
		startColor := colorFromString(viper.GetString("background.gradient.start_color"))
		grad.AddColorStop(0, startColor)
		// Load start color.
		endColor := colorFromString(viper.GetString("background.gradient.end_color"))
		grad.AddColorStop(1, endColor)

		dc.SetFillStyle(grad)
		dc.DrawRectangle(0, 0, float64(W), float64(H))
		dc.Fill()
	}
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
			panic("Failed to pack " + title)
		}
		if viper.GetBool("debug") {
			frame.DebugDraw(dc, float64(rect.X), float64(rect.Y))
		} else {
			frame.Draw(dc, float64(rect.X), float64(rect.Y))
		}
	}

	// Save the file
	filename := fmt.Sprintf("%s.png", time.Now().Format("2006-01-02_15-04-05"))
	// filename := fmt.Sprintf("output.png")
	dc.SavePNG(filename)
	fmt.Println(filename)
}

func colorFromString(s string) color.Color {
	colorParts := strings.Split(s, ",")
	if len(colorParts) != 4 {
		panic("Invalid color format. Must be RGBA csv with valse 0-255")
	}
	o := make([]uint8, 4)
	for i, part := range colorParts {
		v, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			panic(fmt.Errorf("Invalid color format. Must be RGBA csv with valse 0-255: %w", err))
		}
		o[i] = uint8(v)
	}

	return color.RGBA{
		o[0],
		o[1],
		o[2],
		o[3],
	}
}
