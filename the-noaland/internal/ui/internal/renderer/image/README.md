# Terminal image rendering

Render returns inline cells or raw graphics output. UI components use inline
cells in their normal view and reserve a placeholder for raw graphics. Screen
code must not send cursor commands when there is no raw output.

Auto selection uses environment hints, not active protocol probing:
- Ghostty / Kitty: existing Kitty PNG output.
- Other environments (including Windows Terminal, VS Code, and unknown hosts):
  24-bit ANSI colored half-block cells. This is a low-resolution image fallback,
  not native pixel graphics.
- tmux: blocks by default.
- Ghostty through Zellij retains Kitty selection. Older Zellij versions or
  sessions with Kitty graphics disabled may require a manual blocks override.

PowerShell is a shell, not an image protocol. The terminal hosting it must
support ANSI colors and the upper-half-block character for the fallback.
No Sixel support or terminal capability query is implemented in this version.

Set NOALAND_IMAGE_PROTOCOL before starting NoaLand:
auto (default), kitty, blocks, or none.
Unknown values behave as auto. For example in PowerShell:

    $env:NOALAND_IMAGE_PROTOCOL = 'blocks'
    go run ./the-noaland/cmd/noaland

To restore automatic selection:

    Remove-Item Env:NOALAND_IMAGE_PROTOCOL

The fallback composites alpha over RGB(26, 27, 38), assumes two image pixels per
terminal cell, and resamples to the allocated cell area. Protocol encoding and
inline sizing are covered by tests; actual host rendering still needs visual
verification in Ghostty/Zellij and the user's Windows terminal.