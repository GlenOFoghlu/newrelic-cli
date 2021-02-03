// +build unit

package decode

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestDecodeCommand(t *testing.T) {
	assert.Equal(t, "decode", Command.Name())

	testcobra.CheckCobraMetadata(t, cmdDecode)
	testcobra.CheckCobraRequiredFlags(t, cmdDecode, []string{})
	testcobra.CheckCobraCommandAliases(t, cmdDecode, []string{})
}
