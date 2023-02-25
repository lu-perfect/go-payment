package util

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestRandomInt(t *testing.T) {
	min := int64(100)
	max := int64(1000)

	res := RandomInt(min, max)
	require.True(t, res >= min)
	require.True(t, res <= max)
}

func TestRandomString(t *testing.T) {
	n := 10

	res := RandomString(n)
	require.NotEmpty(t, res)
	require.Equal(t, n, len(res))
}

func TestRandomName(t *testing.T) {
	res := RandomName()
	require.NotEmpty(t, res)
	require.Equal(t, 6, len(res))
}

func TestRandomEmail(t *testing.T) {
	res := RandomEmail()
	require.NotEmpty(t, res)
	require.True(t, strings.Contains(res, "@"))
}

func TestRandomMoney(t *testing.T) {
	res := RandomMoney()
	require.NotEmpty(t, res)
	require.True(t, res >= 0)
}

func TestRandomCurrency(t *testing.T) {
	res := RandomCurrency()
	require.NotEmpty(t, res)
	require.True(t, IsSupportedCurrency(res))
}
