// Auto-generated - DO NOT EDIT

package pg_query

func (node List) Deparse(ctx Context) (*string, error) {
	panic("use DeparseList for lists")
}

func (list List) DeparseList(ctx Context) ([]string, error) {
	out := make([]string, len(list.Items))
	for i, node := range list.Items {

		if str, err := node.Deparse(ctx); err != nil {
			return nil, err
		} else {
			out[i] = *str
		}
	}
	return out, nil
}
