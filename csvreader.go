//	.csv file reader with calculating the inputted cell expressions.
//	A test task for YADRO received on 26.01 with the DL of 02.02.

package main

import (
	"fmt"
	"os"
	"strconv"
)

type Matrix struct {
	Columns []string
	Rows    []uint64
	Links   map[string]string
}

func main() {

	// check whether the arg of a filename is the only argument provided
	if len(os.Args) != 2 {
		fmt.Println("The amount of inputted args is invalid.")
		os.Exit(1)
	}

	arg := os.Args[1]
	Input, err := os.ReadFile(arg)

	errorCheck(err)

	// check whether the first character of the .csv is a comma
	if rune(Input[0]) != ',' {
		fmt.Println("The table formatting is invalid since the first cell in the .csv file has to be empty.")
		os.Exit(1)
	}

	Cells := parseInput(Input)
	Cells = expressionHandler(Cells)
	printCSV(Cells)
}

func isCellEmpty(temp string, index int) {
	if temp == "" {
		fmt.Printf("\nOne of the values appears missing on the line #%v.\n", index+1)
		os.Exit(1)
	}
}

func errorCheck(err error) {
	if err != nil {
		fmt.Println("The file is invalid.")
		os.Exit(1)
	}
}

func safeStrToIntConvert(number *int64, value string) {
	var err error
	*number, err = strconv.ParseInt(value, 10, 64)
	if err != nil {
		fmt.Printf("\nAn operand in an expression cell [%v] refers to an invalid value.\n", value)
		os.Exit(1)
	}
}

func parseInput(content []byte) Matrix {
	// parsing the content file

	temp := ""          // a variable for a value in the iterated cell
	columnIndex := 0    // the index of a column iterated
	rowIndexed := false // the flag for an iterated row
	rowIndex := 0       // the index of a row iterated, also the amount of lines in a file after the loop

	tempMatrix := Matrix{
		[]string{},
		[]uint64{},
		make(map[string]string),
	}

	for i := 0; i < len(content); i += 1 {
		if rowIndex == 0 {
			// parsing the first row containing cell columns' names
			if i > 0 {
				switch {
				case content[i] == ',':
					isCellEmpty(temp, rowIndex)
					tempMatrix.Columns = append(tempMatrix.Columns, temp)
					temp = ""

				case content[i] == '\r' || content[i] == '\n':
					isCellEmpty(temp, rowIndex)
					tempMatrix.Columns = append(tempMatrix.Columns, temp)
					temp = ""
					rowIndex += 1
					// extra iteration in order to iterate over "\n" prefix
					if content[i] == '\r' {
						i += 1
					}
					break

				default:
					temp += string(content[i])
				}
			}
		} else {
			// parsing the other rows containing cell columns' indexes
			switch {
			case content[i] == ',':
				isCellEmpty(temp, rowIndex)
				if rowIndexed == false {
					num, err := strconv.ParseUint(temp, 10, 64)
					if err == nil && num != 0 {
						tempMatrix.Rows = append(tempMatrix.Rows, num)
					} else {
						fmt.Printf("\nThe index of a row #%v is not a positive integer value.\n", rowIndex)
						os.Exit(1)
					}
					rowIndexed = true
				} else {
					if len(tempMatrix.Columns) == columnIndex {
						fmt.Printf("There are more cells in a row #%v than columns given.\n", rowIndex)
						os.Exit(1)
					}
					checkIfInt(temp)
					tempMatrix.Links[tempMatrix.Columns[columnIndex-1]+strconv.FormatInt(int64(tempMatrix.Rows[rowIndex-1]), 10)] = temp
				}

				columnIndex += 1
				temp = ""

			// if the file is written on Unix, '\r' prefix may appear missing
			case content[i] == '\r' || content[i] == '\n':
				isCellEmpty(temp, rowIndex)
				checkIfInt(temp)
				tempMatrix.Links[tempMatrix.Columns[columnIndex-1]+strconv.FormatInt(int64(tempMatrix.Rows[rowIndex-1]), 10)] = temp
				temp = ""
				// extra iteration in order to iterate over "\n" prefix
				if content[i] == '\r' {
					i += 1
				}
				rowIndex += 1
				rowIndexed = false
				columnIndex = 0
				break

			default:
				temp += string(content[i])
			}
		}
	}

	// the last cell parsing
	if temp != "" {
		checkIfInt(temp)
		tempMatrix.Links[tempMatrix.Columns[columnIndex-1]+strconv.FormatInt(int64(tempMatrix.Rows[rowIndex-1]), 10)] = temp
	}

	// check whether the formatting of the number of cells is correct
	if len(tempMatrix.Links) != (len(tempMatrix.Columns) * len(tempMatrix.Rows)) {
		fmt.Println("The table is corrupted. Some cells are not presented.")
		os.Exit(1)
	}

	return tempMatrix
}

func expressionHandler(matrix Matrix) Matrix {
	// extract the cells starting with '=' and replace them with the value they refer to

	for i := 0; i < len(matrix.Rows); i += 1 {
		for j := 0; j < len(matrix.Columns); j += 1 {
			// the key containing a link to the iterated cell (ARG)
			cellLink := matrix.Columns[j] + strconv.FormatInt(int64(matrix.Rows[i]), 10)
			matrix.Links[cellLink] = expressionFixer(cellLink, matrix, cellLink)
		}
	}
	return matrix
}

func expressionFixer(cellLink string, matrix Matrix, cellLinkForRecursionCheck string) string {
	//  parsing of a cell with an expression inside
	/*  cellLinkForRecursionCheck is needed to store the genesis value
	and be used for checking whether the recursion in cell expressions is present. */
	if (matrix.Links[cellLink])[0] == '=' {
		opFound := false
		var opIndex int
		var op byte

		for k := 1; k < len(matrix.Links[cellLink]); k += 1 {
			if ((matrix.Links[cellLink])[k] == '+' ||
				(matrix.Links[cellLink])[k] == '-' ||
				(matrix.Links[cellLink])[k] == '*' ||
				(matrix.Links[cellLink])[k] == '/') &&
				opFound == false {
				op = matrix.Links[cellLink][k]
				opIndex = k
				opFound = true
				break
			}
		}

		if opFound == true {
			// a cell link the first operand refers to; in "=B1+Cell30" fstNumArg would be "B1"
			fstNumArg := (matrix.Links[cellLink])[1:(opIndex)]
			// a cell link the second operand refers to; in "=B1+Cell30" sndNumArg would be "Cell30"
			sndNumArg := (matrix.Links[cellLink])[(opIndex + 1):]

			// checking whether any operand in the expression cell refers to itself
			if fstNumArg == cellLink || sndNumArg == cellLink {
				fmt.Printf("\nThere is an expression cell '%v' producing recursion.\n", cellLink)
				os.Exit(1)
			}

			// initializing the variables responsible for storing the integer values of map's elements
			var fstNum, sndNum int64

			// check whether the operands are referring to a non-existing cell
			if matrix.Links[fstNumArg] == "" {
				fmt.Printf("\nThere is no cell with the link '%v' given (expression operand is invalid).\n", fstNumArg)
				os.Exit(1)
			}
			if matrix.Links[sndNumArg] == "" {
				fmt.Printf("\nThere is no cell with the link '%v' given (expression operand is invalid).\n", sndNumArg)
				os.Exit(1)
			}

			// checking whether the reference of an operand is to an expression as well, recursion of expressionFixer()
			recursionFixer(fstNumArg, matrix, cellLinkForRecursionCheck)
			recursionFixer(sndNumArg, matrix, cellLinkForRecursionCheck)

			// checking whether any operand in the expression cell refers to an invalid value
			safeStrToIntConvert(&fstNum, (matrix.Links[fstNumArg]))
			safeStrToIntConvert(&sndNum, (matrix.Links[sndNumArg]))

			switch op {
			case '+':
				matrix.Links[cellLink] = strconv.FormatInt(int64(fstNum+sndNum), 10)
			case '-':
				matrix.Links[cellLink] = strconv.FormatInt(int64(fstNum-sndNum), 10)
			case '*':
				matrix.Links[cellLink] = strconv.FormatInt(int64(fstNum*sndNum), 10)
			case '/':
				if sndNum == 0 {
					fmt.Printf("\nThere is an expression cell %v containing division by zero.\n", cellLink)
					os.Exit(1)
				} else {
					matrix.Links[cellLink] = strconv.FormatInt(int64(fstNum/sndNum), 10)
				}
			}
		} else {
			fmt.Printf("\nInvalid cell expression. In the cell %v no operator was found.\n", cellLink)
			os.Exit(1)
		}
	}
	return matrix.Links[cellLink]
}

func printCSV(matrix Matrix) {
	fmt.Println()
	// column names
	for i := 0; i < len(matrix.Columns); i += 1 {
		fmt.Print(",", matrix.Columns[i])
	}
	// row indexes with the cell values
	for i := 0; i < len(matrix.Rows); i += 1 {
		fmt.Print("\n", matrix.Rows[i])
		for j := 0; j < len(matrix.Columns); j += 1 {
			fmt.Print(",", matrix.Links[matrix.Columns[j]+strconv.FormatInt(int64(matrix.Rows[i]), 10)])
		}
	}
	fmt.Println()
}

func checkIfInt(str string) {
	if str[0] != '=' {
		temp, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			fmt.Printf("\nOne of the inputted cells contains an incorrect value [%s].\n", str)
			os.Exit(1)
			/* 	The next line is just a never executing dummy action (a plug)
			in order to move around the 'declared but not used' of temp */
			temp++
		}
	}
}

func recursionFixer(arg string, matrix Matrix, cellLinkForRecursionCheck string) {
	if (matrix.Links[arg])[0] == '=' {
		if matrix.Links[arg] == matrix.Links[cellLinkForRecursionCheck] {
			fmt.Printf("\nThe structure of cell expressions produces a sempiternal recursion.\n")
			os.Exit(1)
		}
		expressionFixer(arg, matrix, cellLinkForRecursionCheck)
	}
}
