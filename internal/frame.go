package internal

import (
	"image/color"

	"github.com/fogleman/gg"
)

type FrameStyle struct {
	FrameStrokeWidth float64
	TitleFontSize    float64
	ContentFontSize  float64
	TitlePadding     float64
	ContentPadding   float64
	FontPath         string
}

type Frame struct {
	PositionX float64
	PositionY float64
	Title     string
	Contents  []string
	Style     FrameStyle
}

func (f *Frame) GetCalculatedDimensions(dc *gg.Context) (float64, float64) {
	// Get title height.
	dc.LoadFontFace(f.Style.FontPath, f.Style.TitleFontSize)
	_, titleH := dc.MeasureString(f.Title)

	dc.LoadFontFace(f.Style.FontPath, f.Style.ContentFontSize)
	var W, H float64
	currentDrawY := f.PositionY + f.Style.FrameStrokeWidth + f.Style.TitlePadding + titleH
	for _, content := range f.Contents {
		w, h := dc.MeasureString(content)
		H += h + f.Style.ContentPadding
		if w > W {
			W = w
		}
		currentDrawY += h
		currentDrawY += f.Style.ContentPadding
	}

	return W + (f.Style.ContentPadding * 2), H + f.Style.ContentPadding + f.Style.TitlePadding + titleH
}

func (f *Frame) Draw(dc *gg.Context) {
	// Set default styles
	dc.SetColor(color.White)
	dc.SetLineWidth(f.Style.FrameStrokeWidth)

	// Draw the title
	dc.LoadFontFace(f.Style.FontPath, f.Style.TitleFontSize)
	_, titleH := dc.MeasureString(f.Title)
	dc.DrawString(f.Title, f.PositionX, f.PositionY+titleH)

	// Draw the content and find the rectangle size.
	dc.LoadFontFace(f.Style.FontPath, f.Style.ContentFontSize)
	var W, H float64
	currentDrawY := f.PositionY + f.Style.FrameStrokeWidth + f.Style.TitlePadding + titleH
	for _, content := range f.Contents {
		w, h := dc.MeasureString(content)
		H += h + f.Style.ContentPadding
		if w > W {
			W = w
		}
		currentDrawY += h
		currentDrawY += f.Style.ContentPadding
		dc.DrawString(content, f.PositionX+f.Style.ContentPadding, currentDrawY)
	}

	// Draw the rectangle.
	dc.DrawRectangle(
		f.PositionX,
		f.PositionY+f.Style.TitlePadding+titleH,
		W+(f.Style.ContentPadding*2),
		H+f.Style.ContentPadding,
	)
	dc.Stroke()
}
