package ast

type Context int32

const (
	Context_None     Context = 0
	Context_True     Context = 1
	Context_False    Context = 2
	Context_Select   Context = 4
	Context_Update   Context = 8
	Context_AConst   Context = 16
	Context_FuncCall Context = 32
	Context_TypeName Context = 64
	Context_Operator Context = 128
)
