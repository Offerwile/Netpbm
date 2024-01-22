package Netpbm

//library import
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PBM structure
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	// Opening/closing of file and cheking for errors
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return nil, err
	}
	defer file.Close()

	// Definition of variables and scan
	ReadPBM := PBM{}
	scanner := bufio.NewScanner(file)
	var line int = 0

	// Scan of the file and information collection
	for scanner.Scan() {
		// Check for '#' to skip comments
		if strings.HasPrefix(scanner.Text(), "#") {

		} else {
			//magicNumber
			if ReadPBM.magicNumber != "P1" && ReadPBM.magicNumber != "P4" {
				ReadPBM.magicNumber = scanner.Text()

			} else if ReadPBM.width == 0 && ReadPBM.height == 0 {
				// splits the line
				sizes := strings.Split(scanner.Text(), " ")
				//fill struct
				ReadPBM.width, _ = strconv.Atoi(sizes[0])
				ReadPBM.height, _ = strconv.Atoi(sizes[1])

				// Creation of the array
				ReadPBM.data = make([][]bool, ReadPBM.height)
				for i := range ReadPBM.data {
					ReadPBM.data[i] = make([]bool, ReadPBM.width)
				}

			} else {
				// Check the magicNumber
				if ReadPBM.magicNumber == "P1" {
					// split the line
					TheLine := strings.Split(scanner.Text(), " ")

					//fill the array
					for i := 0; i < len(TheLine); i++ {
						if TheLine[i] == "1" {
							ReadPBM.data[line][i] = true
						} else {
							ReadPBM.data[line][i] = false
						}
					}
					// Check the magicNumber
				} else if ReadPBM.magicNumber == "P4" {
					// Number of bytes per line
					var bytes int
					if ReadPBM.width%8 == 0 {
						bytes = (ReadPBM.width / 8)
					} else {
						bytes = (ReadPBM.width / 8) + 1
					}
					//padding
					padding := (bytes * 8) - ReadPBM.width

					// convert line into ZeroAndOne
					var ZeroAndOne string = ""

					// range line
					for i := 0; i < len(scanner.Text()); i++ {
						bin := fmt.Sprintf("%08b", scanner.Text()[i])
						// Remove padding from last bytes of each line
						if bytes != 1 {
							if i != 0 && (i+1)%bytes == 0 {
								bin = bin[:len(bin)-padding]
							}
						} else if bytes == 1 {
							bin = bin[:len(bin)-padding]
						}
						// Add the char converted in ZeroAndOne to final string
						ZeroAndOne += bin
					}

					// Fill the array
					var DataColumn, DataLine int = 0, 0
					for _, Value := range ZeroAndOne {
						if DataColumn == ReadPBM.width {
							DataColumn = 0
							DataLine += 1
						}
						if DataLine != ReadPBM.width {
							if Value == '1' {
								ReadPBM.data[DataLine][DataColumn] = true
							} else {
								ReadPBM.data[DataLine][DataColumn] = false
							}
						}
					}
				}
				// count lines
				line++
			}
		}
	}
	// Final return
	return &PBM{ReadPBM.data, ReadPBM.width, ReadPBM.height, ReadPBM.magicNumber}, nil
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[x][y]
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[x][y] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	// Opening/closing of file and cheking for errors
	file, err := os.Create(filename)
	defer file.Close()

	//write magicNumber, width and height
	fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	//write 0 and 1
	for _, a := range pbm.data {
		for _, b := range a {
			if b == true {
				fmt.Fprint(file, "1 ")
			} else {
				fmt.Fprint(file, "0 ")
			}
		}
		//next line
		fmt.Fprintln(file)
	}
	return err
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	// range the array
	for DataLine := 0; DataLine < pbm.height; DataLine++ {
		for DataColumn := 0; DataColumn < pbm.width; DataColumn++ {
			// Invertdata
			if pbm.data[DataLine][DataColumn] {
				pbm.data[DataLine][DataColumn] = false
			} else {
				pbm.data[DataLine][DataColumn] = true
			}
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	// Count changes
	var Iteration int = (pbm.width / 2)
	//exchange data
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < Iteration; j++ {
			pbm.data[i][j], pbm.data[i][pbm.width-j-1] = pbm.data[i][pbm.width-j-1], pbm.data[i][j]
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	// Count changes
	var Iteration int = (pbm.height / 2)
	//exchange data
	for i := 0; i < pbm.width; i++ {
		for j := 0; j < Iteration; j++ {
			pbm.data[j][i], pbm.data[pbm.height-j-1][i] = pbm.data[pbm.height-j-1][i], pbm.data[j][i]
		}
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
