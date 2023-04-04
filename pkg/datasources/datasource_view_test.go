package datasources

import (
	"context"
	"terraform-materialize/pkg/testhelpers"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestViewDatasource(t *testing.T) {
	r := require.New(t)

	in := map[string]interface{}{
		"schema_name":   "schema",
		"database_name": "database",
	}
	d := schema.TestResourceDataRaw(t, View().Schema, in)
	r.NotNil(d)

	testhelpers.WithMockDb(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		ir := mock.NewRows([]string{"id", "name", "schema", "database"}).
			AddRow("u1", "view", "schema", "database")
		mock.ExpectQuery(`
			SELECT
				mz_views.id,
				mz_views.name,
				mz_schemas.name,
				mz_databases.name
			FROM mz_views
			JOIN mz_schemas
				ON mz_views.schema_id = mz_schemas.id
			JOIN mz_databases
				ON mz_schemas.database_id = mz_databases.id
			WHERE mz_databases.name = 'database' AND mz_schemas.name = 'schema'`).WillReturnRows(ir)

		if err := viewRead(context.TODO(), d, db); err != nil {
			t.Fatal(err)
		}
	})

}