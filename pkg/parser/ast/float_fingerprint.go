// Auto-generated - DO NOT EDIT

package ast

func (node Float) Fingerprint(ctx FingerprintContext, parentNode Node, parentFieldName string) {
	ctx.WriteString("Float")
	if len(node.Str) > 0 {
		ctx.WriteString("str")
		ctx.WriteString(node.Str)
	}
}
