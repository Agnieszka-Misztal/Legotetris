package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "LEGO TETRIS",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	positionX := win.Bounds().Center().X // pozycja klocka
	positionY := win.Bounds().Center().Y

	pic, err := loadPicture("brickBlack.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds()) //obiekt sluzacy do wyswietlenia obrazka

	for !win.Closed() {
		win.Clear(colornames.Lightsteelblue) //wyczyszczenie okna przed wyswietleniem

		if win.Pressed(pixelgl.KeyLeft) {
			positionX--

		}
		if win.Pressed(pixelgl.KeyRight) {
			positionX++

		}
		if win.Pressed(pixelgl.KeyDown) {
			positionY--

		}
		if win.Pressed(pixelgl.KeyUp) {
			positionY++

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
