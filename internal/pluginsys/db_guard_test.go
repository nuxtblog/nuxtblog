package pluginsys

import (
	"encoding/json"
	"testing"
)

func TestSQLGuard(t *testing.T) {
	prefix := "plugin_nuxtblog_plugin_test_"

	tests := []struct {
		name    string
		sql     string
		caps    *DBCap
		trust   TrustLevel
		wantErr bool
	}{
		{
			name:    "basic SELECT own table with own:true",
			sql:     "SELECT * FROM plugin_nuxtblog_plugin_test_data",
			caps:    &DBCap{Own: true},
			trust:   TrustLevelCommunity,
			wantErr: false,
		},
		{
			name:    "SELECT core table without whitelist",
			sql:     "SELECT * FROM users",
			caps:    &DBCap{Own: true},
			trust:   TrustLevelCommunity,
			wantErr: true,
		},
		{
			name:    "whitelisted SELECT on posts",
			sql:     "SELECT * FROM posts",
			caps:    &DBCap{Own: true, Tables: []DBTableRule{{Table: "posts", Ops: []string{"select"}}}},
			trust:   TrustLevelCommunity,
			wantErr: false,
		},
		{
			name:    "INSERT on posts with only select whitelist",
			sql:     "INSERT INTO posts (title) VALUES ('test')",
			caps:    &DBCap{Own: true, Tables: []DBTableRule{{Table: "posts", Ops: []string{"select"}}}},
			trust:   TrustLevelCommunity,
			wantErr: true,
		},
		{
			name:    "DDL on own table",
			sql:     "CREATE TABLE plugin_nuxtblog_plugin_test_foo (id INTEGER PRIMARY KEY)",
			caps:    &DBCap{Own: true},
			trust:   TrustLevelCommunity,
			wantErr: false,
		},
		{
			name:    "DDL on core table",
			sql:     "DROP TABLE posts",
			caps:    &DBCap{Own: true, Tables: []DBTableRule{{Table: "posts", Ops: []string{"select"}}}},
			trust:   TrustLevelCommunity,
			wantErr: true,
		},
		{
			name:    "multi-statement rejected",
			sql:     "SELECT 1; DROP TABLE users",
			caps:    &DBCap{Own: true},
			trust:   TrustLevelCommunity,
			wantErr: true,
		},
		{
			name:    "string literal not misdetected",
			sql:     "SELECT * FROM plugin_nuxtblog_plugin_test_data WHERE name = 'FROM users'",
			caps:    &DBCap{Own: true},
			trust:   TrustLevelCommunity,
			wantErr: false,
		},
		{
			name:    "JOIN multi-table with whitelist",
			sql:     "SELECT * FROM plugin_nuxtblog_plugin_test_data JOIN posts ON plugin_nuxtblog_plugin_test_data.post_id = posts.id",
			caps:    &DBCap{Own: true, Tables: []DBTableRule{{Table: "posts", Ops: []string{"select"}}}},
			trust:   TrustLevelCommunity,
			wantErr: false,
		},
		{
			name:    "official trust bypasses all checks",
			sql:     "DROP TABLE users; DELETE FROM posts",
			caps:    nil,
			trust:   TrustLevelOfficial,
			wantErr: false,
		},
		{
			name:    "nil caps denies everything",
			sql:     "SELECT 1",
			caps:    nil,
			trust:   TrustLevelCommunity,
			wantErr: true,
		},
		{
			name:    "raw capability bypasses checks",
			sql:     "SELECT * FROM users",
			caps:    &DBCap{Raw: true},
			trust:   TrustLevelCommunity,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &sqlGuard{prefix: prefix, caps: tt.caps, trust: tt.trust}
			err := g.validate(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBCapUnmarshalJSON(t *testing.T) {
	t.Run("legacy bool true", func(t *testing.T) {
		var cap DBCap
		if err := json.Unmarshal([]byte(`true`), &cap); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if !cap.Own {
			t.Error("expected Own=true for legacy 'true'")
		}
	})

	t.Run("legacy bool false", func(t *testing.T) {
		var cap DBCap
		if err := json.Unmarshal([]byte(`false`), &cap); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if cap.Own {
			t.Error("expected Own=false for legacy 'false'")
		}
	})

	t.Run("structured object", func(t *testing.T) {
		var cap DBCap
		data := `{"own":true,"tables":[{"table":"posts","ops":["select"]}]}`
		if err := json.Unmarshal([]byte(data), &cap); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if !cap.Own {
			t.Error("expected Own=true")
		}
		if len(cap.Tables) != 1 || cap.Tables[0].Table != "posts" {
			t.Error("expected one table rule for posts")
		}
	})

	t.Run("in PluginCapabilities struct", func(t *testing.T) {
		var caps PluginCapabilities
		data := `{"db":true}`
		if err := json.Unmarshal([]byte(data), &caps); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if caps.DB == nil {
			t.Fatal("expected DB to be non-nil")
		}
		if !caps.DB.Own {
			t.Error("expected DB.Own=true for legacy 'true'")
		}
	})
}

func TestStripStringLiterals(t *testing.T) {
	input := `SELECT * FROM t WHERE name = 'FROM users' AND id = 1`
	result := stripStringLiterals(input)
	if result != `SELECT * FROM t WHERE name = '' AND id = 1` {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestExtractTableNames(t *testing.T) {
	tests := []struct {
		sql    string
		tables []string
	}{
		{"SELECT * FROM users", []string{"users"}},
		{"SELECT * FROM users JOIN posts ON users.id = posts.author_id", []string{"users", "posts"}},
		{"INSERT INTO plugin_test_data (id) VALUES (1)", []string{"plugin_test_data"}},
		{"UPDATE users SET name = 'foo'", []string{"users"}},
		{"CREATE TABLE IF NOT EXISTS my_table (id INT)", []string{"my_table"}},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			got := extractTableNames(tt.sql)
			if len(got) != len(tt.tables) {
				t.Fatalf("expected %d tables, got %d: %v", len(tt.tables), len(got), got)
			}
			for i, want := range tt.tables {
				if got[i] != want {
					t.Errorf("table[%d] = %s, want %s", i, got[i], want)
				}
			}
		})
	}
}
