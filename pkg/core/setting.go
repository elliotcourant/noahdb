package core

type settingContext struct {
	*base
}

type SettingContext interface {
	GetSetting(SettingKey, interface{}) (bool, error)
}

func (ctx *base) Settings() SettingContext {
	return &settingContext{
		ctx,
	}
}

func (ctx *settingContext) GetSetting(key SettingKey, result interface{}) (bool, error) {
	rows, err := ctx.db.Query("SELECT types.id, ")
	panic("not implemented")
}
