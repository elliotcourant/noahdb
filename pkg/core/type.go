package core

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/readystock/goqu"
	"strconv"
	"strings"
)

func (typ Type) PostgresName() string {
	return ""
}

var (
	getTypeQuery = goqu.
		From("types").
		Select("types.type_id")
)

type typeContext struct {
	*base
}

type TypeContext interface {
	GetTypeByName(name string) (Type, bool, error)
}

func (ctx *base) Types() TypeContext {
	return &typeContext{
		ctx,
	}
}

func (ctx *typeContext) parseArray(name string) (string, error) {
	if strings.HasPrefix(name, "[") {
		i := strings.IndexRune(name, ']')
		arraySize := name[1:i]
		if arraySize != "" {
			size, err := strconv.Atoi(arraySize)
			if err != nil {
				return name, fmt.Errorf("could not parse array bounds: %v", err)
			}
			golog.Infof("array size: %d", size)
		}
		name = name[i+1:]
		name = fmt.Sprintf("_%s", name)
	}
	return name, nil
}

func (ctx *typeContext) parseTimes(name string) (string, error) {
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

func (ctx *typeContext) GetTypeByName(name string) (Type, bool, error) {
	name = strings.ToLower(name)

	name, err := ctx.parseArray(name)
	if err != nil {
		return Type_unknown, false, err
	}

	name, err = ctx.parseTimes(name)
	if err != nil {
		return Type_unknown, false, err
	}

	compiledSql, _, _ := getTypeQuery.
		LeftJoin(
			goqu.I("type_aliases"),
			goqu.On(goqu.I("type_aliases.type_id").Eq(goqu.I("types.type_id")))).
		Where(
			goqu.Or(
				goqu.Ex{
					"types.type_name": name,
				},
				goqu.Ex{
					"type_aliases.alias_name": name,
				},
			),
		).
		Limit(1).
		ToSql()
	rows, err := ctx.db.Query(compiledSql)
	if err != nil {
		return Type_unknown, false, err
	}
	ids, err := idArray(rows)
	if err != nil {
		return Type_unknown, false, err
	}
	if len(ids) == 0 {
		return Type_unknown, false, nil
	}
	return Type(ids[0]), true, nil
}
