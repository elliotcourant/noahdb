package types

func (x Type) PostgresName() string {
	return ""
}

func (x Type) Uint32() uint32 {
	return uint32(x)
}

func (x Type) OID() OID {
	return OID(x.Uint32())
}
