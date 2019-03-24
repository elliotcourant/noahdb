// Auto-generated - DO NOT EDIT

package ast

func (node LoadStmt) Fingerprint(ctx FingerprintContext, parentNode Node, parentFieldName string) {
	ctx.WriteString("LoadStmt")

	if node.Filename != nil {
		ctx.WriteString("filename")
		ctx.WriteString(*node.Filename)
	}
}
