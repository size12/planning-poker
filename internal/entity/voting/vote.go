package voting

type Vote string

const (
	Zero       Vote = "ZERO"
	OneHalf    Vote = "ONE_HALF"
	One        Vote = "ONE"
	Two        Vote = "TWO"
	Three      Vote = "THREE"
	Five       Vote = "FIVE"
	Eight      Vote = "EIGHT"
	Thirteen   Vote = "THIRTEEN"
	TwentyOne  Vote = "TWENTY_ONE"
	ThirtyFour Vote = "THIRTY_FOUR"
	Question   Vote = "QUESTION"
	Coffee     Vote = "COFFEE"
)
