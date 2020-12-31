/**
Copyright (C) 2020 Aaron Sky.

This file is part of Cider, a tool for automating submission
of apps to Apple's App Stores.

Cider is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Cider is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Cider.  If not, see <http://www.gnu.org/licenses/>.
*/

package clicommand

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	exit := func(code int) {
		assert.Equal(t, 0, code)
	}
	err := os.Setenv("CI", "1")
	assert.NoError(t, err)
	Execute("TEST", exit, []string{"help", "--debug"})
}

func TestRootCmd_Error(t *testing.T) {
	exit := func(code int) {
		assert.NotEqual(t, 0, code)
	}
	Execute("TEST", exit, []string{"check"})
}