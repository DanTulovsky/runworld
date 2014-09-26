package main

import (
	"bufio"
	"local/world"
	"os"
	"time"
)

func main() {
	s := &world.Settings{
		NewPeep:         1,
		MaxAge:          999,
		RandomDeath:     0.0001,
		NewPeepModifier: 1000,
	}
	w := world.NewWorld("Alpha1", *s)
	go w.Run()

	// Advance world every time user hits enter
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		w.NextTurn()
		w.Show()
		time.Sleep(time.Microsecond * 3000)
	}
}
