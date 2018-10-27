package main

import (
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "LEGO TETRIS",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true, // synchronizacja z predkoscia odswiezania monitora
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	positionX := win.Bounds().Center().X // pozycja klocka
	positionY := win.Bounds().Center().Y

	lastTime := time.Now()
	leftRightTime := 0.0 // zmienna do zliczania czasu
	moveLeftOrRight := 0

	pic, err := loadPicture("brickBlack.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds()) //obiekt sluzacy do wyswietlenia obrazka

	for !win.Closed() {
		dt := time.Since(lastTime).Seconds() // czas ktory uplynal od poprzedniej klatki
		lastTime = time.Now()

		leftRightTime += dt
		moveLeftOrRight = 0

		win.Clear(colornames.Lightsteelblue) //wyczyszczenie okna przed wyswietleniem

		if win.Pressed(pixelgl.KeyLeft) {
			moveLeftOrRight = -1
			//positionX--

		}
		if win.Pressed(pixelgl.KeyRight) {
			moveLeftOrRight = 1
			//positionX++

		}
		if win.Pressed(pixelgl.KeyDown) {
			positionY--

		}
		if win.Pressed(pixelgl.KeyUp) {
			positionY++

		}

		if leftRightTime >= 0.1 && moveLeftOrRight != 0 {
			leftRightTime = 0.0
			positionX += float64(moveLeftOrRight)
		}
		sprite.Draw(win, pixel.IM.Moved(pixel.V(positionX, positionY)))

		win.Update() // odswiezenie okna
	}

}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // defer = finnaly
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func main() {
	pixelgl.Run(run)
}
