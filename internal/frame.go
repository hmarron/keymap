package internal

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/fogleman/gg"
)

type FrameStyle struct {
	FrameStrokeWidth float64
	TitleFontSize    float64
	ContentFontSize  float64
	TitlePadding     float64
	ContentPadding   float64
	FontPath         string
	FramePaddingX    float64
	FramePaddingY    float64
}

type Frame struct {
	Title    string
	Contents []string
	Style    FrameStyle
}

// getCalculatedDimensions returns dimensions not including padding.
func (f *Frame) getCalculatedDimensions(dc *gg.Context) (float64, float64) {
	f.normalizeContent()

	// Get title height.
	dc.LoadFontFace(f.Style.FontPath, f.Style.TitleFontSize)
	_, titleH := dc.MeasureString(f.Title)

	dc.LoadFontFace(f.Style.FontPath, f.Style.ContentFontSize)
	var W, H float64
	currentDrawY := f.Style.FrameStrokeWidth + f.Style.TitlePadding + titleH
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

// GetCalculatedDimensions returns dimensions including padding.
func (f *Frame) GetCalculatedDimensions(dc *gg.Context) (float64, float64) {
	w, h := f.getCalculatedDimensions(dc)
	return (w + (f.Style.FramePaddingX * 2)),
		(h + (f.Style.FramePaddingY * 2))
}

func (f *Frame) DebugDraw(dc *gg.Context, x, y float64) {
	// Draw standard
	f.Draw(dc, x, y)

	// Draw bounding rect including padding.
	w, h := f.getCalculatedDimensions(dc)
	dc.DrawRectangle(
		x,
		y,
		w+(f.Style.FramePaddingX*2),
		h+(f.Style.FramePaddingY*2),
	)
	dc.SetColor(color.RGBA{255, 0, 0, 255})
	dc.Stroke()
}

func (f *Frame) Draw(dc *gg.Context, x, y float64) {
	// Set default styles
	dc.SetColor(color.White)
	dc.SetLineWidth(f.Style.FrameStrokeWidth)

	// Draw the title
	dc.LoadFontFace(f.Style.FontPath, f.Style.TitleFontSize)
	_, titleH := dc.MeasureString(f.Title)
	dc.DrawString(
		f.Title,
		x+f.Style.FramePaddingX,
		y+titleH+f.Style.FramePaddingY,
	)

	// Draw the content and find the rectangle size.
	dc.LoadFontFace(f.Style.FontPath, f.Style.ContentFontSize)

	// Thes vars track the max heigh and widht of all content.
	// used to draw the box around it later.
	var contentW, contentH float64
	// Start drawing just under the title.
	currentDrawY := y + f.Style.FrameStrokeWidth + f.Style.TitlePadding + titleH + f.Style.FramePaddingY
	for _, content := range f.Contents {
		lineW, lineH := dc.MeasureString(content)

		// Always update contentH to add the height of this line.
		contentH += lineH + f.Style.ContentPadding

		// If this is now the new longest line, updated contentW.
		if lineW > contentW {
			contentW = lineW
		}

		// Move the next line.
		currentDrawY += lineH
		currentDrawY += f.Style.ContentPadding

		// Draw the line.
		dc.DrawString(
			content,
			x+f.Style.ContentPadding+f.Style.FramePaddingX,
			currentDrawY,
		)
	}

	// Draw the rectangle.
	dc.DrawRectangle(
		x+f.Style.FramePaddingX,
		y+f.Style.TitlePadding+titleH+f.Style.FramePaddingY,
		contentW+(f.Style.ContentPadding*2),
		contentH+f.Style.ContentPadding,
	)
	dc.Stroke()
}

func (f *Frame) normalizeContent() {
	maxKeyLength := 0
	// Get the longest key length.
	for _, content := range f.Contents {
		if strings.Contains(content, ":") {
			key := strings.SplitN(content, ":", 2)[0]
			if len(key) > maxKeyLength {
				maxKeyLength = len(key)
			}
		}
	}
	// Go over all lines and make sure that the first : is in the same col for every line.
	// Normalize the content by padding keys with spaces.
	for i, content := range f.Contents {
		if strings.Contains(content, ":") {
			parts := strings.SplitN(content, ":", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				paddedKey := fmt.Sprintf("%-*s", maxKeyLength+1, key)
				f.Contents[i] = paddedKey + ":" + value
			}
		}
	}

}
