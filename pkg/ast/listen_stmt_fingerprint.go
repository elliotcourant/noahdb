// Auto-generated - DO NOT EDIT

package ast

func (node ListenStmt) Fingerprint(ctx FingerprintContext, parentNode Node, parentFieldName string) {
	ctx.WriteString("ListenStmt")

	if node.Conditionname != nil {
		ctx.WriteString("conditionname")
		ctx.WriteString(*node.Conditionname)
	}
}
