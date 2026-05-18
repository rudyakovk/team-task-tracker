package migrations

import "testing"

func TestExtractUpSQL(t *testing.T) {
	t.Parallel()

	sql, err := extractUpSQL(`
-- +goose Up
CREATE TABLE example (id uuid PRIMARY KEY);

-- +goose Down
DROP TABLE example;
`)
	if err != nil {
		t.Fatalf("extract up sql: %v", err)
	}

	want := "CREATE TABLE example (id uuid PRIMARY KEY);"
	if sql != want {
		t.Fatalf("unexpected sql:\nwant: %s\n got: %s", want, sql)
	}
}

func TestExtractUpSQLRequiresMarker(t *testing.T) {
	t.Parallel()

	if _, err := extractUpSQL("CREATE TABLE example (id uuid PRIMARY KEY);"); err == nil {
		t.Fatal("expected marker error")
	}
}
