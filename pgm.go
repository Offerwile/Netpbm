package Netpbm

//library import
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PGM structure
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	// Opening/closing of file and cheking for errors
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return nil, err
	}
	defer file.Close()

	// Definition of variables and scan
	ReadPGM := PGM{}
	scanner := bufio.NewScanner(file)
	var line int = 0

	// Scan of the file and information collection
	for scanner.Scan() {
		// Check for '#' to skip comments
		if strings.HasPrefix(scanner.Text(), "#") {

		} else {
			//magicNumber
			if ReadPGM.magicNumber != "P2" && ReadPGM.magicNumber != "P5" {
				ReadPGM.magicNumber = scanner.Text()
			} else if ReadPGM.width == 0 && ReadPGM.height == 0 {
				// splits the line
				sizes := strings.Split(scanner.Text(), " ")
				//fill struct
				ReadPGM.width, _ = strconv.Atoi(sizes[0])
				ReadPGM.height, _ = strconv.Atoi(sizes[1])

				// Creation of the array
				ReadPGM.data = make([][]uint8, ReadPGM.height)
				for i := range ReadPGM.data {
					ReadPGM.data[i] = make([]uint8, ReadPGM.width)
				}

			} else if ReadPGM.max == 0 {
				// convert string -> int -> uint8
				intMax, _ := strconv.Atoi(scanner.Text())
				ReadPGM.max = uint8(intMax)

			} else {
				// Check the magicNumber
				if ReadPGM.magicNumber == "P2" {
					// split the line
					TheLine := strings.Split(scanner.Text(), " ")

					//fill the array
					for i := 0; i < ReadPGM.width; i++ {
						IntValue, _ := strconv.Atoi(TheLine[i])
						ReadPGM.data[line][i] = uint8(IntValue)
					}

				} else if ReadPGM.magicNumber == "P5" {
					// Filling of the array
					var DataLine, DataColumn int = 0, 0

					for _, data := range scanner.Text() {
						// Reset for new line
						if DataLine == ReadPGM.width {
							DataLine = 0
							DataColumn += 1
						}
						ReadPGM.data[DataColumn][DataLine] = uint8(data)
						DataLine++
					}
				}
				// count lines
				line++
			}
		}
	}
	// Final return
	return &PGM{ReadPGM.data, ReadPGM.width, ReadPGM.height, ReadPGM.magicNumber, ReadPGM.max}, nil
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[x][y]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	// Opening/closing of file and cheking for errors
	file, err := os.Create(filename)
	defer file.Close()

	//write magicNumber, width, height and max
	fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	//write data
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			fmt.Fprint(file, pgm.data[i][j], " ")
		}
		//next line
		fmt.Fprintln(file)
	}

	return err
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	// range the array
	for DataLine := 0; DataLine < pgm.height; DataLine++ {
		for DataColumn := 0; DataColumn < pgm.width; DataColumn++ {
			// Invertdata
			pgm.data[DataLine][DataColumn] = uint8(pgm.max) - pgm.data[DataLine][DataColumn]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	// Count changes
	var Iteration int = (pgm.width / 2)
	//exchange data
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < Iteration; j++ {
			pgm.data[i][j], pgm.data[i][pgm.width-j-1] = pgm.data[i][pgm.width-j-1], pgm.data[i][j]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	// Count changes
	var Iteration int = (pgm.height / 2)
	//exchange data
	for i := 0; i < pgm.width; i++ {
		for j := 0; j < Iteration; j++ {
			pgm.data[j][i], pgm.data[pgm.height-j-1][i] = pgm.data[pgm.height-j-1][i], pgm.data[j][i]
		}
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	// range the array
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			//calculate new value
			pgm.data[i][j] = uint8(float64(pgm.data[i][j]) * float64(maxValue) / float64(pgm.max))
		}
	}
	pgm.max = maxValue
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	rotatedData := make([][]uint8, pgm.width)
	for i := range rotatedData {
		rotatedData[i] = make([]uint8, pgm.height)
	}

	// rotate
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			rotatedData[i][j] = pgm.data[(pgm.width-1)-j][i]
		}
	}

	// new values
	pgm.width, pgm.height = pgm.height, pgm.width
	pgm.data = rotatedData
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	// Define bool array
	data := make([][]bool, pgm.width)
	for i := range data {
		data[i] = make([]bool, pgm.height)
	}

	// convert data
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			if pgm.data[i][j] >= (uint8(pgm.max) / 2) {
				data[i][j] = false
			} else if pgm.data[i][j] < (uint8(pgm.max) / 2) {
				data[i][j] = true
			}
		}
	}

	//return converted pbm struct
	return &PBM{data, pgm.width, pgm.height, "P1"}
}
