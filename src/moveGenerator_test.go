package src

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MoveGenTestCase struct {
	tcName        string
	cp            Position
	expectedMoves []Move
}

func genMoves(start, end Square) Move {
	return Move(start + end<<6)
}

func TestGenerateKnightMoves(t *testing.T) {
	startingPosition, _ := Fen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1").Parse()
	positionAGvsMG1, _ := Fen("2r5/2r2ppp/5k2/B1nBb3/2R5/6P1/P4P1P/3R2K1 b - - 2 35").Parse()

	tcs := []MoveGenTestCase{
		{
			"startingPosition",
			startingPosition,
			[]Move{
				genMoves(b1, a3),
				genMoves(b1, c3),
				genMoves(g1, f3),
				genMoves(g1, h3),
			},
		},
		{
			"positionAGvsMG1",
			positionAGvsMG1,
			[]Move{
				genMoves(c5, a4),
				genMoves(c5, b3),
				genMoves(c5, d3),
				genMoves(c5, e4),
				genMoves(c5, e6),
				genMoves(c5, d7),
				genMoves(c5, b7),
				genMoves(c5, a6),
			},
		},
	}

	for i := 0; i < len(tcs); i++ {
		actualMoves := Knight.generateMoves(tcs[i].cp)
		assert.Equal(t, len(tcs[i].expectedMoves), len(actualMoves))
		assert.ElementsMatch(t, tcs[i].expectedMoves, actualMoves)
	}
}

func TestGenerateRookMoves(t *testing.T) {
	positionAGvsMG, _ := Fen("2r5/2r2ppp/5k2/B1nBb3/2R5/6P1/P4P1P/3R2K1 w - - 2 35").Parse()
	tcs := []MoveGenTestCase{
		{
			"pos1",
			positionAGvsMG,
			[]Move{
				genMoves(d1, a1),
				genMoves(d1, b1),
				genMoves(d1, c1),
				genMoves(d1, e1),
				genMoves(d1, f1),
				genMoves(d1, d2),
				genMoves(d1, d3),
				genMoves(d1, d4),
				genMoves(c4, a4),
				genMoves(c4, b4),
				genMoves(c4, d4),
				genMoves(c4, e4),
				genMoves(c4, f4),
				genMoves(c4, g4),
				genMoves(c4, h4),
				genMoves(c4, c3),
				genMoves(c4, c2),
				genMoves(c4, c1),
				genMoves(c4, c5),
			},
		},
	}

	GenerateSquareMasks()

	for i := 0; i < len(tcs); i++ {
		actualMoves := Rook.generateMoves(tcs[i].cp)
		assert.Equal(t, len(tcs[i].expectedMoves), len(actualMoves))
		assert.ElementsMatch(t, tcs[i].expectedMoves, actualMoves)
	}
}

func TestGenerateQueenMoves(t *testing.T) {
	positionPraggvsMG, _ := Fen("Q7/5p1k/7p/P7/4K3/8/8/3q4 w - - 7 68").Parse()

	tcs := []MoveGenTestCase{
		{
			"pos1",
			positionPraggvsMG,
			[]Move{
				genMoves(a8, a7),
				genMoves(a8, a6),
				genMoves(a8, b8),
				genMoves(a8, c8),
				genMoves(a8, d8),
				genMoves(a8, e8),
				genMoves(a8, f8),
				genMoves(a8, g8),
				genMoves(a8, h8),
				genMoves(a8, b7),
				genMoves(a8, c6),
				genMoves(a8, d5),
			},
		},
	}

	GenerateSquareMasks()

	for i := 0; i < len(tcs); i++ {
		actualMoves := Queen.generateMoves(tcs[i].cp)
		assert.Equal(t, len(tcs[i].expectedMoves), len(actualMoves))
		assert.ElementsMatch(t, tcs[i].expectedMoves, actualMoves)
	}
}

func TestGenerateBishopMoves(t *testing.T) {
	positionAGvsMG, _ := Fen("2r5/2r2ppp/5k2/B1nBb3/2R5/6P1/P4P1P/3R2K1 w - - 2 35").Parse()
	GenerateSquareMasks()

	tcs := []MoveGenTestCase{
		{
			"positionAGvsMG",
			positionAGvsMG,
			[]Move{
				genMoves(a5, b4),
				genMoves(a5, c3),
				genMoves(a5, d2),
				genMoves(a5, e1),
				genMoves(a5, b6),
				genMoves(a5, c7),
				genMoves(d5, e4),
				genMoves(d5, f3),
				genMoves(d5, g2),
				genMoves(d5, h1),
				genMoves(d5, c6),
				genMoves(d5, b7),
				genMoves(d5, a8),
				genMoves(d5, e6),
				genMoves(d5, f7),
			},
		},
	}

	for i := 0; i < len(tcs); i++ {
		actualMoves := Bishop.generateMoves(tcs[i].cp)
		assert.Equal(t, len(tcs[i].expectedMoves), len(actualMoves))
		assert.ElementsMatch(t, tcs[i].expectedMoves, actualMoves)
	}
}

func TestGenerateKingMoves(t *testing.T) {
	startingPosition, _ := Fen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1").Parse()
	positionAGvsMG1, _ := Fen("2r5/2r2ppp/5k2/B1NBb3/2R5/6P1/P4P1P/3R2K1 b - - 2 35").Parse()
	pos1, _ := Fen("r3k1br/8/8/8/8/8/8/R1q1K2R w KQkq - 0 1").Parse()
	pos2, _ := Fen("r3k1br/8/8/8/8/8/8/R1q1K2R b KQkq - 0 1").Parse()

	tcs := []MoveGenTestCase{
		{
			"startingPosition",
			startingPosition,
			[]Move{},
		},
		{
			"positionAGvsMG1",
			positionAGvsMG1,
			[]Move{
				genMoves(f6, e6),
				genMoves(f6, g6),
				genMoves(f6, g5),
				genMoves(f6, f5),
				genMoves(f6, e7),
			},
		},
		{
			"pos1",
			pos1,
			[]Move{
				genMoves(e1, d1),
				genMoves(e1, d2),
				genMoves(e1, e2),
				genMoves(e1, f2),
				genMoves(e1, f1),
				Move(e1 + g1<<6 + Square(CastlingM)),
			},
		},
		{
			"pos2",
			pos2,
			[]Move{
				genMoves(e8, d8),
				genMoves(e8, d7),
				genMoves(e8, e7),
				genMoves(e8, f7),
				genMoves(e8, f8),
				Move(e8 + c8<<6 + Square(CastlingM)),
			},
		},
	}

	for i := 0; i < len(tcs); i++ {
		actualMoves := generateKingMoves(tcs[i].cp)
		assert.Equal(t, len(tcs[i].expectedMoves), len(actualMoves))
		assert.ElementsMatch(t, tcs[i].expectedMoves, actualMoves)
	}
}

func TestGeneratePawnMoves(t *testing.T) {
	whiteMove, _ := Fen("2R1b3/pP1PP1p1/8/PpPp1Pq1/4PppP/5n2/P1p3p1/3Q1b1R w - d6 0 1").Parse()
	fmt.Println(whiteMove)
	blackMove, _ := Fen("2R1b3/pP1PP1p1/8/PpPp1Pq1/4PppP/5n2/P1p3p1/3Q1b1R b - h3 0 1").Parse()
	fmt.Println(blackMove)

	tcs := []MoveGenTestCase{
		{
			"whiteMove",
			whiteMove,
			[]Move{
				genMoves(a2, a3),
				genMoves(a2, a4),
				genMoves(e4, d5),
				genMoves(e4, e5),
				genMoves(h4, h5),
				genMoves(h4, g5),
				genMoves(a5, a6),
				genMoves(c5, c6),
				genMoves(c5+Square(EnPassant), d6),
				genMoves(f5, f6),
				genMoves(b7+Square(Promotion)+Square(QueenPromotion), b8),
				genMoves(b7+Square(Promotion)+Square(RookPromotion), b8),
				genMoves(b7+Square(Promotion)+Square(BishopPromotion), b8),
				genMoves(b7+Square(Promotion)+Square(KnightPromotion), b8),
				genMoves(d7+Square(Promotion)+Square(QueenPromotion), d8),
				genMoves(d7+Square(Promotion)+Square(RookPromotion), d8),
				genMoves(d7+Square(Promotion)+Square(BishopPromotion), d8),
				genMoves(d7+Square(Promotion)+Square(KnightPromotion), d8),
				genMoves(d7+Square(Promotion)+Square(QueenPromotion), e8),
				genMoves(d7+Square(Promotion)+Square(RookPromotion), e8),
				genMoves(d7+Square(Promotion)+Square(BishopPromotion), e8),
				genMoves(d7+Square(Promotion)+Square(KnightPromotion), e8),
			},
		},
		{
			"blackMove",
			blackMove,
			[]Move{
				genMoves(c2+Square(Promotion)+Square(QueenPromotion), c1),
				genMoves(c2+Square(Promotion)+Square(RookPromotion), c1),
				genMoves(c2+Square(Promotion)+Square(BishopPromotion), c1),
				genMoves(c2+Square(Promotion)+Square(KnightPromotion), c1),
				genMoves(c2+Square(Promotion)+Square(QueenPromotion), d1),
				genMoves(c2+Square(Promotion)+Square(RookPromotion), d1),
				genMoves(c2+Square(Promotion)+Square(BishopPromotion), d1),
				genMoves(c2+Square(Promotion)+Square(KnightPromotion), d1),
				genMoves(g2+Square(Promotion)+Square(QueenPromotion), g1),
				genMoves(g2+Square(Promotion)+Square(RookPromotion), g1),
				genMoves(g2+Square(Promotion)+Square(BishopPromotion), g1),
				genMoves(g2+Square(Promotion)+Square(KnightPromotion), g1),
				genMoves(g2+Square(Promotion)+Square(QueenPromotion), h1),
				genMoves(g2+Square(Promotion)+Square(RookPromotion), h1),
				genMoves(g2+Square(Promotion)+Square(BishopPromotion), h1),
				genMoves(g2+Square(Promotion)+Square(KnightPromotion), h1),
				genMoves(g4, g3),
				genMoves(g4+Square(EnPassant), h3),
				genMoves(b5, b4),
				genMoves(d5, d4),
				genMoves(d5, e4),
				genMoves(a7, a6),
				genMoves(g7, g6),
			},
		},
	}

	for i := 0; i < len(tcs); i++ {
		actualMoves := generatePawnMoves(tcs[i].cp)
		assert.Equal(t, len(tcs[i].expectedMoves), len(actualMoves))
		assert.ElementsMatch(t, tcs[i].expectedMoves, actualMoves)
	}
}
