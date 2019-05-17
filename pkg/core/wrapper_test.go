package core

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func Test_NetworkResolution(t *testing.T) {
	addr, err := externalIP()
	assert.NoError(t, err)
	golog.Infof("%v", addr)
}
