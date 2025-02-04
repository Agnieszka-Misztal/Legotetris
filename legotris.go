package main

import (
	"image"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

func run() {

	var grid [20][10]int

	figures := [7][4]int{
		{1, 3, 5, 7}, // I  0 1
		{2, 4, 5, 7}, // Z  2 3
		{3, 5, 4, 6}, // S  4 5
		{1, 3, 2, 5}, // T  6 7
		{2, 3, 5, 7}, // L
		{3, 5, 7, 6}, // J
		{2, 3, 4, 5}, // O
	}

	rand.Seed(time.Now().UnixNano())
	figureType := rand.Intn(7)
	figureColor := rand.Intn(6)
	figureTypeNext := rand.Intn(7)
	figureColorNext := rand.Intn(6)

	var figure [4]pixel.Vec //4 pozycje klocka, vectory
	var figureNext [4]pixel.Vec
	var figureTemp [4]pixel.Vec

	score := 0

	//stworzenie klocka
	for i := 0; i < 4; i++ {
		figure[i].X = float64(figures[figureType][i] % 2)    //ustawienie x na 0 lub 1
		figure[i].Y = float64(figures[figureType][i]/2 + 16) //ustawienie y od 0 do 3
	}
	for i := 0; i < 4; i++ {
		figureNext[i].X = float64(figures[figureTypeNext][i] % 2)    //ustawienie x na 0 lub 1
		figureNext[i].Y = float64(figures[figureTypeNext][i]/2 + 16) //ustawienie y od 0 do 3
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

	//win.SetSmooth(true)

	// załadowanie fonta
	face, err := loadTTF("Legothick.ttf", 80)
	if err != nil {
		panic(err)
	}
	// wygenerowanie z fontu obrazka
	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(312, 680), atlas)

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
	pic2, err := loadPicture("brickBlue.png")
	if err != nil {
		panic(err)
	}

	pic3, err := loadPicture("brickGreen.png")
	if err != nil {
		panic(err)
	}

	pic4, err := loadPicture("brickRed.png")
	if err != nil {
		panic(err)
	}

	pic5, err := loadPicture("brickWhite.png")
	if err != nil {
		panic(err)
	}
	pic6, err := loadPicture("brickYellow.png")
	if err != nil {
		panic(err)
	}

	backPic, err := loadPicture("back3.png")
	if err != nil {
		panic(err)
	}

	block := pixel.NewSprite(pic, pic.Bounds()) //obiekt sluzacy do wyswietlenia obrazka
	block2 := pixel.NewSprite(pic2, pic2.Bounds())
	block3 := pixel.NewSprite(pic3, pic3.Bounds())
	block4 := pixel.NewSprite(pic4, pic4.Bounds())
	block5 := pixel.NewSprite(pic5, pic5.Bounds())
	block6 := pixel.NewSprite(pic6, pic6.Bounds())
	backSprite := pixel.NewSprite(backPic, backPic.Bounds())

	coloredBlocks := [6]*pixel.Sprite{block, block2, block3, block4, block5, block6}

	for !win.Closed() {
		dt := time.Since(lastTime).Seconds() // czas ktory uplynal od poprzedniej klatki
		lastTime = time.Now()

		leftRightTime += dt
		moveDwonTime += dt
		moveLeftOrRight = 0

		win.Clear(colornames.Lightsteelblue) //wyczyszczenie okna przed wyswietleniem

		txt.Clear()
		txt.WriteString("SCORE " + strconv.Itoa(score))

		if win.Pressed(pixelgl.KeyLeft) {
			moveLeftOrRight = -1
			//positionX--

		}
		if win.Pressed(pixelgl.KeyRight) {
			moveLeftOrRight = 1
			//positionX++

		}
		if win.Pressed(pixelgl.KeyDown) {
			moveDwonTime = 1.0
			positionY--

		}
		//obracanie klocka!!!!!
		if win.JustPressed(pixelgl.KeyUp) { // just pressed zadziala tylko raz

			for i := 0; i < 4; i++ {

				figureTemp[i] = figure[i]
			}

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
			//sprawdzanie czy poza plansza
			sideMove := checkCollisionSides(grid, figure)
			if sideMove != 0 {

				if sideMove < 0 {
					for i := 0; i < 4; i++ {
						figure[i].X -= sideMove
					}
				}

				if sideMove > 9 {
					for i := 0; i < 4; i++ {
						figure[i].X -= (sideMove - 9)
					}

				}
			}

			if checkCollision(grid, figure) {
				for i := 0; i < 4; i++ {

					figure[i] = figureTemp[i]
				}

			}

		}

		//opadanie klocka
		if moveDwonTime >= 1.0 {
			moveDwonTime = 0.0
			for i := 0; i < 4; i++ {
				figure[i].Y--

			}

			// sprawdzanie kolizji od dolu ze sciena lub klockiem, jezeli jest cofnij o jedno do gory
			// wpisz na plansze, wypelniajac jedynkami
			if checkCollision(grid, figure) {
				score += 10

				for i := 0; i < 4; i++ {
					x := int(figure[i].X)
					y := int(figure[i].Y + 1)

					//tabica od 0, a zakaldamy ze 0 to puste miejsce na tablicy
					grid[y][x] = figureColor + 1

				}

				figureType = figureTypeNext
				figureColor = figureColorNext

				figureTypeNext = rand.Intn(7)
				figureColorNext = rand.Intn(6)

				//stworzenie klocka od nowa do góry
				for i := 0; i < 4; i++ {
					figure[i].X = float64(figures[figureType][i] % 2)    //ustawienie x na 0 lub 1
					figure[i].Y = float64(figures[figureType][i]/2 + 16) //ustawienie y od 0 do 3
				}

				for i := 0; i < 4; i++ {
					figureNext[i].X = float64(figures[figureTypeNext][i] % 2)    //ustawienie x na 0 lub 1
					figureNext[i].Y = float64(figures[figureTypeNext][i]/2 + 16) //ustawienie y od 0 do 3
				}

				//sprawdzanie kolizji czy nowotowrzony klocek nie nachodzi na inny, czy koniec gry
				if checkCollision(grid, figure) {
					score = 0
					for y := 0; y < 20; y++ {
						for x := 0; x < 10; x++ {
							grid[y][x] = 0

						}

					}
				}
			}
		}

		if leftRightTime >= 0.1 && moveLeftOrRight != 0 {
			leftRightTime = 0.0

			for i := 0; i < 4; i++ {
				figure[i].X += float64(moveLeftOrRight)

			}

			//jezeli wykryto kolizje, cofinj przesuniecie
			if checkCollision(grid, figure) {
				for i := 0; i < 4; i++ {
					figure[i].X -= float64(moveLeftOrRight)

				}
			}

		}

		//sprawdzanie czy jest zapelniona linia
		lineToOverwrite := 0

		// for po wszytkich wierszach
		for y := 0; y < 20; y++ {

			columnCount := 0

			for x := 0; x < 10; x++ {

				if grid[y][x] > 0 {
					columnCount++
				}

				grid[lineToOverwrite][x] = grid[y][x]
			}

			if columnCount < 10 {

				lineToOverwrite++
			} else {
				score += 100
			}
		}

		//rysowanie tła
		backSprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		//rysowanie planszy/grid
		for y := 0; y < 20; y++ {
			for x := 0; x < 10; x++ {
				if grid[y][x] > 0 {

					coloredBlocks[grid[y][x]-1].Draw(win, pixel.IM.Moved(pixel.V(float64(x*32+16+352), float64(y*25+16+134))))
				}
			}
		}

		//sortowanie rysowania klocka od najnizszych elementow (wg Y)

		for i := 0; i < 4; i++ {
			figureTemp[i] = figure[i]
		}
		sort.Slice(figureTemp[:], func(i, j int) bool {
			return figureTemp[i].Y < figureTemp[j].Y
		})
		//rysowanie klocka
		for i := 0; i < 4; i++ {
			coloredBlocks[figureColor].Draw(win, pixel.IM.Moved(pixel.V(figureTemp[i].X*32.0+16.0+352, figureTemp[i].Y*25+16.0+134)))
		}
		//rysowanie klocka nastepnego w poczekalni
		for i := 0; i < 4; i++ {
			coloredBlocks[figureColorNext].Draw(win, pixel.IM.Moved(pixel.V(figureNext[i].X*32.0+16.0+64, figureNext[i].Y*25+16.0+134)))
		}
		//rysowanie wyniku
		txt.Draw(win, pixel.IM)

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

func checkCollisionSides(grid [20][10]int, figure [4]pixel.Vec) float64 {
	sideX := 0.0
	for i := 0; i < 4; i++ {

		//sprawdzanie czy klocek nie wychodzi poza boczne sciany
		if figure[i].X < 0 && figure[i].X < sideX {
			sideX = figure[i].X

		}
		if figure[i].X > 9 && figure[i].X > sideX {
			sideX = figure[i].X
		}
	}
	return sideX
}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func main() {
	pixelgl.Run(run)
}
