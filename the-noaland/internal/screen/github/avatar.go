package github

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"strings"
)

const kittyChunkSize = 4096

func avatarPlaceholder(
	width int,
	height int,
) string {
	line := strings.Repeat(" ", width)

	lines := make([]string, height)

	for i := range lines {
		lines[i] = line
	}

	return strings.Join(lines, "\n")
}

func renderKittyImage(
	img image.Image,
	cols int,
	rows int,
) string {
	if img == nil {
		return ""
	}

	var pngBuf bytes.Buffer

	if err := png.Encode(&pngBuf, img); err != nil {
		return ""
	}

	data := base64.StdEncoding.EncodeToString(
		pngBuf.Bytes(),
	)

	var out strings.Builder

	first := true

	for len(data) > 0 {
		n := min(len(data), kittyChunkSize)

		chunk := data[:n]
		data = data[n:]

		more := len(data) > 0

		m := 0
		if more {
			m = 1
		}

		if first {
			fmt.Fprintf(
				&out,
				"\x1b_Ga=T,f=100,c=%d,r=%d,q=2,m=%d;%s\x1b\\",
				cols,
				rows,
				m,
				chunk,
			)

			first = false
		} else {
			fmt.Fprintf(
				&out,
				"\x1b_Gm=%d;%s\x1b\\",
				m,
				chunk,
			)
		}
	}

	return out.String()
}

func moveCursor(
	row int,
	col int,
) string {
	return fmt.Sprintf(
		"\x1b[%d;%dH",
		row,
		col,
	)
}
