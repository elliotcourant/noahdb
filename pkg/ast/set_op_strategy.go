// Auto-generated from postgres/src/include/nodes/nodes.h - DO NOT EDIT

package ast

type SetOpStrategy uint

const (
	SETOP_SORTED SetOpStrategy = iota /* input must be sorted */
	SETOP_HASHED                      /* use internal hashtable */
)
