package src

import (
	"strconv"
)

type Color bool

const (
	White Color = true
	Black Color = false
)

func (c Color) String() string {
	switch c {
	case true:
		return "White"
	case false:
		return "Black"
	default:
		return "No Color"
	}
}

type CastlingType struct {
	kingSide  bool
	queenSide bool
}

type CastlingRights map[Color]CastlingType

func (cr CastlingRights) String() string {
	crRep := "\n"
	for k, v := range cr {
		if v.kingSide {
			crRep += "  - " + k.String() + " King Side Castling available\n"
		}
		if v.queenSide {
			crRep += "  - " + k.String() + " Queen Side Castling available\n"
		}
	}
	return crRep
}

type EnPassantTarget Square

type Position struct {
	// essential information
	piecePlacement  PiecePlacement
	activeColor     Color
	castlingRights  CastlingRights
	enPassantTarget Square
	halfMoveClock   uint16
	fullMoveNumber  uint16

	// auxiliary information
	occupiedSquaresColorWise OccupiedSquaresColorWise
	allOccupiedSquares       Bitboard
	// mailbox                  [64]PieceType

	// auxiliary information for checks and legal moves
	ourKingDangerSquares Bitboard
	kingCheckers         []KingChecker //  if len > 0, king is in check
	captureMask          Bitboard
	pushMask             Bitboard

	// Pinned pieces
}

type KingChecker struct {
	pieceType PieceType
	bitboard  Bitboard
}

func (position Position) generateAuxiliaryInfo() Position {
	// calculating all occupied squares color wise
	// updatedPosition := position
	occupiedSqauresColorWise := make(OccupiedSquaresColorWise)

	for color, pieceBitboard := range position.piecePlacement {
		for i := 0; i < len(pieceBitboard); i++ {
			occupiedSqauresColorWise[color] += pieceBitboard[i]
		}
	}

	position.occupiedSquaresColorWise = occupiedSqauresColorWise

	// calculating all occupied squares by all pieces
	position.allOccupiedSquares = occupiedSqauresColorWise[Black] | occupiedSqauresColorWise[White]

	// calculating our king danger squares & king checkers
	position.ourKingDangerSquares, position.kingCheckers = position.calculateOurKingDangerSquares()

	if len(position.kingCheckers) == 0 {
		position.captureMask, position.pushMask = Bitboard(0xFFFFFFFFFFFFFFFF), Bitboard(0xFFFFFFFFFFFFFFFF)
	} else if len(position.kingCheckers) == 1 {
		position.captureMask, position.pushMask = position.calculateCapturePushMask()
	} else {
		position.captureMask, position.pushMask = Bitboard(0), Bitboard(0)
	}

	return position
}

func (position Position) calculateOurKingDangerSquares() (Bitboard, []KingChecker) {
	// calculating our king danger squares
	/*
		KBb = King Bitboard
		usOS = our occupied squares
		oppOS = opponent occupied squares
		usOCMK = our occupied sqaures without king
		oppPP = opponent Piece Placement
		usKDSBK = our King Danger Squares by opponent King
		kCBQ = king checker by Queen
	*/
	us, opp := position.activeColor, !position.activeColor
	KBb := position.piecePlacement[us][King]
	usOS, oppOS := position.occupiedSquaresColorWise[us], position.occupiedSquaresColorWise[opp]
	usOSMK := usOS - KBb
	oppPP := position.piecePlacement[opp]

	usKDS := Bitboard(0)
	kCs := []KingChecker{}

	for p := Pawn; p < TotalPieceTypes; p++ {
		usKDSBP, kC := p.attackBbMP(opp, oppPP[p], KBb, usOSMK, oppOS)
		usKDS |= usKDSBP
		if kC.bitboard != 0 {
			kCs = append(kCs, kC)
		}
	}

	return usKDS, kCs
}

// Calculate attacking squares for every piece based on Bitboard
// attachBBMP = Bitboard of attacking squares by multiple piece of same type
func (pt PieceType) attackBbMP(color Color, pBB, kBb, usOS, oppOC Bitboard) (Bitboard, KingChecker) {
	attackBb := Bitboard(0)
	kingChecker := KingChecker{}

	for ; pBB > 0; pBB &= pBB - 1 {
		// position of piece: 0-63
		sq := pBB.leftmostSignificantSquare()

		// Calculating attacking squares of single piece including our own pieces
		attackBb |= pt.attackBbSP(color, sq, usOS, oppOC)
		if attackBb&kBb != 0 {
			kingChecker.pieceType = pt
			kingChecker.bitboard = Bitboard(1 << sq)
		}
	}

	return attackBb, kingChecker
}

func (position Position) calculateCapturePushMask() (Bitboard, Bitboard) {
	// kBb := position.piecePlacement[position.activeColor][King]
	checkingPiece := position.kingCheckers[0].pieceType

	captureMask := position.kingCheckers[0].bitboard

	pushMask := Bitboard(0)
	switch checkingPiece {
	case Rook:

	}
	return captureMask, pushMask
}

// https://en.wikipedia.org/wiki/Pin_(chess)#Absolute_pin
func (position Position) calculateAbsolutePinnedPieces() {

}

// utility toString functions
func (position Position) String() string {
	positionRep := ""
	positionRep += position.piecePlacement.String()
	positionRep += "\nActive color: " + position.activeColor.String()
	positionRep += "\nCastlingRights: " + position.castlingRights.String()
	// positionRep += "\nEn Passant target: " + position.enPassantTarget.String()
	positionRep += "\nHalf Move Clock: " + strconv.FormatInt(int64(position.halfMoveClock), 10)
	positionRep += "\nFull Move Number: " + strconv.FormatInt(int64(position.fullMoveNumber), 10) + "\n"

	return positionRep
}

func replacePiece(s []string, rep string) []string {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case "1":
			s[i] = rep
		default:
			s[i] = ""
		}
	}
	return s
}

func generateBoardString(s [][]string) (a [64]string) {
	for i := 0; i < 64; i++ {
		for j := 0; j < len(s); j++ {
			a[i] += s[j][i]
		}
		if a[i] == "" {
			a[i] = " "
		}
	}
	return
}

// // assuming legal moves
// func (oldPosition Position) makeMove(move Move) (Position, error) {
// 	us, them := oldPosition.activeColor, !oldPosition.activeColor
// 	var ourPieces = oldPosition.piecePlacement[us]
// 	var theirPieces = oldPosition.piecePlacement[them]

// 	var _ Square = Square(move & 63)
// 	var moveTo Square = Square((move >> 6) & 63)

// 	var moveFromBB Bitboard = 1 << (move & 63)
// 	var moveToBB Bitboard = 1 << ((move >> 6) & 63)

// 	moveType := move & moveTypeMask

// 	var updatedPiecePlacement PiecePlacementR = oldPosition.piecePlacement
// 	var updatedOurPieces = oldPosition.piecePlacement[us]
// 	var updatedTheirPieces = oldPosition.piecePlacement[them]

// 	var updatedCastlingRights = oldPosition.castlingRights

// 	var updatedEnPassantTarget Square

// 	switch moveType {
// 	case Move(Normal):
// 		// determining which piece is moved
// 		var movedPiece PieceType
// 		if moveFromBB&ourPieces[Pawn] != 0 {
// 			movedPiece = Pawn
// 		} else if moveFromBB&ourPieces[Knight] != 0 {
// 			movedPiece = Knight
// 		} else if moveFromBB&ourPieces[Bishop] != 0 {
// 			movedPiece = Bishop
// 		} else if moveFromBB&ourPieces[Rook] != 0 {
// 			movedPiece = Rook
// 		} else if moveFromBB&ourPieces[Queen] != 0 {
// 			movedPiece = Queen
// 		} else if moveFromBB&ourPieces[King] != 0 {
// 			movedPiece = King
// 		}

// 		// reset castling rights in case rook or king has moved from home square

// 		// set enpassant target in case pawn has moved 2 squares

// 		// changing position of the moved piece
// 		updatedOurPieces[movedPiece] = ourPieces[movedPiece] - moveFromBB + moveToBB
// 		updatedPiecePlacement[us] = updatedOurPieces
// 		// removing opponent piece in case it is present
// 		updatedTheirPieces[Pawn] &= ^(moveToBB)
// 		updatedTheirPieces[Knight] &= ^(moveToBB)
// 		updatedTheirPieces[Bishop] &= ^(moveToBB)
// 		updatedTheirPieces[Rook] &= ^(moveToBB)
// 		updatedTheirPieces[Queen] &= ^(moveToBB)
// 		updatedPiecePlacement[them] = updatedTheirPieces

// 	case Move(Promotion):
// 		var promotionType Move = move & promotionPieceTypeMask

// 		switch promotionType {
// 		case Move(QueenPromotion):
// 			updatedOurPieces[Queen] = ourPieces[Queen] + moveToBB
// 		case Move(RookPromotion):
// 			updatedOurPieces[Rook] = ourPieces[Rook] + moveToBB
// 		case Move(BishopPromotion):
// 			updatedOurPieces[Bishop] = ourPieces[Bishop] + moveToBB
// 		case Move(KnightPromotion):
// 			updatedOurPieces[Knight] = ourPieces[Knight] + moveToBB
// 		}
// 		updatedOurPieces[Pawn] = ourPieces[Pawn] - moveFromBB
// 		updatedPiecePlacement[us] = updatedOurPieces

// 		updatedTheirPieces[Knight] &= ^(moveToBB)
// 		updatedTheirPieces[Bishop] &= ^(moveToBB)
// 		updatedTheirPieces[Rook] &= ^(moveToBB)
// 		updatedTheirPieces[Queen] &= ^(moveToBB)

// 		updatedPiecePlacement[them] = updatedTheirPieces
// 	case Move(CastlingM):
// 		updatedOurPieces[King] = moveToBB
// 		switch moveTo {
// 		case g1:
// 			updatedOurPieces[Rook] -= (1 << h1) + (1 << f1)
// 			updatedCastlingRights[White] = CastlingType{}
// 		case c1:
// 			updatedOurPieces[Rook] -= (1 << a1) + (1 << c1)
// 			updatedCastlingRights[White] = CastlingType{}
// 		case g8:
// 			updatedOurPieces[Rook] -= (1 << h8) + (1 << f8)
// 			updatedCastlingRights[Black] = CastlingType{}
// 		case c8:
// 			updatedOurPieces[Rook] -= (1 << a8) + (1 << c8)
// 			updatedCastlingRights[Black] = CastlingType{}
// 		}
// 		updatedPiecePlacement[us] = updatedOurPieces
// 	case Move(EnPassant):
// 		updatedOurPieces[Pawn] = ourPieces[Pawn] - moveFromBB + moveToBB
// 		if oldPosition.activeColor == White {
// 			updatedTheirPieces[Pawn] = theirPieces[Pawn] - (moveFromBB >> 8)
// 		} else {
// 			updatedTheirPieces[Pawn] = theirPieces[Pawn] - (moveFromBB << 8)
// 		}
// 		updatedEnPassantTarget = 64
// 	}

// 	// only update full move number in case black has made a move
// 	var updatedFullMoveNumber uint16 = oldPosition.fullMoveNumber
// 	if oldPosition.activeColor == Black {
// 		updatedFullMoveNumber += 1
// 	}

// 	// update auxiliary information & return
// 	return Position{
// 		piecePlacement:  updatedPiecePlacement,
// 		activeColor:     !oldPosition.activeColor, // reverse the color White <-> Black
// 		castlingRights:  updatedCastlingRights,
// 		enPassantTarget: updatedEnPassantTarget,
// 		halfMoveClock:   oldPosition.halfMoveClock + 1, // update halfMoveClock since a move is made
// 		fullMoveNumber:  updatedFullMoveNumber,
// 	}.updateAuxiliaryInfo(), nil
// }

// func (position Position) updateAuxiliaryInfo() Position {
// 	position.occupiedSquaresColorWise[White] = position.piecePlacement[White][Pawn] | position.piecePlacement[White][Knight] | position.piecePlacement[White][Bishop] | position.piecePlacement[White][Rook] | position.piecePlacement[White][Queen] | position.piecePlacement[White][King]
// 	position.occupiedSquaresColorWise[Black] = position.piecePlacement[Black][Pawn] | position.piecePlacement[Black][Knight] | position.piecePlacement[Black][Bishop] | position.piecePlacement[Black][Rook] | position.piecePlacement[Black][Queen] | position.piecePlacement[Black][King]

// 	return position
// }
