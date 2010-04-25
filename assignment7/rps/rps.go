package rps

import "games"
import "os"

var (
	player1 *games.View
	player2 *games.View
)

func Play(path string) {
	player1 = games.NewLocalView("A", os.Stdin, os.Stdout)

	file, _ := os.Open(path, os.O_RDWR, 0)
	player2 = games.NewLocalView("B", file, file)
} 
