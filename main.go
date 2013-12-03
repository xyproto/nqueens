package main

/*
 * The N queens problem
 * solved with GA
 * one byte per queen, for the placement.
 * the placement is a number that only counts the free spaces,
 * not the occupied ones.
 * That is: queen 1 can be placed at any position from 0 up to N*N
 * queen 2 can be placed at any position that is available from 0
 * up to N*N etc. The function that tests the fitness keeps track
 * of which positions are occupied or not. The solution is a list of
 * bytes, one byte per queen.
 */

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	FREE           PosType = 0
	QUEEN          PosType = 1
	COVERED        PosType = 2
	QUEENS         uint    = 20
	BOARDSIZE      uint    = 20
	POPSIZE        uint    = 2000
	MAXGENERATIONS uint    = 3000
)

/* All positions on the board, from 0 up to N*N */
type BoardPosIndex uint

/* All free positions on the board, from 0 up to N*N */
type FreePosIndex uint

/* A solution is N index positions */
type Solution []FreePosIndex

/* A population is a collection of solutions */
type Population []Solution

/* A free space, queen or "covered" by a queen */
type PosType byte

/* A board is N*N positions */
type Board struct {
	width uint
	data  []PosType
}

type PopulationFitness []float64

type Position struct {
	x uint
	y uint
}

func NewPopulationFitness(popsize uint) PopulationFitness {
	return make(PopulationFitness, popsize)
}

func NewBoard(w uint) *Board {
	var b Board
	b.width = w
	b.data = make([]PosType, w*w)
	return &b
}

func NewSolution(numqueens uint) Solution {
	return make([]FreePosIndex, numqueens)
}

func NewRandomSolution(numqueens uint, maxboardpos BoardPosIndex) Solution {
	sol := NewSolution(numqueens)
	var i uint
	for i = 0; i < numqueens; i++ {
		sol.set(i, FreePosIndex(rand.Intn(int(maxboardpos))))
	}
	return sol
}

func pos2xy(pos BoardPosIndex) *Position {
	var i BoardPosIndex = 0
	var p Position
	for p.y = 0; p.y < BOARDSIZE; p.y++ {
		for p.x = 0; p.x < BOARDSIZE; p.x++ {
			if i == pos {
				return &Position{p.x, p.y}
			}
			i++
		}
	}
	return &Position{255, 255}
}

/* Place a queen at the Nth free position on the board
 * Returns the board position as an index and possibly an error
 */
func (b *Board) place(targetpos FreePosIndex) (BoardPosIndex, error) {
	var freepos FreePosIndex
	var usepos BoardPosIndex

	width := b.width
	height := width
	maxx := width - 1

	for usepos = 0; usepos < BoardPosIndex(width*width); usepos++ {
		if (targetpos == freepos) && (b.data[usepos] == FREE) {
			// Mark the row, column and the two diagonals too
			o := pos2xy(usepos)
			//fmt.Println("xy", o.x, o.y)
			var x, y uint
			for y = 0; y < height; y++ {
				for x = 0; x < width; x++ {
					/* Mark horizontal and vertical lines as COVERED */
					if o.x == x || o.y == y {
						b.data[y*width+x] = COVERED
					}
					/* Mark diagonal lines from upper left to lower right as COVERED */
					if x == y {
						diagx := (x + o.x) - o.y
						diagy := y
						if (diagx >= 0 && diagy >= 0) && (diagx < width && diagy < height) {
							b.data[diagy*width+diagx] = COVERED
						}
					}
					/* Mark diagonal lines from upper right to lower left as COVERED */
					if (maxx - x) == y {
						diagx := (x - (maxx - o.x)) + o.y
						diagy := (y + o.y) - o.y
						if (diagx >= 0 && diagy >= 0) && (diagx < BOARDSIZE && diagy < height) {
							b.data[diagy*width+diagx] = COVERED
						}
					}

				}
			}
			// Mark the queen
			b.data[usepos] = QUEEN
			return usepos, nil
		}
		if b.data[usepos] == FREE {
			freepos++
		}
	}
	return 0, errors.New("No available position")
}

func (sol *Solution) generateBoard() (*Board, uint) {
	board := NewBoard(BOARDSIZE)
	var queenCounter uint
	for _, queenposindex := range *sol {
		if _, err := board.place(queenposindex); err == nil {
			// Could place queen, increase queenCounter
			queenCounter++
		}
	}
	return board, queenCounter
}

func (b *Board) String() string {
	var (
		s string
		t PosType
		y uint
		x uint
	)
	width := b.width
	height := width
	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			t = b.data[y*width+x]
			if t == FREE {
				s += " "
			} else if t == QUEEN {
				s += "q"
			} else if t == COVERED {
				s += "."
			}
		}
		s += "\n"
	}
	return s + "\n"
}

func (sol Solution) String() string {
	board, _ := sol.generateBoard()
	return board.String()
}

func (sol Solution) set(i uint, freepos FreePosIndex) {
	sol[i] = freepos
}

func (sol *Solution) fitness() float64 {
	_, numqueens := sol.generateBoard()
	return float64(numqueens) / float64(QUEENS)
}

func NewPopulation(size uint) Population {
	t := time.Now()
	rand.Seed(t.UnixNano())
	pop := make([]Solution, size)
	var i uint
	for i = 0; i < size; i++ {
		pop[i] = NewRandomSolution(QUEENS, BoardPosIndex(BOARDSIZE*BOARDSIZE))
	}
	return pop
}

func test_solution() {
	sol := NewSolution(QUEENS)
	sol.set(0, 20)
	sol.set(1, 2)
	sol.set(2, 2)
	fmt.Println(sol)
	//fmt.Println("fitness:", sol.fitness())
}

func crossover(a, b Solution, point uint, numqueens uint) Solution {
	c := NewSolution(numqueens)
	var i uint
	for i = 0; i < point; i++ {
		c[i] = a[i]
	}
	for i = point; i < numqueens; i++ {
		c[i] = b[i]
	}
	return c
}

func (sol Solution) mutate(numqueens uint, maxboardpos BoardPosIndex) {
	randpos := rand.Intn(int(numqueens))
	sol[randpos] = FreePosIndex(rand.Intn(int(maxboardpos)))
}

func sum(scores []float64) float64 {
	var total float64
	for _, score := range scores {
		total += score
	}
	return total
}

func main() {
	bestfitnessindex := 0
	var popsize uint = POPSIZE
	pop := NewPopulation(popsize)
	var generation uint
	var average float64
	for generation = 0; generation < MAXGENERATIONS; generation++ {
		fmt.Println("Generation", generation)
		fit := NewPopulationFitness(popsize)
		for i, s := range pop {
			fit[i] = s.fitness()
		}
		//fmt.Println(fit)
		total := sum(fit)
		fmt.Println("total =", total)
		average = total / float64(popsize)
		fmt.Println("average =", average)
		var bestfitness float64 = 0.0
		var nextbestfitness float64 = 0.0
		nextbestfitnessindex := 0
		for i, _ := range fit {
			if fit[i] >= bestfitness {
				nextbestfitness = bestfitness
				bestfitness = fit[i]
				nextbestfitnessindex = bestfitnessindex
				bestfitnessindex = i
			}
		}
		fmt.Println("best =", bestfitness)
		fmt.Println("nextbest =", nextbestfitness)
		if bestfitness == 1.0 {
			fmt.Println("Found fitness 1")
			break
		}
		var mutrate float64 = 0.0
		var crossrate float64 = 0.1
		var newpoprate float64 = 0.0
		for i, _ := range pop {
			fitness := fit[i]
			if average > 0.7 && fitness < 0.5 {
				pop[i] = NewSolution(QUEENS)
			} else if average > 0.8 && fitness < 0.6 {
				pop[i] = NewSolution(QUEENS)
			} else if average > 0.9 && fitness < 0.7 {
				pop[i] = NewSolution(QUEENS)
			} else if fitness < (average * 0.3) {
				// 50% chance of being replaced with randomness
				if rand.Float64() <= 0.5 {
					pop[i] = NewSolution(QUEENS)
				}
			}
			if bestfitness > average {
				// slow down the mutation rate
				mutrate = 0.15
				crossrate = 0.07
			} else {
				mutrate = 0.4
				crossrate = 0.4
			}
			if bestfitness == nextbestfitness {
				mutrate *= 3.0
			}
			if average > 0.9 {
				newpoprate = 0.4
			} else {
				newpoprate = 0.2
			}
			// An advantage for the best ones
			if fitness > (bestfitness * 0.9) {
				if rand.Float64() <= 0.9 {
					continue
				}
			}
			// A certain chance for mutation
			if rand.Float64() <= mutrate {
				pop[rand.Intn(int(popsize))].mutate(QUEENS, BoardPosIndex(BOARDSIZE*BOARDSIZE))
			}
			// A certain chance for crossover
			if rand.Float64() <= crossrate {
				crossoverpoint := uint(rand.Intn(int(QUEENS)))
				pop[i] = crossover(pop[bestfitnessindex], pop[nextbestfitnessindex], crossoverpoint, QUEENS)
			}
			// A certain chance for new random variations
			if rand.Float64() <= newpoprate {
				pop[i] = NewSolution(QUEENS)
			}
		}
	}
	fmt.Println("generation", generation)
	fmt.Println(pop[bestfitnessindex])
}
