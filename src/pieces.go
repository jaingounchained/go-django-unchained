package src

import "fmt"

type PieceType int

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
	TotalPieceTypes
)

func (pt PieceType) generateMoves(position Position) (moveList []Move) {

	// only for generating moves for knights, bishops, rooks and queen
	if pt == Pawn || pt == King {
		return
	}

	us, opponent := position.activeColor, !position.activeColor
	ourSquares, opponentSquares := position.occupiedSquares[us], position.occupiedSquares[opponent]
	bitboard := position.piecePlacement[us][pt]

	for ; bitboard > 0; bitboard &= bitboard - 1 {
		// position of piece: 0-63
		sq := bitboard.leftmostSignificantSquare()

		attackingSquares := pt.attackingSquares(sq, ourSquares, opponentSquares)

		// Possible squares where piece can jump
		possibleSquares := attackingSquares.removePieces(ourSquares)

		moveList = append(moveList, possibleSquares.spawnMoves(sq)...)
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

	ourSquares, opponentSquares := position.occupiedSquares[us], position.occupiedSquares[opponent]
	bitboard := position.piecePlacement[us][Pawn]
	occ := ourSquares | opponentSquares
	ep := position.enPassantTarget

	for ; bitboard > 0; bitboard &= bitboard - 1 {
		// position of piece: 0-63
		sq := bitboard.leftmostSignificantSquare()

		if sq >= a && sq <= b {
			// Possible squares where pawn can jump 1 square forward
			possibleSquares := Bitboard(1 << sq).shift(up)
			possibleSquares = possibleSquares.removePieces(occ)

			// Possible squares where pawn can jump 2 square forward from starting square
			if possibleSquares != 0 && sq >= c && sq <= d {
				possibleSquares += Bitboard(1 << sq).shift(up).shift(up)
				possibleSquares = possibleSquares.removePieces(occ)
			}

			// Possible squares where pawn can attack
			possibleSquares += (Bitboard(1<<sq).shift(upRight) + Bitboard(1<<sq).shift(upLeft)) & opponentSquares

			fmt.Println(possibleSquares.String())

			moveList = append(moveList, possibleSquares.spawnMoves(sq)...)

			// enpassant capture condition

			if ep < 64 && (Bitboard(1<<ep).shift(-upRight)|Bitboard(1<<ep).shift(-upLeft))&Bitboard(1<<sq) != 0 {
				fmt.Println("Adding enpassant move")
				moveList = append(moveList, Move(Square(EnPassant)+sq+Square(ep<<6)))
			}
		} else {
			// promotion to major piece

			// Possible squares where pawn can jump 1 square forward
			possibleSquares := Bitboard(1 << sq).shift(up)
			possibleSquares = possibleSquares.removePieces(occ)

			// Possible squares where pawn can attack
			possibleSquares += (Bitboard(1<<sq).shift(upRight) + Bitboard(1<<sq).shift(upLeft)) & opponentSquares

			fmt.Println(possibleSquares.String())

			for ; possibleSquares > 0; possibleSquares &= possibleSquares - 1 {
				q := possibleSquares.leftmostSignificantSquare()
				moveList = append(moveList, Move(Square(Promotion)+Square(QueenPromotion)+sq+q<<6))
				moveList = append(moveList, Move(Square(Promotion)+Square(RookPromotion)+sq+q<<6))
				moveList = append(moveList, Move(Square(Promotion)+Square(BishopPromotion)+sq+q<<6))
				moveList = append(moveList, Move(Square(Promotion)+Square(KnightPromotion)+sq+q<<6))
			}
		}
	}
	return moveList
}

func generateKingMoves(position Position) (moveList []Move) {
	us, opponent := position.activeColor, !position.activeColor
	ourSquares, opponentSquares := position.occupiedSquares[us], position.occupiedSquares[opponent]
	occ := ourSquares | opponentSquares
	bitboard := position.piecePlacement[us][King]
	castling := position.castlingRights[us]

	for ; bitboard > 0; bitboard &= bitboard - 1 {
		// position of piece: 0-63
		sq := bitboard.leftmostSignificantSquare()

		attackingSquares := KingAttacks(sq)

		// Possible squares where king can jump
		possibleSquares := attackingSquares.removePieces(ourSquares)

		moveList = append(moveList, possibleSquares.spawnMoves(sq)...)

	}

	if castling != nil {
		// Bitboard for squares between king and respective rook
		var kksr, kqsr Bitboard
		var kingHomeSquare, kingSideCastlingSquare, queenSideCastlingSquare Square
		if us == White {
			kksr, kqsr = 0x60, 0xE
			kingHomeSquare, kingSideCastlingSquare, queenSideCastlingSquare = e1, g1, c1
		} else {
			kksr, kqsr = 0x6000000000000000, 0xE00000000000000
			kingHomeSquare, kingSideCastlingSquare, queenSideCastlingSquare = e8, g8, c8
		}
		// castling
		// is castling available & are there any peices between king and rook
		if castling.kingSide && kksr&occ == 0 {
			moveList = append(moveList, Move(kingHomeSquare+kingSideCastlingSquare<<6+Square(CastlingM)))
		}
		if castling.queenSide && kqsr&occ == 0 {
			moveList = append(moveList, Move(kingHomeSquare+queenSideCastlingSquare<<6+Square(CastlingM)))
		}
	}
	return moveList
}

func (pt PieceType) attackingSquares(sq Square, ourSquares Bitboard, opponentSquares Bitboard) Bitboard {
	allOcc := ourSquares | opponentSquares
	switch pt {
	case Knight:
		return KnightAttacks(sq)
	case Bishop:
		return allOcc.diagonalAttacks(sq) | allOcc.antiDiagonalAttacks(sq)
	case Rook:
		return allOcc.fileAttacks(sq) | allOcc.rankAttacks(sq)
	case Queen:
		return allOcc.diagonalAttacks(sq) | allOcc.antiDiagonalAttacks(sq) | allOcc.fileAttacks(sq) | allOcc.rankAttacks(sq)
	default:
		return 0
	}
}

func (pt PieceType) String() string {
	switch pt {
	case Pawn:
		return "Pawn"
	case Knight:
		return "Knight"
	case Bishop:
		return "Bishop"
	case Rook:
		return "Rook"
	case Queen:
		return "Queen"
	case King:
		return "King"
	case TotalPieceTypes:
		return "All Pieces"
	default:
		return "No Piece"
	}
}

func (pt PieceType) PieceRep(c Color) string {
	switch c {
	case White:
		switch pt {
		case Pawn:
			return "♟︎"
		case Knight:
			return "♞"
		case Bishop:
			return "♝"
		case Rook:
			return "♜"
		case Queen:
			return "♛"
		case King:
			return "♚"
		default:
			return " "
		}
	case Black:
		switch pt {
		case Pawn:
			return "♙"
		case Knight:
			return "♘"
		case Bishop:
			return "♗"
		case Rook:
			return "♖"
		case Queen:
			return "♕"
		case King:
			return "♔"
		default:
			return " "
		}
	default:
		return " "
	}
}
