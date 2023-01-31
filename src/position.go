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

type CastlingRights map[Color]*CastlingType

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
	piecePlacement  PiecePlacementR
	activeColor     Color
	castlingRights  CastlingRights
	enPassantTarget Square
	halfMoveClock   uint16
	fullMoveNumber  uint16

	occupiedSquares OccupiedSquares
}

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
