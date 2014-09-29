package main

import (
	"fmt"
	"local/world"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	// Initialize GUI
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Listen for input events on keyboard
	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	width, length := termbox.Size()
	width = width/2 - 1
	length = length/2 - 1
	var height int32

	s := &world.Settings{
		NewPeep:          1,      // Initial chance of a new peep being spawned at origin
		MaxAge:           2000,   // Any peep reaching this age will die
		MaxPeeps:         500,    // Absolute max peeps in the world, no more can be born after this.
		RandomDeath:      0.0001, // Chances of random death each turn for every peep
		NewPeepMax:       500,    // Once this many peeps exist, no new ones are spawned from origin
		NewPeepModifier:  1000,   // Controls how often new peeps spawn.  Lower is less often
		Size:             &world.Size{int32(width), int32(length), height, int32(-width), int32(-length), -height},
		SpawnProbability: 0.5, // Chances of two meetings peeps spawning a new one
		TurnTime:         time.Millisecond * 100,
	}
	fmt.Fprintln(os.Stderr, s.Size)

	w := world.NewWorld("Alpha1", *s, event_queue)
	// go w.Run()

	// Advance world every time user hits enter
	// scanner := bufio.NewScanner(os.Stdin)
	w.Show()

	for {
		// scanner.Scan()
		if err := w.NextTurn(); err == nil {
			w.Show()
			time.Sleep(s.TurnTime)
		} else {
			fmt.Println(err)
			os.Exit(0)
		}

	}
}
