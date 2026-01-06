//go:build !js

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	p1Type := flag.String("p1", "human", "Player 1 type (human/mcts)")
	p2Type := flag.String("p2", "human", "Player 2 type (human/mcts)")
	p3Type := flag.String("p3", "human", "Player 3 type (human/mcts)")
	iterations := flag.Int("iterations", 1000, "MCTS iterations")
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	seed := flag.Int64("seed", 0, "Random seed (0 for time-based)")
	flag.Parse()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "could not start CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}
	if *seed == 0 {
		xorState = uint64(time.Now().UnixNano())
	} else {
		xorState = uint64(*seed)
	}
	if xorState == 0 {
		xorState = 1
	}
	game := NewSquavaGame()
	createPlayer := func(t, name, symbol string, id int) Player {
		if t == "mcts" {
			p := NewMCTSPlayer(name, symbol, id, *iterations)
			p.Verbose = true
			return p
		}
		return NewHumanPlayer(name, symbol, id)
	}
	game.AddPlayer(createPlayer(*p1Type, "Player 1", "X", 0))
	game.AddPlayer(createPlayer(*p2Type, "Player 2", "O", 1))
	game.AddPlayer(createPlayer(*p3Type, "Player 3", "Z", 2))
	game.Run()
}
