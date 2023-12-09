package src

func GenerateAllMoves(position Position, depth int) (moveList []Move) {

	if len(position.kingCheckers) == 2 {
		// double check
		// only generate king moves
		return
	}

	// if position.isCheck {
	// 	// 1. generate king moves
	// 	return
	// }

	return moveList
}

// https://peterellisjones.com/posts/generating-legal-chess-moves-efficiently/
func (pt PieceType) generateMoves(position Position) (moveList []Move) {
	// capture mask and push mask
	cM, pM := position.captureMask, position.pushMask

	us, opponent := position.activeColor, !position.activeColor
	ourOccupiedSquares, opponentOccupiedSquares := position.occupiedSquaresColorWise[us], position.occupiedSquaresColorWise[opponent]
	pieceBitboard := position.piecePlacement[us][pt]

	for ; pieceBitboard > 0; pieceBitboard &= pieceBitboard - 1 {
		// position of piece: 0-63
		sq := pieceBitboard.leftmostSignificantSquare()

		// Calculating all attacking squares including our own pieces
		attackingSquares := pt.attackBbSP(us, sq, ourOccupiedSquares, opponentOccupiedSquares)

		// Calculating possible squares where piece can attack by removing our pieces
		possibleSquares := attackingSquares.removePieces(ourOccupiedSquares)

		// Calculating possible squares where piece can capture
		captureSquares := possibleSquares & opponentOccupiedSquares & cM

		moveList = append(moveList, captureSquares.spawnMoves(sq, Capture)...)

		// Calculating possible squares where piece can jump
		jumpableSquares := (possibleSquares - captureSquares) & pM

		moveList = append(moveList, jumpableSquares.spawnMoves(sq, Normal)...)
	}
	return
}

func generatePawnMoves(position Position) (moveList []Move) {

	us, opponent := position.activeColor, !position.activeColor

	var up, upRight, upLeft Direction
	var a, b, c, d Square

	if us == White {
		up = north
		upRight = northEast
		upLeft = northWest
		a, b = 8, 47 // rank 2-6
		c, d = 8, 16 // rank 2
	} else {
		up = south
		upRight = southWest
		upLeft = southEast
		a, b = 16, 55 // rank 7,3
		c, d = 48, 56 // rank 7
	}

	_, opponentSquares := position.occupiedSquaresColorWise[us], position.occupiedSquaresColorWise[opponent]
	pieceBitboard := position.piecePlacement[us][Pawn]
	occ := position.allOccupiedSquares
	ep := position.enPassantTarget

	for ; pieceBitboard > 0; pieceBitboard &= pieceBitboard - 1 {
		// position of piece: 0-63
		sq := pieceBitboard.leftmostSignificantSquare()

		if sq >= a && sq <= b {
			// Possible squares where pawn can jump 1 square forward
			jumpableSquares := Bitboard(1 << sq).shift(up)
			jumpableSquares = jumpableSquares.removePieces(occ)

			moveList = append(moveList, jumpableSquares.spawnMoves(sq, Normal)...)

			// Possible squares where pawn can jump 2 square forward from starting square
			if jumpableSquares != 0 && sq >= c && sq <= d {
				jumpable2Squares := Bitboard(1 << sq).shift(up).shift(up)
				jumpable2Squares = jumpable2Squares.removePieces(occ)
				moveList = append(moveList, jumpable2Squares.spawnMoves(sq, DoublePawnPush)...)
			}

			// Possible squares where pawn can attack
			attackingSquares := PawnAttacks[us][sq] & opponentSquares

			moveList = append(moveList, attackingSquares.spawnMoves(sq, Capture)...)

			// enpassant capture condition
			if ep < 64 && (Bitboard(1<<ep).shift(-upRight)|Bitboard(1<<ep).shift(-upLeft))&Bitboard(1<<sq) != 0 {
				moveList = append(moveList, EnPassant+Move(sq+Square(ep<<6)))
			}
		} else {
			// promotion to major piece

			// Possible squares where pawn can jump 1 square forward
			jumpableSquares := Bitboard(1 << sq).shift(up)
			jumpableSquares = jumpableSquares.removePieces(occ)

			if jumpableSquares != 0 {
				q := jumpableSquares.leftmostSignificantSquare()
				moveList = append(moveList, QueenPromotionNormal+Move(sq+q<<6))
				moveList = append(moveList, RookPromotionNormal+Move(sq+q<<6))
				moveList = append(moveList, BishopPromotionNormal+Move(sq+q<<6))
				moveList = append(moveList, KnightPromotionNormal+Move(sq+q<<6))
			}

			// Possible squares where pawn can attack
			attackingSquares := (Bitboard(1<<sq).shift(upRight) + Bitboard(1<<sq).shift(upLeft)) & opponentSquares

			for ; attackingSquares > 0; attackingSquares &= attackingSquares - 1 {
				q := attackingSquares.leftmostSignificantSquare()
				moveList = append(moveList, QueenPromotionCapture+Move(sq+q<<6))
				moveList = append(moveList, RookPromotionCapture+Move(sq+q<<6))
				moveList = append(moveList, BishopPromotionCapture+Move(sq+q<<6))
				moveList = append(moveList, KnightPromotionCapture+Move(sq+q<<6))
			}
		}
	}
	return moveList
}

func generateKingMoves(position Position) (moveList []Move) {
	us, opponent := position.activeColor, !position.activeColor
	ourOccupiedSquares, opponentOccupiedSquares := position.occupiedSquaresColorWise[us], position.occupiedSquaresColorWise[opponent]
	// usKDS = our King Danger Squares
	// usKDS := position.ourKingDangerSquares
	allOccupiedSquares := position.allOccupiedSquares
	bitboard := position.piecePlacement[us][King]
	castling := position.castlingRights[us]

	for ; bitboard > 0; bitboard &= bitboard - 1 {
		// position of piece: 0-63
		sq := bitboard.leftmostSignificantSquare()

		attackingSquares := KingAttacks[sq]

		// Possible squares where king can jump
		possibleSquares := attackingSquares.removePieces(ourOccupiedSquares)

		captureSquares := possibleSquares & opponentOccupiedSquares

		moveList = append(moveList, captureSquares.spawnMoves(sq, Capture)...)

		jumpableSquares := possibleSquares - captureSquares

		moveList = append(moveList, jumpableSquares.spawnMoves(sq, Normal)...)
	}

	// Bitboard for squares between king and respective rook
	var kksr, kqsr Bitboard
	var kingSideCastling, queenSideCastling Move
	if us == White {
		kksr, kqsr = 0x60, 0xE
		kingSideCastling, queenSideCastling = WhiteKingSideCastling, WhiteQueenSideCastling
	} else {
		kksr, kqsr = 0x6000000000000000, 0xE00000000000000
		kingSideCastling, queenSideCastling = BlackKingSideCastling, BlackQueenSideCastling
	}
	// castling
	// is castling available & are there any peices between king and rook
	if castling.kingSide && kksr&allOccupiedSquares == 0 {
		moveList = append(moveList, kingSideCastling)
	}
	if castling.queenSide && kqsr&allOccupiedSquares == 0 {
		moveList = append(moveList, queenSideCastling)
	}

	return moveList
}
