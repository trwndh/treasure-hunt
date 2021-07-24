# Treasure Hunt
This is simple command line game, built using Golang

## Requirements
`- Golang`

## How To Run
Just run `$ make play` or `$ go run main.go` to start the game

## How To Play
Based on given grid, you will have a board game that contains 4 elements: 
```
- # -> as an obstacle
- . (dot) -> as the clearpath, you can only move to these path
- X -> as your current position
- $ -> as possibility treasure location
```

You need to move the `X` to find the treasure. 
To move, you need to input `directions [string]` and `steps [integer]` to the input. 

available directions:
```
1. Up / North 
2. Down / South
3. Right / East
```

You only have 3 steps available to find Treasure.

	1. First you can only move 'up / north', 
	2. After that you can only move 'right / right', 
	3. Last step you can only move 'down / south'

If after these 3 steps you didn't find treasure, then Game is over.

---
### Example
```
the board look like this:

    0 1 2 3 4 5 6 7
0: [# # # # # # # #]
1: [# . . . . . . #]
2: [# . # # # . . #]
3: [# . $ $ # . # #]
4: [# X # $ . $ $ #]
5: [# # # # # # # #]
```
As stated before, you need to go to one of `$` location. Give the input `up 1` to move your position (mark by `X`) to up 1 step.
```
The result will look like this:
    0 1 2 3 4 5 6 7
0: [# # # # # # # #]
1: [# . . . . . . #]
2: [# . # # # . . #]
3: [# X $ $ # . # #]
4: [# . # $ . $ $ #]
5: [# # # # # # # #]

```
 

## Quit the game

input `q` to quit the game. Or you can use the mighty `Ctrl + C`.