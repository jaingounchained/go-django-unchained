package src

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

type Bitboard uint64

type PieceBitboard [TotalPieceTypes]Bitboard

type PiecePlacement map[Color]PieceBitboard

type OccupiedSquaresColorWise map[Color]Bitboard

func (bb Bitboard) leftmostSignificantSquare() Square {
	return Square(Log2n(uint64(bb) & uint64(-bb)))
}

func (bb Bitboard) removePieces(ourSquares Bitboard) Bitboard {
	return bb & ^ourSquares
}

func (bb Bitboard) spawnMoves(initialSq Square, moveType Move) (moveList []Move) {
	for ; bb > 0; bb &= bb - 1 {
		finalSq := bb.leftmostSignificantSquare()
		moveList = append(moveList, moveType+Move(initialSq+finalSq<<6))
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

// const (
// 	Rank1 int = iota
// 	Rank2
// 	Rank3
// 	Rank4
// 	Rank5
// 	Rank6
// 	Rank7
// 	Rank8
// )

// const (
// 	filea int = iota
// 	fileb
// 	filec
// 	filed
// 	filee
// 	filef
// 	fileg
// 	fileh
// )

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

type Slider int

const (
	file Slider = iota
	rank
	diagonal
	antiDiagonal
)

/*
	SqMask = Square Mask
	Structure for holding masks w.r.t. squares
*/
type SqMask struct {
	bitMask      Bitboard
	sliderMaskEx map[Slider]Bitboard
}

var notAFile Bitboard = 0xfefefefefefefefe
var notHFile Bitboard = 0x7f7f7f7f7f7f7f7f
var notABFile Bitboard = 0xFCFCFCFCFCFCFCFC
var notGHFile Bitboard = 0x3F3F3F3F3F3F3F3F

/*
	Calcuated in starting to generate sliding piece attacks
*/
var sqMask [64]SqMask

// precalculate mask array indexed on square
func GenerateSquareMasks() {
	var n Bitboard = 0x0101010101010100
	var s Bitboard = 0x0080808080808080
	var e Bitboard = 0xFE
	var w Bitboard = 0x7F
	var ne Bitboard = 0x8040201008040200
	var sw Bitboard = 0x40201008040201
	var nw Bitboard = 0x102040810204000
	var se Bitboard = 0x2040810204080

	// bit Mask
	for i := 0; i < 64; i++ {
		sqMask[i].bitMask = 1 << i
	}

	// file Mask excluding the square
	for i := 0; i < 64; i++ {
		sqMask[i].sliderMaskEx = make(map[Slider]Bitboard)
		sqMask[i].sliderMaskEx[file] = (n << i) | (s >> (63 - i))
	}

	// rank Mask excluding the square
	for i := 0; i < 8; i, e, w = i+1, e.shift(east), w.shift(west) {
		ea := e
		we := w
		for k := 0; k < 8*8; k, ea, we = k+8, ea<<8, we<<8 {
			sqMask[k+i].sliderMaskEx[rank] |= ea   // east direction
			sqMask[k+7-i].sliderMaskEx[rank] |= we // west direction
		}
	}

	// diagonal mask excluding the square
	for i := 0; i < 8; i, ne, sw = i+1, ne.shift(east), sw.shift(west) {
		noea := ne
		sowe := sw
		for k := 0; k < 8*8; k, noea, sowe = k+8, noea<<8, sowe>>8 {
			sqMask[k+i].sliderMaskEx[diagonal] |= noea // east direction
			sqMask[56-k+7-i].sliderMaskEx[diagonal] |= sowe
		}
	}

	// antidiagonal mask excluding the square
	for i := 7; i >= 0; i, nw, se = i-1, nw.shift(west), se.shift(east) {
		nowe := nw
		soea := se
		for k := 0; k < 8*8; k, nowe, soea = k+8, nowe<<8, soea>>8 {
			sqMask[k+i].sliderMaskEx[antiDiagonal] |= nowe
			sqMask[56-k+7-i].sliderMaskEx[antiDiagonal] |= soea
		}
	}
}

// Calculated at starting
var KingAttacks [64]Bitboard
var KnightAttacks [64]Bitboard
var PawnAttacks map[Color][64]Bitboard

func GenerateNonSlidingPieceTypeAttackingSquares() {
	var WhitePawnAttacks, BlackPawnAttacks [64]Bitboard
	PawnAttacks = make(map[Color][64]Bitboard)

	for i := 0; i < 64; i++ {
		// sqBb = square Bitboard
		sqBb := Bitboard(1 << i)

		// generating King attacking squares
		KingAttacks[i] = sqBb.shift(south) | sqBb.shift(north) |
			((sqBb.bitShift(-9) | sqBb.bitShift(7) | sqBb.bitShift(-1)) & notHFile) |
			((sqBb.bitShift(-7) | sqBb.bitShift(1) | sqBb.bitShift(9)) & notAFile)

		// generating Knight attacking squares
		KnightAttacks[i] = (sqBb.bitShift(-10) & notGHFile) |
			(sqBb.bitShift(-17) & notHFile) |
			(sqBb.bitShift(-15) & notAFile) |
			(sqBb.bitShift(-6) & notABFile) |
			(sqBb.bitShift(10) & notABFile) |
			(sqBb.bitShift(17) & notAFile) |
			(sqBb.bitShift(15) & notHFile) |
			(sqBb.bitShift(6) & notGHFile)

		// generating Pawn attacking squares for White & Black
		// White Pawn Attacks
		WhitePawnAttacks[i] = sqBb.shift(northEast) | sqBb.shift(northWest)
		// Black Pawn Attacks
		BlackPawnAttacks[i] = sqBb.shift(southEast) | sqBb.shift(southWest)
		continue
	}

	PawnAttacks[White], PawnAttacks[Black] = WhitePawnAttacks, BlackPawnAttacks
}

// Calculate attacking squares for every piece based on square
// attackBBSP = Bitboard of attacking squares by a single piece
func (pt PieceType) attackBbSP(color Color, sq Square, ourSquares, opponentSquares Bitboard) Bitboard {
	//	allOcc = all occupied squares
	allOcc := ourSquares | opponentSquares
	switch pt {
	case King:
		return KingAttacks[sq]
	case Queen:
		return diagonal.sliderAttacks(sq, allOcc) | antiDiagonal.sliderAttacks(sq, allOcc) | file.sliderAttacks(sq, allOcc) | rank.sliderAttacks(sq, allOcc)
	case Rook:
		return file.sliderAttacks(sq, allOcc) | rank.sliderAttacks(sq, allOcc)
	case Bishop:
		return diagonal.sliderAttacks(sq, allOcc) | antiDiagonal.sliderAttacks(sq, allOcc)
	case Knight:
		return KnightAttacks[sq]
	case Pawn:
		return PawnAttacks[color][sq]
	}
	return 0
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
// occ = occupied squares in Bitboars
func (slider Slider) sliderAttacks(sq Square, occ Bitboard) Bitboard {
	var forward, reverse Bitboard

	// (o-r): masking the file & subtracting the sqaure
	forward = occ & sqMask[sq].sliderMaskEx[slider]
	reverse = forward.reverseBits()
	// (o-2r)
	forward -= sqMask[sq].bitMask
	reverse -= sqMask[sq].bitMask.reverseBits()
	// (o-2r)^rev(o'-2r')
	forward ^= reverse.reverseBits()
	forward &= sqMask[sq].sliderMaskEx[slider]

	return forward
}

func Log2n(n uint64) uint16 {
	if n > 1 {
		return 1 + Log2n(n/2)
	} else {
		return 0
	}
}

/*
	Utility toString functions
*/
func (pp PiecePlacement) String() string {
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
