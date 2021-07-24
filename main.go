package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	// map of coordinates
	// clearPaths contains every path that user can move to
	clearPaths = make([][]int, 0)
	blockPaths = make([][]int, 0)
	// store all possibility coordinates of treasure
	possibleTreasurePosition = make([][]int, 0)

	// coordinate for current position
	currentPosition = make([]int, 0)

	// coordinate for treasure position
	treasurePosition = make([]int, 0)

	// size of board column
	totalColumn int

	// rowData contains value per row, using row index as key, and slice of string as value
	// to modify a data in a row, you just need to find its coordinate,
	// imagine row as x and column as y.
	// coordinate table will look like this:
	//-----------------------------------------
	//| (x,y) |   0   |   1   |   2   |   3   |  -˃ column (y)
	//-----------------------------------------
	//|   0   |  0,0  |  0,1  |  0,2  |  0,3  |
	//|   1   |  1,0  |  1,1  |  1,2  |  1,3  |
	//|   2   |  2,0  |  2,1  |  2,2  |  2,3  |
	//|   3   |  3,0  |  3,1  |  3,2  |  3,3  |
	//-----------------------------------------
	//   |
	//   ˅
	// row (x)
	// example: if you want to change the value in row 2 and column 3 to 'X',
	// just use this way: rowData[2][3] = 'X'
	rowData = make(map[int][]string, 0)

	doneGoUp    bool
	doneGoRight bool
	doneGoDown  bool
)

func main() {
	reset()
	// load grid data from file
	gridData, err := ioutil.ReadFile("board_grid.txt")
	if err != nil {
		log.Fatal("Failed to read board config data", err.Error())
	}

	// load grid data to variables
	err = setupNewBoard(gridData)
	if err != nil {
		log.Fatal(err.Error())
	}

	// generate treasure location
	generateTreasurePosition()

	printBoardData("")
	play()
}

func setupNewBoard(gridData []byte) (err error) {
	// split each rows by new line (\n)
	rows := strings.Split(string(gridData), "\n")

	// loop rows to get data per row
	for rowIndex, rowValue := range rows {
		values := strings.Split(rowValue, "")
		totalColumn = len(values)
		if totalColumn < 1 {
			return errors.New("Empty board data")
		}

		dataPerRow := make([]string, 0)
		// loop data per row
		for columnIndex, data := range values {
			dataCoordinate := []int{rowIndex, columnIndex}
			switch data {
			case "#":
				// ignore obstacles
				blockPaths = append(blockPaths, dataCoordinate)
			case ".":
				clearPaths = append(clearPaths, dataCoordinate)
			case "X":
				// starting position is clear path too.
				clearPaths = append(clearPaths, dataCoordinate)
				currentPosition = dataCoordinate
			default:
				return errors.New("Invalid character on the grid detected")
			}
			dataPerRow = append(dataPerRow, data)
		}
		// store data each row
		rowData[rowIndex] = dataPerRow
	}

	return nil
}

func generateTreasurePosition() {
	for i := 0; i < len(clearPaths)/3; i++ {
		// generate possible treasure location
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(clearPaths))))
		idx := int(randomIndex.Int64())

		// exclude starting position from treasure possibility location, exclude 4,4 and 4,6 too since it's impossible to get there.
		if !reflect.DeepEqual(currentPosition, clearPaths[idx]) && !reflect.DeepEqual(clearPaths[idx],[]int{4,4}) && !reflect.DeepEqual(clearPaths[idx], []int{4,6}) {
			// if random generated same value each loop, it will append same coordinate.
			alreadyAdded, _ := isInPossibleTreasureLocation(clearPaths[idx])
			if !alreadyAdded {
				// add possibility treasure position only
				possibleTreasurePosition = append(possibleTreasurePosition, clearPaths[idx])
				// change clear path value from '.' to '$'
				rowData[clearPaths[idx][0]][clearPaths[idx][1]] = "$"
			}
		}
	}

	// set treasure from possible treasure location

	randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(possibleTreasurePosition))))
	idx := int(randomIndex.Int64())
	treasurePosition = possibleTreasurePosition[idx]
}

// check position is in the clear path
func isInClearPath(x, y int) bool {
	checkCoordinate := []int{x, y}
	for _, clearPath := range clearPaths {
		equal := reflect.DeepEqual(checkCoordinate, clearPath)
		if equal {
			return true
		}
	}
	return false
}
// check position is in the clear path
func isInBlockPath(x, y int) bool {
	checkCoordinate := []int{x, y}
	for _, blockPath := range blockPaths {
		equal := reflect.DeepEqual(checkCoordinate, blockPath)
		if equal {
			return true
		}
	}
	return false
}
// misc. to clear terminal
func clear() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// print information about current condition in board
func printBoardData(message string) {
	clear()
	fmt.Printf("  TREASURE HUNT BOARD\n\n")
	fmt.Printf("    ")

	// printing every column
	for i := 0; i < totalColumn; i++ {
		fmt.Printf("%d", i)
		if i != totalColumn-1 {
			fmt.Print(" ")
		}
	}
	fmt.Printf("\n")

	// printing every rows with its data
	for i := 0; i < len(rowData); i++ {
		fmt.Printf("%d: %v\n", i, rowData[i])
	}

	fmt.Printf("\ncurrent position (X): %+v\n", currentPosition)
	fmt.Printf("\npossible treasure location left:\n %d\n\n", possibleTreasurePosition)

	//fmt.Printf("\n treasure location:\n %d\n\n", treasurePosition)
	// print message if any
	if message != "" {
		fmt.Printf("%s\n", message)
	}
}

func play() {
	// waiting for user input..
	for {
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("input > ")
		input, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		} else {
			printBoardData("")
			processInput(string(input))
		}
	}
}

func processInput(input string) {
	input = strings.TrimSuffix(input, "\n")
	if input == "q" {
		fmt.Println("quitting Treasure Hunt..")
		os.Exit(1)
	}
	if input == "help" {
		printHelp()
		return
	}
	input = strings.ToLower(input)
	split := strings.Split(input, " ")
	if len(split) != 2 {
		fmt.Println("invalid input, required: direction and step(s), see 'help'")
		return
	}
	direction := split[0]

	// generalize direction into 4: up, right, down, left
	switch direction {
	case "up", "north":
		if doneGoUp{
			fmt.Println("cannot go up anymore")
			return
		}
		direction = "up"
	case "right", "east":
		if doneGoRight{
			fmt.Println("cannot go right anymore")
			return
		}
		direction = "right"
	case "down", "south":
		direction = "down"
	default:
		fmt.Println("invalid direction input!")
		return
	}

	step, err := strconv.Atoi(split[1])
	if err != nil {
		fmt.Println("invalid step input! needed: integer", err.Error())
		return
	}

	// navigate this input
	msg := navigate(direction, step)
	printBoardData(msg)
	// checking possibility user found the treasure
	check, index := isInPossibleTreasureLocation(currentPosition)
	if check {
		// remove this position from possible treasure location
		removeVisitedLocationFromPossibleTreasureList(index)
		fmt.Print("checking treasure..")
		time.Sleep(1 * time.Second)
		if !reflect.DeepEqual(treasurePosition, currentPosition) {
			printBoardData("no treasure here..")
		} else {
			// treasure found!
			endGame()
			return
		}
	}

	if doneGoUp && !doneGoRight && !doneGoDown {
		// check right, if next is obstacle, then game over
		if isInBlockPath(currentPosition[0], currentPosition[1]+1){
			gameOver()
		}
	}
	if doneGoUp && doneGoRight && !doneGoDown{
		// check down, if next is obstacle, then game over
		if isInBlockPath(currentPosition[0]+1, currentPosition[1]){
			gameOver()
		}
	}

	if doneGoUp && doneGoRight && doneGoDown{
		gameOver()
	}
}

func removeVisitedLocationFromPossibleTreasureList(index int) {
	copy(possibleTreasurePosition[index:], possibleTreasurePosition[index+1:])
	possibleTreasurePosition[len(possibleTreasurePosition)-1] = nil
	possibleTreasurePosition = possibleTreasurePosition[:len(possibleTreasurePosition)-1]
}

// navigate will move current position, only if new position is valid
func navigate(direction string, step int) (message string) {
	x := currentPosition[0]
	y := currentPosition[1]

	// change coordinate regarding direction and step,
	for i := 1; i <= step; i++ {
		switch direction {
		case "up":
			x = x - 1
			doneGoUp = true
		case "right":
			y = y + 1
			doneGoRight = true
		case "down":
			x = x + 1
			doneGoDown = true
		}

		// check per step validity.
		if !isInClearPath(x, y) {
			return "ups, you cannot move there"
		}
	}

	// replace current position value to '.' in board
	rowData[currentPosition[0]][currentPosition[1]] = "."

	// update current position with new coordinate
	currentPosition = []int{x, y}
	// set X as new coordinate in board
	rowData[x][y] = "X"

	return ""
}
func isInPossibleTreasureLocation(curPos []int) (bool, int) {
	var ret bool
	for i, v := range possibleTreasurePosition {
		eq := reflect.DeepEqual(curPos, v)
		if eq {

			return true, i
		}
	}
	return ret, 0
}

// finish the game with option to quit or play again
func endGame() {
	printBoardData("Congratulations! You found the treasure!")
	fmt.Printf("found treasure location at:\n %d\n\n", treasurePosition)
	for {
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("press Enter to play again, or 'q' to quit > ")
		input, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		} else {
			if strings.TrimSuffix(strings.ToLower(string(input)), "\n") == "q" {
				os.Exit(1)
			} else {
				main()
			}
		}
	}
}

// finish the game with option to quit or play again
func gameOver() {
	printBoardData("Game Over, no more steps available for you!")
	fmt.Printf("treasure location:\n %d\n\n", treasurePosition)
	for {
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("press Enter to play again, or 'q' to quit > ")
		input, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		} else {
			if strings.TrimSuffix(strings.ToLower(string(input)), "\n") == "q" {
				os.Exit(1)
			} else {
				main()
			}
		}
	}
}

// reset variable values
func reset() {
	// map of coordinates
	clearPaths = make([][]int, 0)
	possibleTreasurePosition = make([][]int, 0)
	currentPosition = make([]int, 0)
	treasurePosition = make([]int, 0)
	rowData = make(map[int][]string, 0)
	doneGoDown,doneGoRight, doneGoUp = false, false,false
}

func printHelp() {
	fmt.Println(`
============================= H E L P ==================================
Input format: directions [space] step 

Available directions:
	input (string): up, north, down, south, right, east

Step:
	input (integer) : 0 ... positive value 

Example: 
	I want to go upper for 2 steps,
	just type: 'up 2' or 'north 2'

Rule:
	You only have 3 steps available to turn.

	First you can only move 'up', 
	After that you can only move 'right', 
	Last step you can only move 'down'
========================================================================`)
}
