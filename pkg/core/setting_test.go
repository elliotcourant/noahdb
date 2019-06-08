package core_test

import (
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSettingContext_GetSetting(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()
	t.Run("get setting simple", func(t *testing.T) {
		setting, ok, err := colony.Setting().GetSetting(core.SettingKeyOptions_MaxPoolSize)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NotNil(t, setting.IValue)
		val, ok := setting.IValue.(*core.Setting_IntegerValue)
		assert.True(t, ok)
		assert.True(t, val.IntegerValue > 0)
	})

	t.Run("get setting value", func(t *testing.T) {
		setting, ok, err := colony.Setting().GetSettingValue(core.SettingKeyOptions_MaxPoolSize)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NotNil(t, setting)
		assert.Equal(t, int64(5), setting)
	})
}
