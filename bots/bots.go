package bots

/*
Represents heuristic values for certain squares on a Board

From: https://courses.cs.washington.edu/courses/cse573/04au/Project/mini1/RUSSIA/Final_Paper.pdf
*/
var StaticWeights = [...]int{
	4, -3, 2, 2, 2, 2, -3, 4, -3, -4, -1, -1, -1, -1, -4, -3, 2, -1, 1, 0, 0,
	1, -1, 2, 2, -1, 0, 1, 1, 0, -1, 2, 2, -1, 0, 1, 1, 0, -1, 2, 2, -1, 1, 0,
	0, 1, -1, 2, -3, -4, -1, -1, -1, -1, -4, -3, 4, -3, 2, 2, 2, 2, -3, 4,
}

// Get a pseudorandom move for the Randy bot
func RandyMove(moves map[int]bool) int {
	for k := range moves {
		return k
	}
	return 0
}
