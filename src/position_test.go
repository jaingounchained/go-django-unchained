package src

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

type UpdatePositionTestCase struct {
	desc                     string
	move                     Move
	initialPositionFen       Fen
	expectedFinalPositionFen Fen
}

// func TestUpdatePosition_Normal(t *testing.T) {
// 	tcs := []UpdatePositionTestCase{
// 		{
// 			"1",
// 			Normal + squaresToMove(d2, d4),
// 			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
// 			"rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1",
// 		},
// 		{
// 			"2",
// 			Normal + squaresToMove(d4, e5),
// 			"rnbqkbnr/pppp1ppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 2",
// 			"rnbqkbnr/pppp1ppp/8/4P3/8/8/PPP1PPPP/RNBQKBNR b KQkq - 0 2",
// 		},
// 	}

// 	for _, tc := range tcs {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			initialPosition, _ := tc.initialPositionFen.Parse()
// 			expectedFinalPosition, _ := tc.expectedFinalPositionFen.Parse()
// 			actualFinalPosition, _ := initialPosition.makeMove(tc.move)
// 			assert.Equal(t, actualFinalPosition.activeColor, expectedFinalPosition.activeColor)
// 			assert.Equal(t, actualFinalPosition.piecePlacement, expectedFinalPosition.piecePlacement)
// 		})
// 	}
// }

// func TestUpdatePosition_Promotion(t *testing.T) {
// 	tcs := []UpdatePositionTestCase{
// 		{
// 			"1",
// 			Promotion + QueenPromotion + squaresToMove(d7, d8),
// 			"2K1R3/R2P2k1/8/p7/8/3n2q1/1P6/6r1 w - - 0 1",
// 			"2KQR3/R5k1/8/p7/8/3n2q1/1P6/6r1 b - - 0 1",
// 		},
// 		{
// 			"2",
// 			Promotion + BishopPromotion + squaresToMove(d7, e8),
// 			"4r3/RK1P2k1/8/p7/8/3n2q1/1P6/6r1 w - - 0 1",
// 			"4B3/RK4k1/8/p7/8/3n2q1/1P6/6r1 b - - 0 1",
// 		},
// 	}

// 	for _, tc := range tcs {
// 		t.Run(tc.desc, func(t *testing.T) {
// 			initialPosition, _ := tc.initialPositionFen.Parse()
// 			expectedFinalPosition, _ := tc.expectedFinalPositionFen.Parse()
// 			actualFinalPosition, _ := initialPosition.makeMove(tc.move)
// 			assert.Equal(t, actualFinalPosition.activeColor, expectedFinalPosition.activeColor)
// 			assert.Equal(t, actualFinalPosition.piecePlacement, expectedFinalPosition.piecePlacement)
// 		})
// 	}
// }

// // write test cases for castling & enpassant

func TestCalculateOurKingDangerSquares(t *testing.T) {
	// KDSTC = King Danger Squares Test Cases
	type KDSTC struct {
		desc          string
		positionFen   Fen
		expectedKDSBb Bitboard      // expected King Danger Squares Bitboard
		expectedKCs   []KingChecker // expected King Checkers
	}

	GenerateNonSlidingPieceTypeAttackingSquares()
	GenerateSquareMasks()

	tcs := []KDSTC{
		{
			"1",
			"4k3/8/8/5R2/8/8/8/4K3 b - - 0 1",
			0x202020df20203828,
			[]KingChecker{},
		},
		{
			"2",
			"8/4k3/8/4R3/8/8/8/4K3 b - - 0 1",
			0x101010ef10103838,
			[]KingChecker{
				{
					Rook,
					Bitboard(0x1000000000),
				},
			},
		},
		{
			"3",
			"2r5/2r2ppp/5k2/B1nBb3/2R5/6P1/P4P1P/3R2K1 b - - 2 35",
			0x126160cff7eecf7,
			[]KingChecker{},
		},
		{
			"4",
			"4k3/8/5N2/8/8/8/8/4K3 b - - 0 1",
			0x5088008850003828,
			[]KingChecker{
				{
					Knight,
					Bitboard(0x200000000000),
				},
			},
		},
		{
			"5",
			"4k3/6N1/5b2/4R3/8/8/8/4K3 b - - 0 1",
			0x101010ef10103838,
			[]KingChecker{
				{
					Knight,
					Bitboard(0x40000000000000),
				},
				{
					Rook,
					Bitboard(0x1000000000),
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			position, _ := tc.positionFen.Parse()
			fmt.Println(position)
			actualKDSBb, actualKCs := position.calculateOurKingDangerSquares()
			assert.Equal(t, actualKDSBb, tc.expectedKDSBb)
			if len(tc.expectedKCs) != 0 {
				assert.Equal(t, len(actualKCs), len(tc.expectedKCs))
				assert.ElementsMatch(t, actualKCs, tc.expectedKCs)
			} else {
				assert.Len(t, actualKCs, 0)
			}
		})
	}

}
