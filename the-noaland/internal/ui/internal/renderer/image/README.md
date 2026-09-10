# Terminal image rendering

Render returns inline cells or raw graphics output. UI components use inline
cells in their normal view and reserve a placeholder for raw graphics. Screen
code must not send cursor commands when there is no raw output.

Auto selection uses environment hints, not active protocol probing:
- Ghostty / Kitty: existing Kitty PNG output.
- Windows Terminal (WT_SESSION): Sixel, requiring a version with Sixel support.
- Other environments (including VS Code and unknown hosts):
  24-bit ANSI colored half-block cells. This is a low-resolution image fallback,
  not native pixel graphics.
- tmux: blocks by default.
- Ghostty through Zellij retains Kitty selection. Older Zellij versions or
  sessions with Kitty graphics disabled may require a manual blocks override.

PowerShell is a shell, not an image protocol. The terminal hosting it must
support ANSI colors and the upper-half-block character for the fallback.
Sixel is supported; active terminal capability queries are not yet implemented.

Set NOALAND_IMAGE_PROTOCOL before starting NoaLand:
auto (default), kitty, sixel, blocks, or none.
Unknown values behave as auto. For example in PowerShell:

    $env:NOALAND_IMAGE_PROTOCOL = 'blocks'
    go run ./the-noaland/cmd/noaland

To restore automatic selection:

    Remove-Item Env:NOALAND_IMAGE_PROTOCOL

The fallback composites alpha over RGB(26, 27, 38), assumes two image pixels per
terminal cell, and resamples to the allocated cell area. Protocol encoding and
inline sizing are covered by tests; actual host rendering still needs visual
verification in Ghostty/Zellij and the user's Windows terminal.
## Native Windows images

Yazi reports Sixel for the user's Windows session. Select the same backend:

    $env:NOALAND_IMAGE_PROTOCOL = 'sixel'
    go run ./the-noaland/cmd/noaland

This replaces any earlier persistent blocks override. Auto also selects Sixel
when WT_SESSION is present (outside multiplexers). Windows Terminal requires
version 1.22.10352.0 or newer for Sixel support.

Sixel uses a quantized pixel image rather than text blocks. Its raster must be
sized in pixels, while the layout allocates terminal cells. When terminal pixel
dimensions are unavailable, this implementation estimates 8x16 pixels per cell.
If the image is too small or extends past its placeholder, set the actual cell
width and height for the terminal's font, zoom, and DPI, for example:

    $env:NOALAND_IMAGE_CELL_SIZE = '10x20'

Accepted cell dimensions are 1..128 by 1..256. Invalid values use 8x16.
The estimate is not automatic font/DPI detection. Sixel output still needs
visual verification on the user's terminal, including resizing and refresh.

References:
- https://yazi-rs.github.io/docs/image-preview/
- https://devblogs.microsoft.com/commandline/windows-terminal-preview-1-22-release/
## Defaults without environment setup

No environment variables are required: auto protocol selection is the default.
Windows Terminal selects Sixel, and Ghostty/Kitty select Kitty. Unknown terminals
retain the blocks fallback; an operating system alone does not imply support.

Cell size precedence:
1. Valid NOALAND_IMAGE_CELL_SIZE override.
2. Visible native Windows console font metrics, if available.
3. 8x16 pixel estimate.

Windows Terminal, remote sessions, and hidden ConPTY consoles deliberately do
not use legacy console font metrics, which may differ from the terminal font.
This does not automatically detect Windows Terminal font zoom or DPI.

Remove earlier overrides once to use defaults in the current PowerShell session:

    Remove-Item Env:NOALAND_IMAGE_PROTOCOL, Env:NOALAND_IMAGE_CELL_SIZE -ErrorAction SilentlyContinue

If overrides were added to a PowerShell profile or persistent user environment,
remove them there as well. NoaLand never changes the user's environment settings.