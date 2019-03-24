// Auto-generated - DO NOT EDIT

package ast

func (node ColumnRef) Fingerprint(ctx FingerprintContext, parentNode Node, parentFieldName string) {
	ctx.WriteString("ColumnRef")

	if len(node.Fields.Items) > 0 {
		subCtx := FingerprintSubContext{}
		node.Fields.Fingerprint(&subCtx, node, "Fields")

		if len(subCtx.parts) > 0 {
			ctx.WriteString("fields")
			for _, part := range subCtx.parts {
				ctx.WriteString(part)
			}
		}
	}
	// Intentionally ignoring node.Location for fingerprinting
}
