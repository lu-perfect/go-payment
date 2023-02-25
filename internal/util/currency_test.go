package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsSupportedCurrency(t *testing.T) {
	res := IsSupportedCurrency("Invalid Parameter")
	require.Equal(t, false, res)

	res = IsSupportedCurrency(RUB)
	require.Equal(t, true, res)

	res = IsSupportedCurrency(EUR)
	require.Equal(t, true, res)

	res = IsSupportedCurrency(USD)
	require.Equal(t, true, res)
}

// TODO: test register validator
