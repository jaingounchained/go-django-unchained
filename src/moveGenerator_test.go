package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MoveGenTestCase struct {
	tcName        string
	cp            Position
	expectedMoves []Move
}

func TestGenerateKnightMoves(t *testing.T) {
	startingPosition, _ := Fen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1").Parse()
	positionAGvsMG1, _ := Fen("2r5/2r2ppp/5k2/B1nBb3/2R5/6P1/P4P1P/3R2K1 b - - 2 35").Parse()

	tcs := []MoveGenTestCase{
		{
			"startingPosition",
			startingPosition,
			[]Move{
				squaresToMove(b1, a3, Normal),
				squaresToMove(b1, c3, Normal),
				squaresToMove(g1, f3, Normal),
				squaresToMove(g1, h3, Normal),
			},
		},
		{
			"positionAGvsMG1",
			positionAGvsMG1,
			[]Move{
				squaresToMove(c5, a4, Normal),
				squaresToMove(c5, b3, Normal),
				squaresToMove(c5, d3, Normal),
				squaresToMove(c5, e4, Normal),
				squaresToMove(c5, e6, Normal),
				squaresToMove(c5, d7, Normal),
				squaresToMove(c5, b7, Normal),
				squaresToMove(c5, a6, Normal),
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
				squaresToMove(d1, a1, Normal),
				squaresToMove(d1, b1, Normal),
				squaresToMove(d1, c1, Normal),
				squaresToMove(d1, e1, Normal),
				squaresToMove(d1, f1, Normal),
				squaresToMove(d1, d2, Normal),
				squaresToMove(d1, d3, Normal),
				squaresToMove(d1, d4, Normal),
				squaresToMove(c4, a4, Normal),
				squaresToMove(c4, b4, Normal),
				squaresToMove(c4, d4, Normal),
				squaresToMove(c4, e4, Normal),
				squaresToMove(c4, f4, Normal),
				squaresToMove(c4, g4, Normal),
				squaresToMove(c4, h4, Normal),
				squaresToMove(c4, c3, Normal),
				squaresToMove(c4, c2, Normal),
				squaresToMove(c4, c1, Normal),
				squaresToMove(c4, c5, Capture),
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
				squaresToMove(a8, a7, Normal),
				squaresToMove(a8, a6, Normal),
				squaresToMove(a8, b8, Normal),
				squaresToMove(a8, c8, Normal),
				squaresToMove(a8, d8, Normal),
				squaresToMove(a8, e8, Normal),
				squaresToMove(a8, f8, Normal),
				squaresToMove(a8, g8, Normal),
				squaresToMove(a8, h8, Normal),
				squaresToMove(a8, b7, Normal),
				squaresToMove(a8, c6, Normal),
				squaresToMove(a8, d5, Normal),
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
				squaresToMove(a5, b4, Normal),
				squaresToMove(a5, c3, Normal),
				squaresToMove(a5, d2, Normal),
				squaresToMove(a5, e1, Normal),
				squaresToMove(a5, b6, Normal),
				squaresToMove(a5, c7, Capture),
				squaresToMove(d5, e4, Normal),
				squaresToMove(d5, f3, Normal),
				squaresToMove(d5, g2, Normal),
				squaresToMove(d5, h1, Normal),
				squaresToMove(d5, c6, Normal),
				squaresToMove(d5, b7, Normal),
				squaresToMove(d5, a8, Normal),
				squaresToMove(d5, e6, Normal),
				squaresToMove(d5, f7, Capture),
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
				squaresToMove(f6, e6, Normal),
				squaresToMove(f6, g6, Normal),
				squaresToMove(f6, g5, Normal),
				squaresToMove(f6, f5, Normal),
				squaresToMove(f6, e7, Normal),
			},
		},
		{
			"pos1",
			pos1,
			[]Move{
				squaresToMove(e1, d1, Normal),
				squaresToMove(e1, d2, Normal),
				squaresToMove(e1, e2, Normal),
				squaresToMove(e1, f2, Normal),
				squaresToMove(e1, f1, Normal),
				WhiteKingSideCastling,
			},
		},
		{
			"pos2",
			pos2,
			[]Move{
				squaresToMove(e8, d8, Normal),
				squaresToMove(e8, d7, Normal),
				squaresToMove(e8, e7, Normal),
				squaresToMove(e8, f7, Normal),
				squaresToMove(e8, f8, Normal),
				BlackQueenSideCastling,
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
	blackMove, _ := Fen("2R1b3/pP1PP1p1/8/PpPp1Pq1/4PppP/5n2/P1p3p1/3Q1b1R b - h3 0 1").Parse()

	tcs := []MoveGenTestCase{
		{
			"whiteMove",
			whiteMove,
			[]Move{
				squaresToMove(a2, a3, Normal),
				squaresToMove(a2, a4, DoublePawnPush),
				squaresToMove(e4, d5, Capture),
				squaresToMove(e4, e5, Normal),
				squaresToMove(h4, h5, Normal),
				squaresToMove(h4, g5, Capture),
				squaresToMove(a5, a6, Normal),
				squaresToMove(c5, c6, Normal),
				squaresToMove(c5, d6, EnPassant),
				squaresToMove(f5, f6, Normal),
				squaresToMove(b7, b8, QueenPromotionNormal),
				squaresToMove(b7, b8, RookPromotionNormal),
				squaresToMove(b7, b8, BishopPromotionNormal),
				squaresToMove(b7, b8, KnightPromotionNormal),
				squaresToMove(d7, d8, QueenPromotionNormal),
				squaresToMove(d7, d8, RookPromotionNormal),
				squaresToMove(d7, d8, BishopPromotionNormal),
				squaresToMove(d7, d8, KnightPromotionNormal),
				squaresToMove(d7, e8, QueenPromotionCapture),
				squaresToMove(d7, e8, BishopPromotionCapture),
				squaresToMove(d7, e8, KnightPromotionCapture),
				squaresToMove(d7, e8, RookPromotionCapture),
			},
		},
		{
			"blackMove",
			blackMove,
			[]Move{
				squaresToMove(c2, c1, QueenPromotionNormal),
				squaresToMove(c2, c1, RookPromotionNormal),
				squaresToMove(c2, c1, BishopPromotionNormal),
				squaresToMove(c2, c1, KnightPromotionNormal),
				squaresToMove(c2, d1, QueenPromotionCapture),
				squaresToMove(c2, d1, RookPromotionCapture),
				squaresToMove(c2, d1, BishopPromotionCapture),
				squaresToMove(c2, d1, KnightPromotionCapture),
				squaresToMove(g2, g1, QueenPromotionNormal),
				squaresToMove(g2, g1, RookPromotionNormal),
				squaresToMove(g2, g1, BishopPromotionNormal),
				squaresToMove(g2, g1, KnightPromotionNormal),
				squaresToMove(g2, h1, QueenPromotionCapture),
				squaresToMove(g2, h1, RookPromotionCapture),
				squaresToMove(g2, h1, BishopPromotionCapture),
				squaresToMove(g2, h1, KnightPromotionCapture),
				squaresToMove(g4, g3, Normal),
				squaresToMove(g4, h3, EnPassant),
				squaresToMove(b5, b4, Normal),
				squaresToMove(d5, d4, Normal),
				squaresToMove(d5, e4, Capture),
				squaresToMove(a7, a6, Normal),
				squaresToMove(g7, g6, Normal),
			},
		},
	}

	for i := 0; i < len(tcs); i++ {
		actualMoves := generatePawnMoves(tcs[i].cp)
		assert.Equal(t, len(tcs[i].expectedMoves), len(actualMoves))
		assert.ElementsMatch(t, tcs[i].expectedMoves, actualMoves)
	}
}
