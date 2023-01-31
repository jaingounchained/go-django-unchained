package src

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

type Bitboard uint64

type PieceBitboard [TotalPieceTypes]Bitboard

type PiecePlacementR map[Color]*PieceBitboard

type OccupiedSquares map[Color]Bitboard

func (bb Bitboard) leftmostSignificantSquare() Square {
	return Square(Log2n(uint64(bb) & uint64(-bb)))
}

func (bb Bitboard) removePieces(ourSquares Bitboard) Bitboard {
	return bb & ^ourSquares
}

func (bb Bitboard) spawnMoves(initialSq Square) (moveList []Move) {
	for ; bb > 0; bb &= bb - 1 {
		finalSq := bb.leftmostSignificantSquare()
		moveList = append(moveList, Move(initialSq+finalSq<<6))
	}
	return
}

func (bb Bitboard) reverseBits() Bitboard {
	return Bitboard(bits.Reverse64(uint64(bb)))
}

func (bb Bitboard) isBitSet(sq Square) bool {
	nn := bb >> (sq - 1)
	return (nn & 1) == 1
}

type Direction int

const (
	north     Direction = 8
	south     Direction = -8
	west      Direction = -1
	east      Direction = 1
	northWest Direction = 7
	southWest Direction = -9
	northEast Direction = 9
	southEast Direction = -7
)

const (
	Rank1 int = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

const (
	filea int = iota
	fileb
	filec
	filed
	filee
	filef
	fileg
	fileh
)

func (bb Bitboard) shift(d Direction) Bitboard {
	switch d {
	case west, northWest, southWest:
		return bb.bitShift(d) & notHFile
	case east, northEast, southEast:
		return bb.bitShift(d) & notAFile
	default:
		return bb.bitShift(d)
	}
}

func (bb Bitboard) bitShift(i Direction) Bitboard {
	if i >= 0 {
		return bb << i
	} else {
		return bb >> -i
	}
}

// Structure for holding masks w.r.t. squares
type SqMask struct {
	bitMask            Bitboard
	rankMaskEx         Bitboard
	fileMaskEx         Bitboard
	diagonalMaskEx     Bitboard
	antiDiagonalMaskEx Bitboard
}

var notAFile Bitboard = 0xfefefefefefefefe
var notHFile Bitboard = 0x7f7f7f7f7f7f7f7f
var notABFile Bitboard = 0xFCFCFCFCFCFCFCFC
var notGHFile Bitboard = 0x3F3F3F3F3F3F3F3F

var sqMask [64]SqMask

// precalculate mask array indexed on square
func GenerateSquareMasks() {

	var north Bitboard = 0x0101010101010100
	var south Bitboard = 0x0080808080808080
	var east Bitboard = 0xFE
	var west Bitboard = 0x7F
	var northEast Bitboard = 0x8040201008040200
	var southWest Bitboard = 0x40201008040201
	var northWest Bitboard = 0x102040810204000
	var southEast Bitboard = 0x2040810204080

	// bit Mask
	for i := 0; i < 64; i++ {
		sqMask[i].bitMask = 1 << i
	}

	// file Mask excluding the square
	for i := 0; i < 64; i++ {
		sqMask[i].fileMaskEx = (north << i) | (south >> (63 - i))
	}

	// rank Mask excluding the square
	for i := 0; i < 8; i, east, west = i+1, east.eastOne(), west.westOne() {
		ea := east
		we := west
		for k := 0; k < 8*8; k, ea, we = k+8, ea<<8, we<<8 {
			sqMask[k+i].rankMaskEx |= ea   // east direction
			sqMask[k+7-i].rankMaskEx |= we // west direction
		}
	}

	// diagonal mask excluding the square
	for i := 0; i < 8; i, northEast, southWest = i+1, northEast.eastOne(), southWest.westOne() {
		noea := northEast
		sowe := southWest
		for k := 0; k < 8*8; k, noea, sowe = k+8, noea<<8, sowe>>8 {
			sqMask[k+i].diagonalMaskEx |= noea // east direction
			sqMask[56-k+7-i].diagonalMaskEx |= sowe
		}
	}

	// antidiagonal mask excluding the square
	for i := 7; i >= 0; i, northWest, southEast = i-1, northWest.westOne(), southEast.eastOne() {
		nowe := northWest
		soea := southEast
		for k := 0; k < 8*8; k, nowe, soea = k+8, nowe<<8, soea>>8 {
			sqMask[k+i].antiDiagonalMaskEx |= nowe
			sqMask[56-k+7-i].antiDiagonalMaskEx |= soea
		}
	}
}

func (bb Bitboard) eastOne() Bitboard {
	return (bb << 1) & notAFile
}

func (bb Bitboard) westOne() Bitboard {
	return (bb >> 1) & notHFile
}

func KnightAttacks(sq Square) Bitboard {
	var NB Bitboard = 1 << sq
	NAS :=
		(NB.bitShift(-10) & notGHFile) |
			(NB.bitShift(-17) & notHFile) |
			(NB.bitShift(-15) & notAFile) |
			(NB.bitShift(-6) & notABFile) |
			(NB.bitShift(10) & notABFile) |
			(NB.bitShift(17) & notAFile) |
			(NB.bitShift(15) & notHFile) |
			(NB.bitShift(6) & notGHFile)
	return NAS
}

func KingAttacks(sq Square) Bitboard {
	var KB Bitboard = 1 << sq
	KAS :=
		KB.shift(south) | KB.shift(north) |
			((KB.bitShift(-9) | KB.bitShift(7) | KB.bitShift(-1)) & notHFile) |
			((KB.bitShift(-7) | KB.bitShift(1) | KB.bitShift(9)) & notAFile)

	return KAS
}

/*
	Calculating
	 - file attacks
	 - rank attacks
	 - diagonal attacks
	 - antidiagonal attacks
	using Hyperbola Quintessence
	https://www.chessprogramming.org/Hyperbola_Quintessence
*/
func (occ Bitboard) fileAttacks(sq Square) Bitboard {
	var forward, reverse Bitboard

	// (o-r): masking the file & subtracting the sqaure
	forward = occ & sqMask[sq].fileMaskEx
	reverse = forward.reverseBits()
	// (o-2r)
	forward -= sqMask[sq].bitMask
	reverse -= sqMask[sq].bitMask.reverseBits()
	// (o-2r)^rev(o'-2r')
	forward ^= reverse.reverseBits()
	forward &= sqMask[sq].fileMaskEx

	return forward
}

func (occ Bitboard) rankAttacks(sq Square) Bitboard {
	var forward, reverse Bitboard

	// (o-r): masking the file & subtracting the sqaure
	forward = occ & sqMask[sq].rankMaskEx
	reverse = forward.reverseBits()
	// (o-2r)
	forward -= sqMask[sq].bitMask
	reverse -= sqMask[sq].bitMask.reverseBits()
	// (o-2r)^rev(o'-2r')
	forward ^= reverse.reverseBits()
	forward &= sqMask[sq].rankMaskEx

	return forward
}

func (occ Bitboard) diagonalAttacks(sq Square) Bitboard {
	var forward, reverse Bitboard

	// (o-r): masking the file & subtracting the sqaure
	forward = occ & sqMask[sq].diagonalMaskEx
	reverse = forward.reverseBits()
	// (o-2r)
	forward -= sqMask[sq].bitMask
	reverse -= sqMask[sq].bitMask.reverseBits()
	// (o-2r)^rev(o'-2r')
	forward ^= reverse.reverseBits()
	forward &= sqMask[sq].diagonalMaskEx

	return forward
}

func (occ Bitboard) antiDiagonalAttacks(sq Square) Bitboard {
	var forward, reverse Bitboard

	// (o-r): masking the file & subtracting the sqaure
	forward = occ & sqMask[sq].antiDiagonalMaskEx
	reverse = forward.reverseBits()
	// (o-2r)
	forward -= sqMask[sq].bitMask
	reverse -= sqMask[sq].bitMask.reverseBits()
	// (o-2r)^rev(o'-2r')
	forward ^= reverse.reverseBits()
	forward &= sqMask[sq].antiDiagonalMaskEx

	return forward
}

func Log2n(n uint64) uint16 {
	if n > 1 {
		return 1 + Log2n(n/2)
	} else {
		return 0
	}
}

func (pp PiecePlacementR) String() string {
	boardRep := make([][]string, 0)

	boardRep = append(boardRep,
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[Black][Pawn]), ""), Pawn.PieceRep(Black)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[Black][Knight]), ""), Knight.PieceRep(Black)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[Black][Bishop]), ""), Bishop.PieceRep(Black)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[Black][Rook]), ""), Rook.PieceRep(Black)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[Black][Queen]), ""), Queen.PieceRep(Black)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[Black][King]), ""), King.PieceRep(Black)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[White][Pawn]), ""), Pawn.PieceRep(White)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[White][Knight]), ""), Knight.PieceRep(White)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[White][Bishop]), ""), Bishop.PieceRep(White)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[White][Rook]), ""), Rook.PieceRep(White)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[White][Queen]), ""), Queen.PieceRep(White)),
		replacePiece(strings.Split(fmt.Sprintf("%64b", pp[White][King]), ""), King.PieceRep(White)),
	)

	boardString := generateBoardString(boardRep)

	boardStringRep := "\nPiece placement: \n"
	for i := 0; i < 8; i++ {
		boardStringRep += "\t -————-————-————-————-————-————-————-————-\n\t"
		for j := 7; j >= 0; j-- {
			index := i*8 + j
			boardStringRep += " | " + boardString[index] + " "
		}
		boardStringRep += " | " + strconv.Itoa(8-i) + "\n"
	}
	boardStringRep += "\t -————-————-————-————-————-————-————-————-\n"
	boardStringRep += "\t    a    b    c    d    e    f    g    h  \n"

	return boardStringRep
}

func (bb Bitboard) String() string {
	bitboardString := ""
	bitboardString += strconv.FormatUint(uint64(bb), 2)
	bitboardString = fmt.Sprintf("%064s", bitboardString)
	bitboardString = strings.Replace(bitboardString, "0", ".", -1)

	finalBitboardString := ""
	for i := 0; i < 8; i++ {
		for j := 7; j >= 0; j-- {
			index := i*8 + j
			finalBitboardString += bitboardString[index:index+1] + " "
		}
		finalBitboardString += "\n"
	}

	return finalBitboardString
}
