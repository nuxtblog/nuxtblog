package pluginsys

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// dbDialect returns the normalized dialect key for the current database.
func dbDialect() string {
	t := strings.ToLower(g.DB().GetConfig().Type)
	switch t {
	case "pgsql", "postgres", "postgresql":
		return "pgsql"
	case "mysql", "mariadb":
		return "mysql"
	default:
		return "sqlite"
	}
}

// resolveSQL picks the SQL string for the current database from a dialect map.
// Only exact dialect match is accepted; no fallback.
func resolveSQL(dialects map[string]string) (string, error) {
	if len(dialects) == 0 {
		return "", nil
	}
	d := dbDialect()
	if sql, ok := dialects[d]; ok {
		return sql, nil
	}
	return "", fmt.Errorf("no migration SQL for dialect %q (available: %v)",
		d, mapKeys(dialects))
}

func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
