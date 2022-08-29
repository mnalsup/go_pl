package main

import (
	"image"
	"math/rand"
	"os"
	"time"

	"image/color"
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const GAME_WIDTH = 1020
const GAME_HEIGHT = 1020

const DOT_HEIGHT = 10
const DOT_WIDTH = 10

func loadImage(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game Of Life",
		Bounds: pixel.R(0, 0, GAME_WIDTH, GAME_HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	dot, err := loadImage("assets/dot.png")
	if err != nil {
		panic(err)
	}

	gridHeight := (GAME_HEIGHT / DOT_HEIGHT)
	gridWidth := (GAME_WIDTH / DOT_WIDTH)

	gridA := make([][]bool, gridWidth)
	for i := range gridA {
		gridA[i] = make([]bool, gridHeight)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range gridA {
		for k := range gridA[i] {
			if n := r.Intn(10); n < 2 {
				gridA[i][k] = true
			}
		}
	}
	gridB := make([][]bool, gridWidth)
	for i := range gridB {
		gridB[i] = make([]bool, gridHeight)
	}

	currentGrid := &gridA
	nextGrid := &gridB

	for !win.Closed() {

		win.Clear(color.Black)

		for i := range *currentGrid {
			for k := range (*currentGrid)[i] {
				if (*currentGrid)[i][k] {
					dotSprite := pixel.NewSprite(dot, dot.Bounds())
					dotSprite.Draw(win, pixel.IM.Moved(pixel.V(float64(i*DOT_WIDTH), float64(k*DOT_HEIGHT))))
				}
			}
		}

		win.Update()

		for i := range *currentGrid {
			for k := range (*currentGrid)[i] {
				neighbours := 0
				// right
				if i < gridWidth-1 && (*currentGrid)[i+1][k] {
					neighbours += 1
				}
				// up
				if k < gridHeight-1 && (*currentGrid)[i][k+1] {
					neighbours += 1
				}
				// left
				if i > 0 && (*currentGrid)[i-1][k] {
					neighbours += 1
				}
				// down
				if k > 0 && (*currentGrid)[i][k-1] {
					neighbours += 1
				}
				// right-up
				if i < gridWidth-1 && k < gridHeight-1 && (*currentGrid)[i+1][k+1] {
					neighbours += 1
				}
				// right-down
				if i < gridWidth-1 && k > 0 && (*currentGrid)[i+1][k-1] {
					neighbours += 1
				}
				// left-up
				if i > 0 && k < gridHeight-1 && (*currentGrid)[i-1][k+1] {
					neighbours += 1
				}
				// left-down
				if i > 0 && k > 0 && (*currentGrid)[i-1][k-1] {
					neighbours += 1
				}

				// survive
				if (*currentGrid)[i][k] && (neighbours == 3 || neighbours == 2) {
					(*nextGrid)[i][k] = true
				} else if !(*currentGrid)[i][k] && neighbours == 3 {
					// reproduce
					(*nextGrid)[i][k] = true
				} else {
					// dead
					(*nextGrid)[i][k] = false
				}
			}
		}
		currentGrid, nextGrid = nextGrid, currentGrid

	}
}

func main() {
	pixelgl.Run(run)
}
