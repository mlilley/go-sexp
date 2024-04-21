package sexpr

type SexprParamKind int

const (
	SexprParamKindString SexprParamKind = iota
	SexprParamKindSexpr
)
