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
	width           = 1080
	height          = 720
	mapWidth        = 24
	mapHeight       = 24
	fontPath        = "../assets/fonts/dogica.ttf"
	fontSize        = 24
	textureWidth    = 64
	textureHeight   = 64
	halfPi          = math.Pi * .5
	highlightBorder = 2
	lights          = false
)

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
	{8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 4, 4, 6, 4, 4, 6, 4, 6, 4, 4, 4, 6, 4},
	{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	{8, 0, 3, 3, 0, 0, 0, 0, 0, 8, 8, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6},
	{8, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6},
	{8, 0, 3, 3, 0, 0, 0, 0, 0, 8, 8, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 4, 0, 0, 0, 0, 0, 6, 6, 6, 0, 6, 4, 6},
	{8, 8, 8, 8, 0, 8, 8, 8, 8, 8, 8, 4, 4, 4, 4, 4, 4, 6, 0, 0, 0, 0, 0, 6},
	{7, 7, 7, 7, 0, 7, 7, 7, 7, 0, 8, 0, 8, 0, 8, 0, 8, 4, 0, 4, 0, 6, 0, 6},
	{7, 7, 0, 0, 0, 0, 0, 0, 7, 8, 0, 8, 0, 8, 0, 8, 8, 6, 0, 0, 0, 0, 0, 6},
	{7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 6, 0, 0, 0, 0, 0, 4},
	{7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 6, 0, 6, 0, 6, 0, 6},
	{7, 7, 0, 0, 0, 0, 0, 0, 7, 8, 0, 8, 0, 8, 0, 8, 8, 6, 4, 6, 0, 6, 6, 6},
	{7, 7, 7, 7, 0, 7, 7, 7, 7, 8, 8, 4, 0, 6, 8, 4, 8, 3, 3, 3, 0, 3, 3, 3},
	{2, 2, 2, 2, 0, 2, 2, 2, 2, 4, 6, 4, 0, 0, 6, 0, 6, 3, 0, 0, 0, 0, 0, 3},
	{2, 2, 0, 0, 0, 0, 0, 2, 2, 4, 0, 0, 0, 0, 0, 0, 4, 3, 0, 0, 0, 0, 0, 3},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 4, 0, 0, 0, 0, 0, 0, 4, 3, 0, 0, 0, 0, 0, 3},
	{1, 0, 0, 0, 0, 0, 0, 0, 1, 4, 4, 4, 4, 4, 6, 0, 6, 3, 3, 0, 0, 0, 3, 3},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 1, 2, 2, 2, 6, 6, 0, 0, 5, 0, 5, 0, 5},
	{2, 2, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 0, 2, 2, 0, 5, 0, 5, 0, 0, 0, 5, 5},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 2, 5, 0, 5, 0, 5, 0, 5, 0, 5},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5},
	{2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 2, 5, 0, 5, 0, 5, 0, 5, 0, 5},
	{2, 2, 0, 0, 0, 0, 0, 2, 2, 2, 0, 0, 0, 2, 2, 0, 5, 0, 5, 0, 0, 0, 5, 5},
	{2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 5, 5, 5, 5, 5, 5, 5, 5, 5},
}

func generateTextures(width int32, height int32) *[8][]int32 {
	textures := [8][]int32{}
	for i := range textures {
		textures[i] = make([]int32, width*height)
	}
	//generate some textures
	for x := int32(0); x < width; x++ {
		for y := int32(0); y < height; y++ {
			xorcolor := (x * 0x100 / width) ^ (y * 0x100 / height)
			ycolor := y * 0x100 / height
			xycolor := y*0x80/height + x*0x80/width
			cross := int32(0)
			if x != y && x != width-y {
				cross = 1
			}
			redGradient := int32(0)
			if x%16 != 0 && y%16 != 0 {
				redGradient = 1
			}
			textures[0][width*y+x] = 65536 * 254 * cross                      // flat red texture with black cross
			textures[1][width*y+x] = xycolor + 256*xycolor + 65536*xycolor    // sloped greyscale
			textures[2][width*y+x] = 256*xycolor + 65536*xycolor              // sloped yellow gradient
			textures[3][width*y+x] = xorcolor + 256*xorcolor + 65536*xorcolor // xor greyscale
			textures[4][width*y+x] = 256 * xorcolor                           // xor green
			textures[5][width*y+x] = 65536 * 192 * redGradient                // red bricks
			textures[6][width*y+x] = 65536 * ycolor                           // red gradient
			textures[7][width*y+x] = 128 + 256*128 + 65536*128                // flat grey texture
		}
	}

	return &textures
}

func run() (err error) {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var surface *sdl.Surface
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

	// if surface, err = sdl.CreateRGBSurfaceWithFormat(0, width, height, 32, uint32(sdl.PIXELFORMAT_BGRA32)); err != nil {
	if surface, err = sdl.CreateRGBSurface(0, width, height, 32, 0, 0, 0, 0); err != nil {
		return
	}
	defer surface.Free()

	// Load the font for our text
	if font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		return
	}
	defer font.Close()

	fmt.Println("Starting...")
	running := true

	// game state
	posX := 22.0
	posY := 11.5 // x and y start position

	dirX := -1.0
	dirY := 0.0 // initial direction vector

	planeX := 0.0
	planeY := 0.66 // the 2d raycaster version of camera plane

	time := sdl.GetTicks64() // time of current frame
	oldTime := time          // time of previous frame

	// Which keys are currently being pressed
	keys := map[Direction]interface{}{}
	var hovering struct {
		x     int32
		y     int32
		side  int
		mouse *sdl.MouseMotionEvent
	}

	textures := generateTextures(textureWidth, textureHeight)
	pixels := surface.Pixels()

	// Vertical position of the camera.
	posZ := 0.5 * height

	floorTexture := 3
	ceilingTexture := 5

	for running {
		// surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: width, H: height}, 0)

		// for x := 0; x < 8*texturewidth; x++ {
		// 	for y := 0; y < textureheight; y++ {
		// 		color := textures[x/texturewidth][texturewidth*y+(x%texturewidth)]
		// 		set(x+256, y, color, surface)
		// 	}
		// }

		// rayDir for leftmost ray (x = 0) and rightmost ray (x = w)
		rayDirX0 := dirX - planeX
		rayDirY0 := dirY - planeY
		rayDirX1 := dirX + planeX
		rayDirY1 := dirY + planeY

		// floor casting
		for y := int32(0); y < height; y++ {
			// Current y position compared to the center of the screen (the horizon)
			p := y - height/2

			// Horizontal distance from the camera to the floor for the current row.
			// 0.5 is the z position exactly in the middle between floor and ceiling.
			rowDistance := posZ / float64(p)

			// calculate the real world step vector we have to add for each x (parallel to camera plane)
			// adding step by step avoids multiplications with a weight in the inner loop
			floorStepX := rowDistance * (rayDirX1 - rayDirX0) / width
			floorStepY := rowDistance * (rayDirY1 - rayDirY0) / width

			// real world coordinates of the leftmost column. This will be updated as we step to the right.
			floorX := posX + rowDistance*rayDirX0
			floorY := posY + rowDistance*rayDirY0

			for x := int32(0); x < width; x++ {
				// the cell coord is simply got from the integer parts of floorX and floorY
				cellX := int32(floorX)
				cellY := int32(floorY)

				// get the texture coordinate from the fractional part
				tx := int32(textureWidth*(floorX-float64(cellX))) & (textureWidth - 1)
				ty := int32(textureHeight*(floorY-float64(cellY))) & (textureHeight - 1)

				floorX += floorStepX
				floorY += floorStepY

				dist := max(math.Abs(posX-floorX), math.Abs(posY-floorY))

				// floor
				color := textures[floorTexture][textureWidth*ty+tx]
				color = (color >> 1) & 0x7f7f7f // make a bit darker
				set(x, y, color, dist, surface, pixels)

				// ceiling (symmetrical, at screenHeight - y - 1 instead of y)
				color = textures[ceilingTexture][textureWidth*ty+tx]
				color = (color >> 1) & 0x7f7f7f // make a bit darker
				set(x, height-y-1, color, dist, surface, pixels)
			}
		}

		// walls
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

			// length of ray from one x or y-side to next x or y-side
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

			// texturing calculations
			texNum := worldMap[mapX][mapY] - 1 // 1 subtracted from it so that texture 0 can be used!

			wallX := 0.0 // where exactly the wall was hit
			if side == 0 {
				wallX = posY + perpWallDist*rayDirY
			} else {
				wallX = posX + perpWallDist*rayDirX
			}
			wallX -= math.Floor(wallX)

			// x coordinate on the texture
			texX := int32(wallX * float64(textureWidth))
			if side == 0 && rayDirX > 0 {
				texX = textureWidth - texX - 1
			}
			if side == 1 && rayDirY < 0 {
				texX = textureWidth - texX - 1
			}

			if hovering.mouse != nil &&
				hovering.mouse.X == x &&
				hovering.mouse.Y >= drawStart &&
				hovering.mouse.Y <= drawEnd {
				hovering.x = mapX
				hovering.y = mapY
				hovering.side = side
			}

			step := float64(textureHeight) / float64(lineHeight)

			// Starting texture coordinate
			texPos := float64(drawStart-height/2+lineHeight/2) * step

			for y := drawStart; y < drawEnd; y++ {
				// Cast the texture coordinate to integer, and mask with (textureHeight - 1) in case of overflow
				texY := int32(texPos) & (textureHeight - 1)
				texPos += step
				color := textures[texNum][textureHeight*texY+texX]
				if hovering.side == side && hovering.x == mapX && hovering.y == mapY &&
					((texY < highlightBorder || texY > textureHeight-highlightBorder) || (texX < highlightBorder || texX > textureWidth-highlightBorder)) {
					color = 0xff00ff
				}

				// make color darker for y-sides: R, G and B byte each divided through two with a "shift" and an "and"
				if side == 1 {
					color = (color >> 1) & 0x7f7f7f
				}

				set(x, y, color, perpWallDist, surface, pixels)
			}
		}

		// draw a map
		for x := int32(0); x < mapWidth; x++ {
			for y := int32(0); y < mapHeight; y++ {
				color := int32(0xffffff)
				if x == int32(posX) && y == int32(posY) {
					color = 0x10fdff
				} else if worldMap[x][y] == 0 {
					color = 0
				}

				scaledX := x * 4
				scaledY := y * 4
				for sx := int32(0); sx < 4; sx++ {
					sy := scaledY + 2 + sx
					set(scaledX+2, sy, color, 0, surface, pixels)
					set(scaledX+3, sy, color, 0, surface, pixels)
					set(scaledX+4, sy, color, 0, surface, pixels)
					set(scaledX+5, sy, color, 0, surface, pixels)
				}
			}
		}

		// timing for input and FPS counter
		oldTime = time
		time = sdl.GetTicks64()
		frameTime := float64(time-oldTime) / 1000.0 // frameTime is the time this frame has taken, in milliseconds

		// output the rendered surface
		surfaceTexture, _ := renderer.CreateTextureFromSurface(surface)
		renderer.Copy(surfaceTexture, nil, &sdl.Rect{X: 1, Y: 1, W: width, H: height})
		surfaceTexture.Destroy()

		// render fps counter
		fps, _ := font.RenderUTF8Blended(fmt.Sprintf("%.0f", 1.0/frameTime), white)
		fpsTexture, _ := renderer.CreateTextureFromSurface(fps)
		renderer.Copy(fpsTexture, nil, &sdl.Rect{X: width - fps.W, Y: 2, W: fps.W, H: fps.H})
		fpsTexture.Destroy()
		fps.Free()

		moveSpeed := frameTime * 5.0 // the constant value is in squares/second
		rotSpeed := frameTime * 3.0  // the constant value is in radians/second

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

		// Process events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Bye, Felicia")
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Printf("Mouse %d moved by %d %d at %d %d\n", t.Which, t.XRel, t.YRel, t.X, t.Y)
				hovering.mouse = t
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

func set(x int32, y int32, color int32, dist float64, surface *sdl.Surface, pixels []byte) {
	i := y*surface.Pitch + x*int32(surface.Format.BytesPerPixel)

	var r, g, b byte
	if lights {
		b = byte(color >> 0 & 0xff)
		g = byte(color >> 8 & 0xff)
		r = byte(color >> 16 & 0xff)
	} else {
		scale := 1 - min(.3*math.Sqrt(dist), 1)
		b = byte(float64(color>>0&0xff) * scale)
		g = byte(float64(color>>8&0xff) * scale)
		r = byte(float64(color>>16&0xff) * scale)

	}

	pixels[i] = b
	pixels[i+1] = g
	pixels[i+2] = r
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
