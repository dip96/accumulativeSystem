package balance

//TODO может подобавлять интерфейсов и вызывать через них структуру?
// https://www.youtube.com/watch?v=vR-UVn-5AOs&list=PLP19RjSHH4aE9pB77yT1PbXzftGsXFiGl&index=9

type UserBalance struct {
	ID               int64
	UserID           int
	Balance          float64
	WithdrawnBalance float64
}
