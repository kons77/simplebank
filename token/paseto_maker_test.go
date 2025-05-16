package token

import (
	"testing"
	"time"

	"github.com/kons77/simplebank/util"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	key := util.RandomBytes(32)

	maker, err := NewPasetoMaker(key, nil)
	require.NoError(t, err)

	username := util.RandomUsername()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)

}

func TestExpiredPasetoToken(t *testing.T) {
	key := util.RandomBytes(32)

	maker, err := NewPasetoMaker(key, nil)
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomUsername(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidKeySize(t *testing.T) {
	key := util.RandomBytes(30)

	_, err := NewPasetoMaker(key, nil)
	require.Error(t, err)
}

func TestInvalidToken(t *testing.T) {
	key := util.RandomBytes(32)

	maker, err := NewPasetoMaker(key, nil)
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomUsername(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	token += "aa"

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestWrongKey(t *testing.T) {
	key1 := util.RandomBytes(32)
	maker1, err := NewPasetoMaker(key1, nil)
	require.NoError(t, err)

	key2 := util.RandomBytes(32)
	maker2, err := NewPasetoMaker(key2, nil)
	require.NoError(t, err)

	token, err := maker1.CreateToken(util.RandomUsername(), time.Minute)
	require.NoError(t, err)

	payload, err := maker2.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
