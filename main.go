package main

import (
	"fmt"

	"github.com/bhavya5jain/go-django-unchained/src"
)

func main() {

	// src.UpdatePositionTest_Promotion()

	// fenPosition1 := src.Fen("8/8/8/8/8/8/PPP2PPP/R6R w KQkq - 0 1")
	// randomPosition := src.Fen("2r5/2r2ppp/5k2/B1nBb3/2R5/6P1/P4P1P/3R2K1 b - - 2 35")
	// // firstMove := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"

	// random := src.FenParser(randomPosition)

	// random.PrintChessPosition()

	// moveList := src.GenerateMoves(random, 1)

	// fmt.Println(moveList)
	// // src.FenParser(randomPosition).PrintChessPosition()
	// src.FenParser(fenPosition1).PrintChessPosition()
	// src.FenParser(firstMove).PrintChessPosition()

	// src.PrintAttackingSquares(0x2040810204080)

	// src.GenerateSquareMasks()

	// fmt.Println("Piecetype: ", src.Pawn)

	// position, err := randomPosition.Parse()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(position)

	// for k, v := range position {
	// 	fmt.Printf("%v: \n%v\n", k, v)
	// }

	// var move src.Move = 3<<14 + 3

	// // moveType, _ := src.UpdatePosition(move)

	// fmt.Println(moveType)

	fmt.Println(src.Bitboard(0x1000000000))
}
