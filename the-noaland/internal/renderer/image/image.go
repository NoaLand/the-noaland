package image

import "image"

// Render encodes an image for terminal output without moving the cursor.
// It currently uses Kitty; terminal capability selection is not implemented yet.
func Render(img image.Image, cols, rows int) string {
	return renderKitty(img, cols, rows)
}
