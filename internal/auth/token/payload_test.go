package token

import (
	"github.com/stretchr/testify/require"
	"gobank/internal/util"
	"testing"
	"time"
)

func createPayload(t *testing.T, duration time.Duration) *Payload {
	username := util.RandomName()

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)

	return payload
}

func TestPayloadValid(t *testing.T) {
	// TODO
}
