package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	width  = 800
	height = 800
)

var white = sdl.Color{R: 255, G: 255, B: 255, A: 255}
var black = sdl.Color{R: 0, G: 0, B: 0, A: 255}
var red = sdl.Color{R: 255, G: 0, B: 0, A: 255}
var green = sdl.Color{R: 0, G: 255, B: 0, A: 255}

func pointp(X int32, Y int32) *sdl.Point {
	return &sdl.Point{X: X, Y: Y}
}

func run() (err error) {
	var window *sdl.Window
	var renderer *sdl.Renderer

	if err = ttf.Init(); err != nil {
		return
	}
	defer ttf.Quit()

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return
	}
	defer sdl.Quit()

	if window, err = sdl.CreateWindow("Lights", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN); err != nil {
		return
	}
	defer window.Destroy()

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return
	}

	fmt.Println("Starting...")
	running := true

	var str string
	var hovered *sdl.Point = nil
	var marked *sdl.Point = nil
	gridSize := int32(100)
	squareColor := red

	rows := height / gridSize
	cols := width / gridSize
	grid := make([][]int, cols)
	for i := int32(0); i < cols; i++ {
		grid[i] = make([]int, rows)
	}

	for running {
		renderer.SetDrawColor(black.R, black.G, black.B, 255)
		renderer.Clear()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Bye, Felicia")
				str = "Bye, Felicia"
				running = false
			case *sdl.MouseMotionEvent:
				str = fmt.Sprintf("Mouse %d moved by %d %d at %d %d", t.Which, t.XRel, t.YRel, t.X, t.Y)
				fmt.Println(str)
				hovered = pointp(t.X, t.Y)
			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED {
					squareColor = green
					str = fmt.Sprintf("Mouse %d button %d pressed at %d %d", t.Which, t.Button, t.X, t.Y)
					fmt.Println(str)
				} else {
					str = fmt.Sprintf("Mouse %d button %d released at %d %d", t.Which, t.Button, t.X, t.Y)
					fmt.Println(str)
					squareColor = red
					marked = pointp(t.X, t.Y)
				}
			case *sdl.MouseWheelEvent:
				if t.X != 0 {
					str = fmt.Sprintf("Mouse %d wheel scrolled horizontally by %d", t.Which, t.X)
					fmt.Println(str)
				} else {
					str = fmt.Sprintf("Mouse %d wheel scrolled vertically by %d", t.Which, t.Y)
				}
			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym
				keys := ""

				// Modifier keys
				switch t.Keysym.Mod {
				case sdl.KMOD_LALT:
					keys += "Left Alt"
				case sdl.KMOD_LCTRL:
					keys += "Left Control"
				case sdl.KMOD_LSHIFT:
					keys += "Left Shift"
				case sdl.KMOD_LGUI:
					keys += "Left Meta or Windows key"
				case sdl.KMOD_RALT:
					keys += "Right Alt"
				case sdl.KMOD_RCTRL:
					keys += "Right Control"
				case sdl.KMOD_RSHIFT:
					keys += "Right Shift"
				case sdl.KMOD_RGUI:
					keys += "Right Meta or Windows key"
				case sdl.KMOD_NUM:
					keys += "Num Lock"
				case sdl.KMOD_CAPS:
					keys += "Caps Lock"
				case sdl.KMOD_MODE:
					keys += "AltGr Key"
				}

				if keyCode < 10000 {
					if keys != "" {
						keys += " + "
					}

					// If the key is held down, this will fire
					if t.Repeat > 0 {
						keys += string(keyCode) + " repeating"
					} else {
						if t.State == sdl.RELEASED {
							keys += string(keyCode) + " released"
						} else if t.State == sdl.PRESSED {
							keys += string(keyCode) + " pressed"
						}
					}

				}

				if keys != "" {
					fmt.Println(keys)
				}
			}
		}

		renderer.SetDrawColor(white.R, white.G, white.B, 255)

		for y := int32(0); y < height; y += gridSize {
			renderer.DrawLine(0, y, width, y)
			for x := int32(0); x < width; x += gridSize {
				renderer.DrawLine(x, 0, x, height)
			}
		}

		if marked != nil {
			col := marked.X / gridSize
			row := marked.Y / gridSize
			grid[col][row] = (grid[col][row] + 1) % 2

			if col > 0 {
				grid[col-1][row] = (grid[col-1][row] + 1) % 2
			}
			if col < cols-1 {
				grid[col+1][row] = (grid[col+1][row] + 1) % 2
			}
			if row > 0 {
				grid[col][row-1] = (grid[col][row-1] + 1) % 2
			}
			if row < rows-1 {
				grid[col][row+1] = (grid[col][row+1] + 1) % 2
			}

			marked = nil
		}

		renderer.SetDrawColor(white.R, white.G, white.B, 127)
		for col := range grid {
			for row := range grid[col] {
				if grid[col][row] == 1 {
					renderer.FillRect(&sdl.Rect{X: int32(col)*gridSize + 1, Y: int32(row)*gridSize + 1, W: gridSize - 1, H: gridSize - 1})
				}
			}
		}

		if hovered != nil {
			renderer.SetDrawColor(squareColor.R, squareColor.G, squareColor.B, 127)
			rounded := sdl.Point{X: hovered.X - hovered.X%gridSize + 1, Y: hovered.Y - hovered.Y%gridSize + 1}
			renderer.FillRect(&sdl.Rect{X: rounded.X, Y: rounded.Y, W: gridSize - 1, H: gridSize - 1})
		}

		renderer.Present()

		sdl.Delay(16)
	}

	return
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
