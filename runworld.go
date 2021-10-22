package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/DanTulovsky/world"
	termbox "github.com/nsf/termbox-go"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	debug      = flag.Bool("debug", false, "If true, turns don't advanced automatically.")
)

func main() {
	flag.Parse()

	// profiling support
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// end profiling

	// Initialize GUI
	err := termbox.Init()
	if err != nil {
		// panic(err)
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
		NewPeep:                1,    // Initial chance of a new peep being spawned at spawn points
		MaxAge:                 3000, // Any peep reaching this age will die
		MaxPeeps:               1000, // Absolute max peeps in the world, no more can be born after this.
		RandomDeath:            0,    // Chances of random death each turn for every peep
		NewPeepMax:             3000, // Once this many peeps exist, no new ones are spawned from spawn points
		NewPeepModifier:        100,  // Controls how often new peeps spawn.  Lower is less often
		Size:                   &world.Size{int32(width), int32(length), height, int32(-width), int32(-length), -height},
		SpawnProbability:       1, // Chances of two meetings peeps spawning a new one
		TurnTime:               time.Millisecond * 10,
		PeepRememberTurns:      2000, // can remember what's around them for X turns; right now they look before moving though
		PeepViewDistance:       20,   // can see this many squares away
		KillIfSurroundByOther:  true,
		KillIfSurroundedBySame: true,
		KillIfSurrounded:       false,
		MaxGenders:             2, // 1 - 4
	}
	s.PeepSpawnInterval = 0   //world.Turn(s.MaxAge / 10) // Spawn 10 times in a life time.
	s.YoungHightlightAge = 10 // Highlighted in the GUI while young
	s.SpawnAge = 30           // s.YoungHightlightAge + 100 // s.MaxAge / 10 // Can spawn after this age

	w := world.NewWorld("Alpha1", *s, event_queue, *debug)
	w.Run() // starts https server, other things later

	// Set homebase locations for each gender
	locations := w.SpawnLocations()
	var usedLocations []world.Location // avoid spawning in the same place
	var spawnLocation world.Location

	for _, gender := range w.Genders() {
		if len(usedLocations) != len(locations) {
			for {
				spawnLocation = locations[rand.Intn(len(locations))]
				if world.ListContains(usedLocations, spawnLocation) {
					continue // pick another one
				}
				usedLocations = append(usedLocations, spawnLocation)
				break
			}
		}
		world.Log(gender, spawnLocation)
		w.SetHomebase(gender, spawnLocation)
	}
	w.Show(os.Stderr)

	for {
		if err := w.NextTurn(); err == nil {
			// w.Show()
			time.Sleep(s.TurnTime)
		} else {
			fmt.Println(err)
			break
		}

	}
}
