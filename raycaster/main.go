package main

import (
	"fmt"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Direction = int

const (
	Forward Direction = iota
	Backward
	Left
	Right
	StrafeLeft
	StrafeRight
)

const (
	width     = 1080
	height    = 720
	mapWidth  = 24
	mapHeight = 24
	fontPath  = "../assets/fonts/dogica.ttf"
	fontSize  = 24
)

const halfPi = math.Pi * .5

var white = sdl.Color{R: 255, G: 255, B: 255, A: 255}
var black = sdl.Color{R: 0, G: 0, B: 0, A: 255}
var red = sdl.Color{R: 255, G: 0, B: 0, A: 255}
var green = sdl.Color{R: 0, G: 255, B: 0, A: 255}
var blue = sdl.Color{R: 0, G: 0, B: 255, A: 255}
var yellow = sdl.Color{R: 255, G: 255, B: 0, A: 255}
var sky = sdl.Color{R: 0, G: 175, B: 240, A: 255}

func pointp(X int32, Y int32) *sdl.Point {
	return &sdl.Point{X: X, Y: Y}
}

var worldMap = [mapWidth][mapHeight]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

func run() (err error) {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var font *ttf.Font

	if err = ttf.Init(); err != nil {
		return
	}
	defer ttf.Quit()

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return
	}
	defer sdl.Quit()

	if window, err = sdl.CreateWindow("Raycaster", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN); err != nil {
		return
	}
	defer window.Destroy()

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return
	}
	defer renderer.Destroy()

	// Load the font for our text
	if font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		return
	}
	defer font.Close()

	fmt.Println("Starting...")
	running := true

	// game state
	posX := 22.0
	posY := 12.0 // x and y start position

	dirX := -1.0
	dirY := 0.0 // initial direction vector

	planeX := 0.0
	planeY := 0.66 // the 2d raycaster version of camera plane

	time := sdl.GetTicks64() // time of current frame
	oldTime := time - 1      // time of previous frame

	keys := map[Direction]interface{}{}

	for running {
		renderer.SetDrawColor(sky.R, sky.G, sky.B, sky.A)
		renderer.Clear()

		renderer.SetDrawColor(black.R, black.G, black.B, black.A)
		renderer.FillRect(&sdl.Rect{X: 0, Y: height / 2, W: width, H: height})

		for x := int32(0); x < width; x++ {
			// calculate ray position and direction
			cameraX := float64(2*x)/float64(width) - 1 // x-coordinate in camera space
			rayDirX := dirX + planeX*cameraX
			rayDirY := dirY + planeY*cameraX

			// which box of the map we're in
			mapX := int32(posX)
			mapY := int32(posY)

			// length of ray from current position to next x or y-side
			sideDistX := 0.0
			sideDistY := 0.0

			//l ength of ray from one x or y-side to next x or y-side
			deltaDistX := math.Abs(1 / rayDirX)

			if rayDirX == 0 {
				deltaDistX = math.MaxFloat64
			}

			deltaDistY := math.Abs(1 / rayDirY)

			if rayDirY == 0 {
				deltaDistY = math.MaxFloat64
			}

			perpWallDist := 0.0

			// what direction to step in x or y-direction (either +1 or -1)
			stepX := int32(0)
			stepY := int32(0)

			hit := 0  // was there a wall hit?
			side := 0 // was a NS or a EW wall hit?

			// calculate step and initial sideDist
			if rayDirX < 0 {
				stepX = -1
				sideDistX = (posX - float64(mapX)) * deltaDistX
			} else {
				stepX = 1
				sideDistX = (float64(mapX) + 1.0 - posX) * deltaDistX
			}

			if rayDirY < 0 {
				stepY = -1
				sideDistY = (posY - float64(mapY)) * deltaDistY
			} else {
				stepY = 1
				sideDistY = (float64(mapY) + 1.0 - posY) * deltaDistY
			}

			for hit == 0 {
				// jump to next map square, either in x-direction, or in y-direction
				if sideDistX < sideDistY {
					sideDistX += deltaDistX
					mapX += stepX
					side = 0
				} else {
					sideDistY += deltaDistY
					mapY += stepY
					side = 1
				}
				// Check if ray has hit a wall
				if worldMap[mapX][mapY] > 0 {
					hit = 1
				}
			}

			// Calculate distance projected on camera direction (Euclidean distance would give fisheye effect!)
			if side == 0 {
				perpWallDist = (sideDistX - deltaDistX)
			} else {
				perpWallDist = (sideDistY - deltaDistY)
			}

			// Calculate height of line to draw on screen
			lineHeight := int32(height / perpWallDist)

			// calculate lowest and highest pixel to fill in current stripe
			drawStart := -lineHeight/2 + height/2
			if drawStart < 0 {
				drawStart = 0
			}
			drawEnd := lineHeight/2 + height/2
			if drawEnd >= height {
				drawEnd = height - 1
			}

			var color sdl.Color
			switch worldMap[mapX][mapY] {
			case 1:
				color = red
			case 2:
				color = green
			case 3:
				color = blue
			case 4:
				color = white
			default:
				color = yellow
			}

			// give x and y sides different brightness
			if side == 1 {
				color = sdl.Color{R: color.R / 2, G: color.G / 2, B: color.B / 2, A: color.A / 2}
			}

			// draw the pixels of the stripe as a vertical line
			renderer.SetDrawColor(color.R, color.G, color.B, color.A)
			renderer.DrawLine(x, drawStart, x, drawEnd)
		}

		// timing for input and FPS counter
		oldTime = time - 1
		time = sdl.GetTicks64()
		frameTime := float64(time-oldTime) / 1000.0 // frameTime is the time this frame has taken, in milliseconds

		fps, _ := font.RenderUTF8Blended(fmt.Sprintf("%.0f", (1.0/float64(frameTime))), white)
		fpsTexture, _ := renderer.CreateTextureFromSurface(fps)
		renderer.Copy(fpsTexture, nil, &sdl.Rect{X: 2, Y: 2, W: fps.W, H: fps.H})

		moveSpeed := float64(frameTime) * 5.0 // the constant value is in squares/second
		rotSpeed := float64(frameTime) * 3.0  // the constant value is in radians/second

		if _, exists := keys[Forward]; exists {
			if worldMap[int(posX+dirX*moveSpeed)][int(posY)] == 0 {
				posX += dirX * moveSpeed
			}
			if worldMap[int(posX)][int(posY+dirY*moveSpeed)] == 0 {
				posY += dirY * moveSpeed
			}
		}

		if _, exists := keys[Backward]; exists {
			if worldMap[int(posX-dirX*moveSpeed)][int(posY)] == 0 {
				posX -= dirX * moveSpeed
			}
			if worldMap[int(posX)][int(posY-dirY*moveSpeed)] == 0 {
				posY -= dirY * moveSpeed
			}
		}

		if _, exists := keys[StrafeRight]; exists {
			perpDirX := dirX*math.Cos(-halfPi) - dirY*math.Sin(-halfPi)
			perpDirY := dirX*math.Sin(-halfPi) + dirY*math.Cos(-halfPi)
			if worldMap[int(posX+perpDirX*moveSpeed)][int(posY)] == 0 {
				posX += perpDirX * moveSpeed
			}
			if worldMap[int(posX)][int(posY+perpDirY*moveSpeed)] == 0 {
				posY += perpDirY * moveSpeed
			}
		}

		if _, exists := keys[StrafeLeft]; exists {
			perpDirX := dirX*math.Cos(halfPi) - dirY*math.Sin(halfPi)
			perpDirY := dirX*math.Sin(halfPi) + dirY*math.Cos(halfPi)
			if worldMap[int(posX+perpDirX*moveSpeed)][int(posY)] == 0 {
				posX += perpDirX * moveSpeed
			}
			if worldMap[int(posX)][int(posY+perpDirY*moveSpeed)] == 0 {
				posY += perpDirY * moveSpeed
			}
		}

		if _, exists := keys[Right]; exists {
			oldDirX := dirX
			dirX = dirX*math.Cos(-rotSpeed) - dirY*math.Sin(-rotSpeed)
			dirY = oldDirX*math.Sin(-rotSpeed) + dirY*math.Cos(-rotSpeed)
			oldPlaneX := planeX
			planeX = planeX*math.Cos(-rotSpeed) - planeY*math.Sin(-rotSpeed)
			planeY = oldPlaneX*math.Sin(-rotSpeed) + planeY*math.Cos(-rotSpeed)
		}

		if _, exists := keys[Left]; exists {
			oldDirX := dirX
			dirX = dirX*math.Cos(rotSpeed) - dirY*math.Sin(rotSpeed)
			dirY = oldDirX*math.Sin(rotSpeed) + dirY*math.Cos(rotSpeed)
			oldPlaneX := planeX
			planeX = planeX*math.Cos(rotSpeed) - planeY*math.Sin(rotSpeed)
			planeY = oldPlaneX*math.Sin(rotSpeed) + planeY*math.Cos(rotSpeed)
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Bye, Felicia")
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Printf("Mouse %d moved by %d %d at %d %d\n", t.Which, t.XRel, t.YRel, t.X, t.Y)
			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED {
					fmt.Printf("Mouse %d button %d pressed at %d %d\n", t.Which, t.Button, t.X, t.Y)
				} else {
					fmt.Printf("Mouse %d button %d released at %d %d\n", t.Which, t.Button, t.X, t.Y)
				}
			case *sdl.MouseWheelEvent:
				if t.X != 0 {
					fmt.Printf("Mouse %d wheel scrolled horizontally by %d\n", t.Which, t.X)
				} else {
					fmt.Printf("Mouse %d wheel scrolled vertically by %d\n", t.Which, t.Y)
				}
			case *sdl.KeyboardEvent:
				switch t.Keysym.Sym {
				case sdl.K_UP, sdl.K_w:
					if t.State == sdl.RELEASED {
						delete(keys, Forward)
					} else {
						keys[Forward] = nil
					}
				case sdl.K_DOWN, sdl.K_s:
					if t.State == sdl.RELEASED {
						delete(keys, Backward)
					} else {
						keys[Backward] = nil
					}
				case sdl.K_RIGHT, sdl.K_d:
					if t.State == sdl.RELEASED {
						delete(keys, Right)
					} else {
						keys[Right] = nil
					}
				case sdl.K_LEFT, sdl.K_a:
					if t.State == sdl.RELEASED {
						delete(keys, Left)
					} else {
						keys[Left] = nil
					}
				case sdl.K_q:
					if t.State == sdl.RELEASED {
						delete(keys, StrafeLeft)
					} else {
						keys[StrafeLeft] = nil
					}
				case sdl.K_e:
					if t.State == sdl.RELEASED {
						delete(keys, StrafeRight)
					} else {
						keys[StrafeRight] = nil
					}
				}
			}
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
