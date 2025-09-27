package database

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"template/internal/config"
)

func TestMain(m *testing.M) {
	dbConfig := os.Getenv("DB_CONFIG")
	err := yaml.Unmarshal([]byte(dbConfig), &config.C.Database)
	if err != nil {
		fmt.Println("failed to unmarshal db config", err)
		os.Exit(1)
	}

	err = Connect(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestDBConnection(t *testing.T) {
	t.Parallel()

	var (
		expNum int64 = rand.Int64()
		actNum int64 = -1
	)

	rows := pool.QueryRow(t.Context(), "SELECT $1::BIGINT;", expNum)
	err := rows.Scan(&actNum)

	require.NoError(t, err)
	require.Equal(t, expNum, actNum)
}
