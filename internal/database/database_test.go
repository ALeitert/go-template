package database

import (
	"context"
	"errors"
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

func TestDBTransaction(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), TransactionTimeout)
	defer cancel()

	rndString := make([]byte, 10)
	for i := range rndString {
		rndString[i] = 'a' + byte(rand.IntN(26))
	}
	tblName := "tbl_" + string(rndString)

	_, err := model.Exec(
		ctx,
		`CREATE TABLE `+tblName+` (
			key  INT  PRIMARY KEY,
			val  INT  NOT NULL
		);`,
	)
	require.NoError(t, err)

	//
	// Test 1: Successful transaction.

	t1Key := rand.Int32()
	t1Val := rand.Int32()

	err = TxExec(ctx, func(ctx context.Context, model Model) error {
		tag, err := model.Exec(ctx,
			"INSERT INTO "+tblName+"(key, val) VALUES ($1, $2);",
			t1Key, t1Val,
		)
		require.NoError(t, err)
		require.EqualValues(t, 1, tag.RowsAffected())

		return nil
	})
	require.NoError(t, err)

	row := model.QueryRow(ctx, "SELECT val FROM "+tblName+" WHERE key = $1;", t1Key)
	val := int32(-1)
	err = row.Scan(&val)
	require.NoError(t, err)

	require.EqualValues(t, t1Val, val)

	//
	// Test 2: Transaction with error.

	t2Key := rand.Int32()
	t2Val := rand.Int32()

	err = TxExec(ctx, func(ctx context.Context, model Model) error {
		tag, err := model.Exec(ctx,
			"INSERT INTO "+tblName+"(key, val) VALUES ($1, $2);",
			t2Key, t2Val,
		)
		require.NoError(t, err)
		require.EqualValues(t, 1, tag.RowsAffected())

		return errors.New("dummy error to trigger rollback")
	})
	require.Error(t, err)

	rows, err := model.Query(ctx, "SELECT val FROM "+tblName+" WHERE key = $1;", t2Key)
	require.NoError(t, err)
	defer rows.Close()

	require.False(t, rows.Next())
	require.NoError(t, rows.Err())
}
