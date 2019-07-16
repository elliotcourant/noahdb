package types

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
)

func Decode(format pgwirebase.FormatCode, p Type, data []byte) (Value, error) {
	var t interface{}
	switch p {
	case Type_aclitem_array:
		t = &ACLItemArray{}
	case Type_bool_array:
		t = &BoolArray{}
	case Type_bpchar_array, Type_char_array:
		t = &BPCharArray{}
	case Type_bytea_array:
		t = &ByteaArray{}
	case Type_cidr_array:
		t = &CIDRArray{}
	case Type_date_array:
		t = &DateArray{}
	case Type_float4_array:
		t = &Float4Array{}
	case Type_hstore_array:
		t = &HstoreArray{}
	case Type_int2_array:
		t = &Int2Array{}
	case Type_int4_array:
		t = &Int4Array{}
	case Type_int8_array:
		t = &Int8Array{}
	case Type_numeric_array:
		t = &NumericArray{}
	case Type_text_array:
		t = &TextArray{}
	case Type_time_array, Type_timestamp_array:
		t = &TimestampArray{}
	case Type_timetz_array, Type_timestamptz_array:
		t = &TimestamptzArray{}
	case Type_uuid_array:
		t = &UUIDArray{}
	case Type_varchar_array:
		t = &VarcharArray{}
	case Type_aclitem:
		t = &ACLItem{}
	case Type_bool:
		t = &Bool{}
	case Type_box:
		t = &Box{}
	case Type_bpchar:
		t = &BPChar{}
	case Type_bytea:
		t = &Bytea{}
	case Type_cid:
		t = &CID{}
	case Type_cidr:
		t = &CIDR{}
	case Type_circle:
		t = &Circle{}
	case Type_date:
		t = &Date{}
	case Type_daterange:
		t = &Daterange{}
	case Type_float4:
		t = &Float4{}
	case Type_float8:
		t = &Float8{}
	case Type_hstore:
		t = &Hstore{}
	case Type_inet:
		t = &Inet{}
	case Type_int2:
		t = &Int2{}
	case Type_int4:
		t = &Int4{}
	case Type_int8:
		t = &Int8{}
	case Type_interval:
		t = &Interval{}
	case Type_json:
		t = &JSON{}
	case Type_jsonb:
		t = &JSONB{}
	case Type_line:
		t = &Line{}
	case Type_macaddr:
		t = &Macaddr{}
	case Type_numeric:
		t = &Numeric{}
	case Type_oid:
		o := OID(0)
		t = &o
	case Type_path:
		t = &Path{}
	case Type_point:
		t = &Point{}
	case Type_polygon:
		t = &Polygon{}
	case Type_text:
		t = &Text{}
	case Type_time, Type_timestamp:
		t = &Timestamp{}
	case Type_timetz, Type_timestamptz:
		t = &Timestamptz{}
	case Type_uuid:
		t = &UUID{}
	case Type_varbit:
		t = &Varbit{}
	case Type_varchar:
		t = &Varchar{}
	case Type_xid:
		t = &XID{}
	default:
		return nil, fmt.Errorf("cannot handle type [%s]", p)
	}

	if t == nil {
		return nil, fmt.Errorf("could not determine type [%s]", p)
	}

	switch format {
	case pgproto.TextFormat:
		if textDecoder, ok := t.(TextDecoder); !ok {
			return nil, fmt.Errorf("cannot decode [%s] format [%s] not implemented", p, format)
		} else if err := textDecoder.DecodeText(nil, data); err != nil {
			return nil, fmt.Errorf("failed to decode [%s] format [%s]: %v", p, format, err)
		}
	case pgproto.BinaryFormat:
		if binaryDecoder, ok := t.(BinaryDecoder); !ok {
			return nil, fmt.Errorf("cannot decode [%s] format [%s] not implemented", p, format)
		} else if err := binaryDecoder.DecodeBinary(nil, data); err != nil {
			return nil, fmt.Errorf("failed to decode [%s] format [%s]: %v", p, format, err)
		}
	}

	return t.(Value), nil
}
