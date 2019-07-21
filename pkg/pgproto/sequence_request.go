package pgproto

type SequenceRequest struct {
	SequenceName string
}

func (item *SequenceRequest) Decode(src []byte) error {
	*item = SequenceRequest{}
	// buf := bytes.NewBuffer(src)
	return nil
}
