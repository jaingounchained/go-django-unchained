package src

type Square uint16

const (
	a1 Square = iota
	b1
	c1
	d1
	e1
	f1
	g1
	h1
	a2
	b2
	c2
	d2
	e2
	f2
	g2
	h2
	a3
	b3
	c3
	d3
	e3
	f3
	g3
	h3
	a4
	b4
	c4
	d4
	e4
	f4
	g4
	h4
	a5
	b5
	c5
	d5
	e5
	f5
	g5
	h5
	a6
	b6
	c6
	d6
	e6
	f6
	g6
	h6
	a7
	b7
	c7
	d7
	e7
	f7
	g7
	h7
	a8
	b8
	c8
	d8
	e8
	f8
	g8
	h8
)

type Move uint16

// Move Types
const (
	// Ordinary moves
	Normal  Move = 0 << 12
	Capture Move = 1 << 12

	// Promotions
	KnightPromotionNormal  Move = 2 << 12
	KnightPromotionCapture Move = 3 << 12
	BishopPromotionNormal  Move = 4 << 12
	BishopPromotionCapture Move = 5 << 12
	RookPromotionNormal    Move = 6 << 12
	RookPromotionCapture   Move = 7 << 12
	QueenPromotionNormal   Move = 8 << 12
	QueenPromotionCapture  Move = 9 << 12

	// Double Pawn Push
	DoublePawnPush Move = 10 << 12

	// Castling
	WhiteKingSideCastling  Move = 11 << 12
	WhiteQueenSideCastling Move = 12 << 12
	BlackKingSideCastling  Move = 13 << 12
	BlackQueenSideCastling Move = 14 << 12

	// En Passant
	EnPassant Move = 15 << 12
)

// Promotion Piece Type of type Move
// const (
// 	KnightPromotion Move = 0 << 12
// 	BishopPromotion Move = 1 << 12
// 	RookPromotion   Move = 2 << 12
// 	QueenPromotion  Move = 3 << 12
// )
// const promotionPieceTypeMask = Move(3) << 12

// // Move Type of type Move
// const (
// 	Normal    Move = 0 << 14
// 	Promotion Move = 1 << 14
// 	CastlingM Move = 2 << 14
// 	EnPassant Move = 3 << 14
// )
// const moveTypeMask = Move(3) << 14

// func (sq Square) String() string {
// 	if sq >= 0 && sq <= 63 {
// 		return string((sq%8)+97) + string((sq/8)+49)
// 	} else {
// 		return "No square"
// 	}
// }

func squaresToMove(start, end Square, moveType Move) Move {
	return Move(start+end<<6) + moveType
}
