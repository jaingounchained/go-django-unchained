package src

type Move uint16

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

type PromotionPieceType uint16

const (
	KnightPromotion PromotionPieceType = 0 << 12
	BishopPromotion PromotionPieceType = 1 << 12
	RookPromotion   PromotionPieceType = 2 << 12
	QueenPromotion  PromotionPieceType = 3 << 12
)

type MoveType uint16

const (
	Normal    MoveType = 0 << 14
	Promotion MoveType = 1 << 14
	CastlingM MoveType = 2 << 14
	EnPassant MoveType = 3 << 14
)

// func (sq Square) String() string {
// 	if sq >= 0 && sq <= 63 {
// 		return string((sq%8)+97) + string((sq/8)+49)
// 	} else {
// 		return "No square"
// 	}
// }
