package operations

type Strategy int

const (
	StrategyNone Strategy = iota
	StrategyTab
	StrategyWindIn
	StrategyWindOut
	StrategyLawnMower
)

type Type int

const (
	TypeNone Type = iota
	TypeDrill
	TypeRoute
)

type Operation struct {
	StrategyNormal Strategy
	StrategySmall  Strategy
	Type           Type
	GrainMatched   bool
}
