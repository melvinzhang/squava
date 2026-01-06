//go:build !wasm

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

// --- Human Player ---
type HumanPlayer struct {
	info PlayerInfo
}

func NewHumanPlayer(name, symbol string, id int) *HumanPlayer {
	return &HumanPlayer{info: PlayerInfo{name: name, symbol: symbol, id: id}}
}
func (h *HumanPlayer) Name() string   { return h.info.name }
func (h *HumanPlayer) Symbol() string { return h.info.symbol }
func (h *HumanPlayer) ID() int        { return h.info.id }
func (h *HumanPlayer) GetMove(board Board, players []int, turnIdx int) Move {
	forcedMoves := GetForcedMoves(board, players, turnIdx)
	reader := bufio.NewReader(os.Stdin)
	for {
		prompt := fmt.Sprintf("%s (%s), enter your move (e.g., A1): ", h.info.name, h.info.symbol)
		if forcedMoves != 0 {
			forcedStr := []string{}
			temp := forcedMoves
			for temp != 0 {
				idx := bits.TrailingZeros64(uint64(temp))
				m := MoveFromIndex(idx)
				forcedStr = append(forcedStr, fmt.Sprintf("%c%d", int(m.c)+65, int(m.r)+1))
				temp &= Bitboard(^(uint64(1) << idx))
			}
			fmt.Printf("FORCED MOVE! You must block the next player. Valid moves: %s\n", strings.Join(forcedStr, ", "))
		}
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToUpper(input))
		r, c, err := parseInput(input)
		if err != nil {
			fmt.Println("Invalid format. Use algebraic (A1).")
			continue
		}
		if !isValidCoord(r, c) {
			fmt.Println("Move out of bounds.")
			continue
		}
		idx := r*8 + c
		mask := uint64(1) << idx
		if (uint64(board.Occupied) & mask) != 0 {
			fmt.Println("Cell already occupied.")
			continue
		}
		move := Move{int8(r), int8(c)}
		if forcedMoves != 0 && (forcedMoves&(Bitboard(1)<<idx)) == 0 {
			fmt.Println("Invalid move. You must block the opponent or win immediately.")
			continue
		}
		return move
	}
}
func parseInput(inp string) (int, int, error) {
	if len(inp) < 2 {
		return 0, 0, fmt.Errorf("invalid length")
	}
	colChar := inp[0]
	rowStr := inp[1:]
	if colChar < 'A' || colChar > 'H' {
		return 0, 0, fmt.Errorf("invalid column")
	}
	col := int(colChar - 'A')
	row, err := strconv.Atoi(rowStr)
	if err != nil {
		return 0, 0, err
	}
	return row - 1, col, nil
}
func isValidCoord(r, c int) bool {
	return r >= 0 && r < BoardSize && c >= 0 && c < BoardSize
}

// --- Game Engine ---
type SquavaGame struct {
	gs      GameState
	players []Player
}

func NewSquavaGame() *SquavaGame {
	return &SquavaGame{
		gs: GameState{WinnerID: -1},
	}
}

func (g *SquavaGame) AddPlayer(p Player) {
	g.players = append(g.players, p)
}

func (g *SquavaGame) GetPlayer(id int) Player {
	for _, p := range g.players {
		if p.ID() == id {
			return p
		}
	}
	return nil
}

func (g *SquavaGame) PrintBoard() {
	fmt.Print("   ")
	for i := 0; i < BoardSize; i++ {
		fmt.Printf("%c ", 'A'+i)
	}
	fmt.Println()
	for r := 0; r < BoardSize; r++ {
		fmt.Printf("%2d ", r+1)
		for c := 0; c < BoardSize; c++ {
			symbol := "."
			idx := r*8 + c
			mask := Bitboard(uint64(1) << idx)
			if (g.gs.Board.P[0] & mask) != 0 {
				symbol = "X"
			} else if (g.gs.Board.P[1] & mask) != 0 {
				symbol = "O"
			} else if (g.gs.Board.P[2] & mask) != 0 {
				symbol = "Z"
			}
			fmt.Printf("%s ", symbol)
		}
		fmt.Println()
	}
}

func (g *SquavaGame) Run() {
	fmt.Println("Starting 3-Player Squava!")
	fmt.Printf("Random Seed: %d\n", xorState)
	fmt.Println("Board Size: 8x8")
	fmt.Println("Rules: 4-in-a-row wins. 3-in-a-row loses.")

	activeMask := uint8(0)
	for _, p := range g.players {
		activeMask |= 1 << uint(p.ID())
	}
	g.gs = NewGameState(g.gs.Board, g.players[0].ID(), activeMask)

	moveCount := 1
	for {
		winnerID, ok := g.gs.IsTerminal()
		if ok {
			g.PrintBoard()
			if winnerID != -1 {
				isWin, _ := CheckBoard(g.gs.Board.P[winnerID])
				if isWin {
					fmt.Printf("Result: %s Wins (4-in-a-row)\n", g.GetPlayer(winnerID).Name())
				} else {
					fmt.Printf("Result: %s Wins (Last Standing)\n", g.GetPlayer(winnerID).Name())
				}
			} else {
				fmt.Println("Result: Draw")
			}
			return
		}

		currentPlayer := g.GetPlayer(g.gs.PlayerID)
		g.PrintBoard()
		fmt.Printf("Move %d: %s (%s)\n", moveCount, currentPlayer.Name(), currentPlayer.Symbol())

		if _, ok := currentPlayer.(*MCTSPlayer); ok {
			fmt.Printf("%s is thinking...\n", currentPlayer.Name())
		}

		activeIDs := g.gs.ActiveIDs()
		var turnIdx int
		for i, id := range activeIDs {
			if id == g.gs.PlayerID {
				turnIdx = i
				break
			}
		}

		move := currentPlayer.GetMove(g.gs.Board, activeIDs, turnIdx)

		if _, ok := currentPlayer.(*MCTSPlayer); ok {
			fmt.Printf("%s chooses %c%d\n", currentPlayer.Name(), int(move.c)+65, int(move.r)+1)
		}

		prevMask := g.gs.ActiveMask
		g.gs.ApplyMove(move)
		moveCount++

		if g.gs.ActiveMask != prevMask {
			// Find who was eliminated
			eliminatedID := -1
			for i := 0; i < 3; i++ {
				if (prevMask&(1<<uint(i))) != 0 && (g.gs.ActiveMask&(1<<uint(i))) == 0 {
					eliminatedID = i
					break
				}
			}
			fmt.Printf("Result: %s Eliminated (3-in-a-row)\n", g.GetPlayer(eliminatedID).Name())
		}
	}
}
