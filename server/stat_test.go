package server

import (
	"errors"
	"testing"
	"time"

	"github.com/lazybark/go-tls-server/conn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatistic(t *testing.T) {
	srv := GetEmptyTestServer()
	srv.connPool["test_conn"] = &conn.Connection{}
	srv.connPool["test_conn2"] = &conn.Connection{}

	e := 15
	r := 44
	s := 72

	srv.addErrors(e)
	srv.addRecBytes(r)
	srv.addSentBytes(s)

	now := time.Now()
	sentBytes, recievedBytes, errs, err := srv.Stats(now.Year(), int(now.Month()), now.Day())
	require.NoError(t, err)

	assert.Equal(t, e, errs)
	assert.Equal(t, r, recievedBytes)
	assert.Equal(t, s, sentBytes)

	connections, err := srv.StatsConnections()
	require.NoError(t, err)
	assert.Equal(t, 2, connections)

	e2 := 15
	r2 := 44
	s2 := 72

	srv.addErrors(e2)
	srv.addRecBytes(r2)
	srv.addSentBytes(s2)
	srv.connPool["test_conn3"] = &conn.Connection{}

	sentBytes, recievedBytes, errs, err = srv.Stats(now.Year(), int(now.Month()), now.Day())
	require.NoError(t, err)

	assert.Equal(t, e+e2, errs)
	assert.Equal(t, r+r2, recievedBytes)
	assert.Equal(t, s+s2, sentBytes)

	connections, err = srv.StatsConnections()
	require.NoError(t, err)
	assert.Equal(t, 3, connections)

	// Adding negative does nothing.
	e3 := -15
	r3 := -44
	s3 := -72

	srv.addErrors(e3)
	srv.addRecBytes(r3)
	srv.addSentBytes(s3)
	srv.connPool["test_conn3"] = &conn.Connection{}

	sentBytes, recievedBytes, errs, err = srv.Stats(now.Year(), int(now.Month()), now.Day())
	require.NoError(t, err)

	assert.Equal(t, e+e2, errs)
	assert.Equal(t, r+r2, recievedBytes)
	assert.Equal(t, s+s2, sentBytes)

	connections, err = srv.StatsConnections()
	require.NoError(t, err)
	assert.Equal(t, 3, connections)

	// Error if stat doesn't exist for the day.
	now = time.Now().AddDate(1, 0, 0)
	sentBytes, recievedBytes, errs, err = srv.Stats(now.Year(), int(now.Month()), now.Day())
	assert.Equal(t, true, errors.Is(ErrNoStatForTheDay, err))

	assert.Equal(t, 0, errs)
	assert.Equal(t, 0, recievedBytes)
	assert.Equal(t, 0, sentBytes)

	now = time.Now().AddDate(-5, 0, 0)
	sentBytes, recievedBytes, errs, err = srv.Stats(now.Year(), int(now.Month()), now.Day())
	assert.Equal(t, true, errors.Is(ErrNoStatForTheDay, err))

	assert.Equal(t, 0, errs)
	assert.Equal(t, 0, recievedBytes)
	assert.Equal(t, 0, sentBytes)
}
