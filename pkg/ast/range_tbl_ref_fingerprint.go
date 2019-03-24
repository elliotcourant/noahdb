// Auto-generated - DO NOT EDIT

package ast

import "strconv"

func (node RangeTblRef) Fingerprint(ctx FingerprintContext, parentNode Node, parentFieldName string) {
	ctx.WriteString("RangeTblRef")

	if node.Rtindex != 0 {
		ctx.WriteString("rtindex")
		ctx.WriteString(strconv.Itoa(int(node.Rtindex)))
	}
}
