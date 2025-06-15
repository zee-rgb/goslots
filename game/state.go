package game

type GameState struct {
	Symbols     map[string]uint
	Multipliers map[string]uint
	SymbolArr   []string
	Balance     uint
}

func NewGameState() *GameState {
	return &GameState{
		Symbols: map[string]uint{
			"A": 4,
			"B": 7,
			"C": 12,
			"D": 20,
		},
		Multipliers: map[string]uint{
			"A": 20,
			"B": 10,
			"C": 5,
			"D": 2,
		},
		Balance: 200,
	}
}

func (g *GameState) InitSymbolArray() {
	g.SymbolArr = make([]string, 0)
	for symbol, count := range g.Symbols {
		for i := uint(0); i < count; i++ {
			g.SymbolArr = append(g.SymbolArr, symbol)
		}
	}
}
