package types

import (
	"fmt"
	"github.com/elliotcourant/timber"
	"strconv"
	"strings"
)

// GetTypeByName will return a NoahDB/PostgreSQL type for the name provided. If a matching type
// cannot be found then this will return false. If a problem is encountered in parsing the type name
// then an error will be returned.
func GetTypeByName(name string) (Type, bool, error) {
	name = strings.ToLower(name)
	name = strings.TrimPrefix(name, "pg_catalog.")

	name, err := parseArrayTypeName(name)
	if err != nil {
		return Type_unknown, false, err
	}

	name, err = parseTimeTypeName(name)
	if err != nil {
		return Type_unknown, false, err
	}

	t, ok := typesByName[name]
	return t, ok, nil
}

// GetTypeByOid will return a type if the OID provided is valid for NoahDB and PostgreSQL. If the
// OID provided is not valid or is not handled then GetTypeByOid will return false.
func GetTypeByOid(oid OID) (Type, bool) {
	i := int(oid)
	t := Type(i)
	s := strconv.Itoa(i)
	if s == t.String() {
		return t, false
	}
	return t, true
}

func parseArrayTypeName(name string) (string, error) {
	if strings.HasSuffix(name, "]") {
		i := strings.IndexRune(name, '[')
		arraySize := name[i+1 : len(name)-1]
		if arraySize != "" {
			// TODO (elliotcourant) find a way to return the parsed array size if it's specified.
			size, err := strconv.Atoi(arraySize)
			if err != nil {
				return name, fmt.Errorf("could not parse array bounds: %v", err)
			}
			timber.Verbosef("array size: %d", size)
		}
		name = name[:i]
		name = fmt.Sprintf("_%s", name)
	}
	return name, nil
}

func parseTimeTypeName(name string) (string, error) {
	i := strings.IndexRune(name, ' ')
	if i < 0 {
		return name, nil
	}
	first, second := name[:i], name[i+1:]
	switch first {
	case "time", "_time", "timestamp", "_timestamp":
		if strings.HasSuffix(second, "without time zone") {
			return first, nil
		} else if strings.HasSuffix(second, "with time zone") {
			return fmt.Sprintf("%s with time zone", first), nil
		} else {
			return first, nil
		}
	case "interval", "_interval":
	default:
		return name, nil
	}
	return name, nil
}
