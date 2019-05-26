package core

import (
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/readystock/goqu"
)

type settingContext struct {
	*base
}

type SettingContext interface {
	GetSetting(SettingKeyOptions) (Setting, bool, error)
	GetSettingValue(key SettingKeyOptions) (interface{}, bool, error)
}

func (ctx *base) Setting() SettingContext {
	return &settingContext{
		ctx,
	}
}

//
// func (ctx *settingContext) SetSetting(key SettingKeyOptions, value interface{}) error {
// 	key.EnumDescriptor()
// }

func (ctx *settingContext) GetSettingValue(key SettingKeyOptions) (interface{}, bool, error) {
	val, ok, err := ctx.GetSetting(key)
	if err != nil || !ok {
		return nil, ok, err
	}
	if val.IValue != nil {
		v, ok := val.IValue.(*Setting_IntegerValue)
		if ok {
			return v.IntegerValue, true, nil
		}
		return nil, false, nil
	} else if val.BValue != nil {
		v, ok := val.BValue.(*Setting_BooleanValue)
		if ok {
			return v.BooleanValue, true, nil
		}
		return nil, false, nil
	} else if val.TValue != nil {
		v, ok := val.TValue.(*Setting_TextValue)
		if ok {
			return v.TextValue, true, nil
		}
		return nil, false, nil
	}
	return nil, ok, err
}

func (ctx *settingContext) GetSetting(key SettingKeyOptions) (Setting, bool, error) {
	compiledSql, _, _ := goqu.
		From("settings").
		Select("*").
		Where(goqu.Ex{
			"setting_id": key,
		}).
		Limit(1).
		ToSql()
	rows, err := ctx.db.Query(compiledSql)
	settings, err := ctx.settingsFromRows(rows)
	if err != nil {
		return Setting{}, false, err
	}
	if len(settings) == 0 {
		return Setting{}, false, nil
	}
	return settings[0], true, nil
}

func (ctx *settingContext) settingsFromRows(response *frunk.QueryResponse) ([]Setting, error) {
	rows := rqliter.NewRqlRows(response)
	settings := make([]Setting, 0)
	for rows.Next() {
		setting := Setting{}
		var (
			ivalue *int64  = nil
			bvalue *bool   = nil
			tvalue *string = nil
		)
		if err := rows.Scan(
			&setting.SettingID,
			&setting.SettingKey,
			&setting.TypeID,
			&ivalue,
			&bvalue,
			&tvalue,
		); err != nil {
			return nil, err
		}
		if ivalue != nil {
			setting.IValue = &Setting_IntegerValue{IntegerValue: *ivalue}
		} else if bvalue != nil {
			setting.BValue = &Setting_BooleanValue{BooleanValue: *bvalue}
		} else if tvalue != nil {
			setting.TValue = &Setting_TextValue{TextValue: *tvalue}
		}
		settings = append(settings, setting)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return settings, nil
}
