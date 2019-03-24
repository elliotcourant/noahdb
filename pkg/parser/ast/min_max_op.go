// Auto-generated from postgres/src/include/nodes/primnodes.h - DO NOT EDIT

package ast

/*
 * MinMaxExpr - a GREATEST or LEAST function
 */
type MinMaxOp uint

const (
	IS_GREATEST MinMaxOp = iota
	IS_LEAST
)
