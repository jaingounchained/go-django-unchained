package src

type PieceType int

const (
	_ PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
	TotalPieceTypes
)

func (pt PieceType) isSlider() bool {
	switch pt {
	case Pawn:
		return false
	case Knight:
		return false
	case Bishop:
		return true
	case Rook:
		return true
	case Queen:
		return true
	case King:
		return false
	default:
		return false
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
