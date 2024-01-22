package Netpbm

//library import
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PPM, Pixel, Point structures
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Pixel struct {
	R, G, B uint8
}

type Point struct {
	X, Y int
}

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	// Opening/closing of file and cheking for errors
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return nil, err
	}
	defer file.Close()

	// Definition of variables and scan
	ReadPPM := PPM{}
	scanner := bufio.NewScanner(file)
	var line, DataColumn, RGBPlacement int = 0, 0, 0

	// Scan of the file and information collection
	for scanner.Scan() {
		// Check for '#' to skip comments
		if strings.HasPrefix(scanner.Text(), "#") {

		} else {
			//magicNumber
			if ReadPPM.magicNumber != "P3" && ReadPPM.magicNumber != "P6" {
				ReadPPM.magicNumber = scanner.Text()
			} else if ReadPPM.width == 0 && ReadPPM.height == 0 {
				// splits the line
				sizes := strings.Split(scanner.Text(), " ")
				//fill struct
				ReadPPM.width, _ = strconv.Atoi(sizes[0])
				ReadPPM.height, _ = strconv.Atoi(sizes[1])

				// Creation of the array
				ReadPPM.data = make([][]Pixel, ReadPPM.height)
				for i := range ReadPPM.data {
					ReadPPM.data[i] = make([]Pixel, ReadPPM.width)
				}

			} else if ReadPPM.max == 0 {
				// convert string -> int -> uint8
				intMax, _ := strconv.Atoi(scanner.Text())
				ReadPPM.max = uint8(intMax)

			} else {
				// Check the magicNumber
				if ReadPPM.magicNumber == "P3" {
					// split the line
					TheLine := strings.Split(scanner.Text(), " ")

					Pixel := Pixel{}

					//fill the array
					for i := 0; i < len(TheLine); i++ {
						nombre, _ := strconv.Atoi(TheLine[i])

						switch RGBPlacement {
						case 0:
							Pixel.R = uint8(nombre)
							RGBPlacement++
						case 1:
							Pixel.G = uint8(nombre)
							RGBPlacement++
						case 2:
							Pixel.B = uint8(nombre)
							RGBPlacement = 0

							ReadPPM.data[line][DataColumn] = Pixel
							DataColumn++
						}
					}

				} else if ReadPPM.magicNumber == "P6" {
					// ...
				}
				// count lines
				line++
				// Reset DataColumn and RGBPlacement for next line
				DataColumn = 0
				RGBPlacement = 0
			}
		}
	}
	// Final return
	return &PPM{ReadPPM.data, ReadPPM.width, ReadPPM.height, ReadPPM.magicNumber, ReadPPM.max}, nil
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[x][y] = value
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error {
	// Opening/closing of file and cheking for errors
	file, err := os.Create(filename)
	defer file.Close()

	//write magicNumber, width, height and max
	fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	//write data
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			ThePixel := ppm.data[i][j]
			fmt.Fprint(file, ThePixel.R, " ", ThePixel.G, " ", ThePixel.B, " ")
		}
		//next line
		fmt.Fprintln(file)
	}

	return err
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	// range the array
	for DataLine := 0; DataLine < ppm.height; DataLine++ {
		for DataColumn := 0; DataColumn < ppm.width; DataColumn++ {
			// Invertdata
			ppm.data[DataLine][DataColumn].R = uint8(ppm.max) - ppm.data[DataLine][DataColumn].R
			ppm.data[DataLine][DataColumn].G = uint8(ppm.max) - ppm.data[DataLine][DataColumn].G
			ppm.data[DataLine][DataColumn].B = uint8(ppm.max) - ppm.data[DataLine][DataColumn].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	// Count changes
	var Iteration int = (ppm.width / 2)
	//exchange data
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < Iteration; j++ {
			ppm.data[i][j], ppm.data[i][ppm.width-j-1] = ppm.data[i][ppm.width-j-1], ppm.data[i][j]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	// Count changes
	var Iteration int = (ppm.height / 2)
	//exchange data
	for i := 0; i < ppm.width; i++ {
		for j := 0; j < Iteration; j++ {
			ppm.data[j][i], ppm.data[ppm.height-j-1][i] = ppm.data[ppm.height-j-1][i], ppm.data[j][i]
		}
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	// range the array
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			//calculate new value
			ppm.data[i][j].R = uint8(float64(ppm.data[i][j].R) * float64(maxValue) / float64(ppm.max))
			ppm.data[i][j].G = uint8(float64(ppm.data[i][j].G) * float64(maxValue) / float64(ppm.max))
			ppm.data[i][j].B = uint8(float64(ppm.data[i][j].B) * float64(maxValue) / float64(ppm.max))
		}
	}
	ppm.max = maxValue
}

// Rotate90CW rotates the PPM image 90° clockwise.
func (ppm *PPM) Rotate90CW() {
	rotatedData := make([][]Pixel, ppm.width)
	for i := range rotatedData {
		rotatedData[i] = make([]Pixel, ppm.height)
	}

	// rotate
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			rotatedData[i][j] = ppm.data[(ppm.width-1)-j][i]
		}
	}

	// new values
	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = rotatedData
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Define uint8 array
	data := make([][]uint8, ppm.width)
	for i := range data {
		data[i] = make([]uint8, ppm.height)
	}

	// convert data
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			R := int(ppm.data[i][j].R)
			G := int(ppm.data[i][j].G)
			B := int(ppm.data[i][j].B)
			average := (R + G + B) / 3
			data[i][j] = uint8(average)
		}
	}

	//return converted pgm struct
	return &PGM{data, ppm.width, ppm.height, "P2", ppm.max}
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	// Define bool array
	data := make([][]bool, ppm.width)
	for i := range data {
		data[i] = make([]bool, ppm.height)
	}

	// convert data
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			if uint8((int(ppm.data[i][j].R)+int(ppm.data[i][j].G)+int(ppm.data[i][j].B))/3) < ppm.max/2 {
				data[i][j] = true
			} else {
				data[i][j] = false
			}
		}
	}

	//return converted pbm struct
	return &PBM{data, ppm.width, ppm.height, "P1"}
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {

	// Define the number of squares up/down and right/left
	dx := Abs(p2.X - p1.X)
	dy := Abs(p2.Y - p1.Y)

	//Bresenham algorithm
	sx := -1
	if p1.X < p2.X {
		sx = 1
	}

	sy := -1
	if p1.Y < p2.Y {
		sy = 1
	}

	err := dx - dy

	for {
		// Scope safe check
		if (p1.Y < ppm.height) && (p1.X < ppm.width) {
			ppm.data[p1.Y][p1.X] = color
		}

		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p1.X += sx
		}
		if e2 < dx {
			err += dx
			p1.Y += sy
		}
	}
}

// Func to get absolute value
func Abs(nb int) int {
	if nb < 0 {
		nb = -nb
		return nb
	} else {
		return nb
	}
}

// DrawRectangle draws a rectangle.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	// Define all 4 corners of the rectangle
	TopLeft := Point{
		X: p1.X,
		Y: p1.Y,
	}

	TopRight := Point{
		X: (p1.X + width),
		Y: p1.Y,
	}

	BottomLeft := Point{
		X: p1.X,
		Y: (p1.Y + height),
	}

	BottomRight := Point{
		X: (p1.X + width),
		Y: (p1.Y + height),
	}

	ppm.DrawLine(TopLeft, TopRight, color)
	ppm.DrawLine(BottomLeft, BottomRight, color)
	ppm.DrawLine(TopLeft, BottomLeft, color)
	ppm.DrawLine(TopRight, BottomRight, color)
}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	ppm.DrawRectangle(p1, width, height, color)
	if width == 0 && height == 0 {
		return
	} else if width == 0 && height > 0 {
		ppm.DrawFilledRectangle(p1, width, height-1, color)
	} else if width > 0 && height == 0 {
		ppm.DrawFilledRectangle(p1, width-1, height, color)
	} else if width > 0 && height > 0 {
		ppm.DrawFilledRectangle(p1, width-1, height-1, color)
	}
}

// DrawTriangle draws a triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p1, p3, color)
	ppm.DrawLine(p2, p3, color)
}
