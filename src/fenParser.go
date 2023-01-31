package src

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Fen string

func (fen Fen) Parse() (Position, error) {
	fenSplit := strings.Split(string(fen), " ")

	// fen string should contain 6 components
	if len(fenSplit) != 6 {
		return Position{}, errors.New("number of components in fen != 6")
	}

	piecePlacement, err := parsePiecePlacement(fenSplit[0])
	if err != nil {
		return Position{}, err
	}

	activeColor, err := parseActiveColor(fenSplit[1])
	if err != nil {
		return Position{}, err
	}

	castlingRights, err := parseCastlingRights(fenSplit[2])
	if err != nil {
		return Position{}, err
	}

	enPassantTarget, err := parseEnPassantTarget(fenSplit[3])
	if err != nil {
		return Position{}, err
	}

	halfMoveClock, err := strconv.Atoi(fenSplit[4])
	if err != nil {
		return Position{}, err
	}

	fullMoveNumber, err := strconv.Atoi(fenSplit[5])
	if err != nil {
		return Position{}, err
	}

	return Position{
		piecePlacement:  piecePlacement,
		activeColor:     activeColor,
		castlingRights:  castlingRights,
		enPassantTarget: enPassantTarget,
		halfMoveClock:   uint16(halfMoveClock),
		fullMoveNumber:  uint16(fullMoveNumber),

		occupiedSquares: generateOccupiedSquares(piecePlacement),
	}, nil
}

func parsePiecePlacement(pp string) (PiecePlacementR, error) {
	ranks := strings.Split(pp, "/")
	// return error if number of ranks != 8
	if len(ranks) != 8 {
		return nil, errors.New("number of ranks in fen !=8")
	}

	var whitePieces, blackPieces PieceBitboard
	for i := 0; i < 8; i++ {
		k := 0
		for j := 0; j < len(ranks[i]); j++ {
			index := 8*(7-i) + k
			switch ranks[i][j] {
			case 'P':
				whitePieces[Pawn] += 1 << index
				k++
			case 'N':
				whitePieces[Knight] += 1 << index
				k++
			case 'B':
				whitePieces[Bishop] += 1 << index
				k++
			case 'R':
				whitePieces[Rook] += 1 << index
				k++
			case 'Q':
				whitePieces[Queen] += 1 << index
				k++
			case 'K':
				whitePieces[King] += 1 << index
				k++
			case 'p':
				blackPieces[Pawn] += 1 << index
				k++
			case 'n':
				blackPieces[Knight] += 1 << index
				k++
			case 'b':
				blackPieces[Bishop] += 1 << index
				k++
			case 'r':
				blackPieces[Rook] += 1 << index
				k++
			case 'q':
				blackPieces[Queen] += 1 << index
				k++
			case 'k':
				blackPieces[King] += 1 << index
				k++
			default:
				emptyPositions, err := strconv.ParseInt(string(ranks[i][j]), 10, 32)
				if emptyPositions < 1 || emptyPositions > 8 || err != nil {
					return nil, fmt.Errorf("char %s: invalid empty space", string(ranks[i][j]))
				}
				k += int(emptyPositions)
			}
		}
	}

	// Piece Placement map
	return PiecePlacementR{
		White: &whitePieces,
		Black: &blackPieces,
	}, nil
}

func parseActiveColor(ac string) (Color, error) {
	switch ac {
	case "w":
		return White, nil
	case "b":
		return Black, nil
	default:
		return false, errors.New("active color invalid")
	}
}

func parseCastlingRights(cr string) (CastlingRights, error) {
	blackCastlingType, WhiteCastlingType := CastlingType{}, CastlingType{}

	for i := 0; i < len(cr); i++ {
		switch cr[i] {
		case 'K':
			WhiteCastlingType.kingSide = true
		case 'Q':
			WhiteCastlingType.queenSide = true
		case 'k':
			blackCastlingType.kingSide = true
		case 'q':
			blackCastlingType.queenSide = true
		case '-':
			return make(CastlingRights), nil
		default:
			return nil, errors.New("castling rights invalid format")
		}
	}
	return CastlingRights{
		White: &WhiteCastlingType,
		Black: &blackCastlingType,
	}, nil
}

func parseEnPassantTarget(ept string) (Square, error) {
	enPassantTargetSquare := rune(64)
	chars := []rune(ept)

	// checking if move is in right format
	if len(chars) == 2 && chars[0] <= 104 && chars[0] >= 97 && chars[1] <= 56 && chars[1] >= 49 {
		file := chars[0] - 97
		rank := chars[1] - 49
		enPassantTargetSquare = rank*8 + file
		return Square(enPassantTargetSquare), nil
	} else if ept == "-" {
		return 64, nil
	} else {
		return 0, errors.New("en passant target in wrong format")
	}
}

func generateOccupiedSquares(pp PiecePlacementR) OccupiedSquares {
	os := make(OccupiedSquares)
	os[White] = 0
	os[Black] = 0

	for color, pieceBitboard := range pp {
		for i := 0; i < len(pieceBitboard); i++ {
			os[color] += pieceBitboard[i]
		}
	}
	return os
}
