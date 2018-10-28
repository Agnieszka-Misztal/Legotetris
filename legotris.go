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

	//var greed [20][10]int

	figures := [7][4]int{
		{1, 3, 5, 7}, // I  0 1
		{2, 4, 5, 7}, // Z  2 3
		{3, 5, 4, 6}, // S  4 5
		{1, 3, 2, 5}, // T  6 7
		{2, 3, 5, 7}, // L
		{3, 5, 7, 6}, // J
		{2, 3, 4, 5}, // O
	}

	var figure [4]pixel.Vec //4 pozycje klocka, vectory

	//stworzenie klocka
	for i := 0; i < 4; i++ {
		figure[i].X = float64(figures[3][i] % 2)    //ustawienie x na 0 lub 1
		figure[i].Y = float64(figures[3][i]/2 + 16) //ustawienie y od 0 do 3
	}

	cfg := pixelgl.WindowConfig{
		Title:  "LEGO TETRIS",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true, // synchronizacja z predkoscia odswiezania monitora
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//positionX := win.Bounds().Center().X // pozycja klocka
	positionY := win.Bounds().Center().Y

	lastTime := time.Now()
	leftRightTime := 0.0 // zmienna do zliczania czasu
	moveLeftOrRight := 0
	moveDwonTime := 0.0

	pic, err := loadPicture("brickBlack.png")
	if err != nil {
		panic(err)
	}

	block := pixel.NewSprite(pic, pic.Bounds()) //obiekt sluzacy do wyswietlenia obrazka

	for !win.Closed() {
		dt := time.Since(lastTime).Seconds() // czas ktory uplynal od poprzedniej klatki
		lastTime = time.Now()

		leftRightTime += dt
		moveDwonTime += dt
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
		if win.JustPressed(pixelgl.KeyUp) { // just pressed zadziala tylko raz
			// pobranie lokalizacji srodka do obracania klocka, zawsze drugi element
			centerBlockX := figure[1].X
			centerBlockY := figure[1].Y

			for i := 0; i < 4; i++ {
				// wyznaczenie wkatora kierunku dla kazdego elemntu
				x := figure[i].X - centerBlockX
				y := figure[i].Y - centerBlockY

				//obrot wektora kierunku o 90 st w pawo, przemnozenie przez macierz obrotu
				//nowe wpolrzedne
				// mnozenie macierzy obroty przez wektor
				x1 := 0.0*x + 1.0*y
				y1 := -1.0*x + 0.0*y

				//ustalenie polozenia klocka na planszy (dodanie do pozycji srodka)
				figure[i].X = x1 + centerBlockX
				figure[i].Y = y1 + centerBlockY

			}

		}

		if moveDwonTime >= 1.0 {
			moveDwonTime = 0.0
			for i := 0; i < 4; i++ {
				figure[i].Y--

			}
		}

		if leftRightTime >= 0.1 && moveLeftOrRight != 0 {
			leftRightTime = 0.0
			for i := 0; i < 4; i++ {
				figure[i].X += float64(moveLeftOrRight)

			}

		}

		// for y := 0; y < 20; y++ {
		// 	for x := 0; x < 10; x++ {

		// 		block.Draw(win, pixel.IM.Moved(pixel.V(float64(x*32+16+400), float64(y*25+16+50))))
		// 	}
		// }

		//rysowanie klocka
		for i := 0; i < 4; i++ {
			block.Draw(win, pixel.IM.Moved(pixel.V(figure[i].X*32.0+16.0, figure[i].Y*25+16.0)))
		}

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

func checkCollision(grid [20][10]int, figure [4]pixel.Vec) bool {
	for i := 0; i < 4; i++ {

		//sprawdzanie czy klocek nie wychodzi poza boczne sciany
		if figure[i].X < 0 || figure[i].X > 9 {
			return true
		}

		//sprawdzanie czy klocek nie wyszedl za nisko, y
		if figure[i].Y < 0 {
			return true
		}

		//sprawdzanie czy na miejscu klocka jest juz inny klocek
		if grid[int(figure[i].Y)][int(figure[i].X)] > 0 {
			return true
		}

	}
	return false
}

func main() {
	pixelgl.Run(run)
}
