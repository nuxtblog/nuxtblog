package pluginsys

import (
	"fmt"
	"regexp"
	"strings"
)

// sqlGuard validates SQL statements against a plugin's declared DB capabilities.
type sqlGuard struct {
	prefix string     // e.g. "plugin_nuxtblog_plugin_hello_js_"
	caps   *DBCap
	trust  TrustLevel
}

// identPattern captures a table identifier (backtick-quoted, double-quoted, or bare).
const identPattern = `(?:` + "`" + `([^` + "`" + `]+)` + "`" + `|"([^"]+)"|([a-zA-Z_][a-zA-Z0-9_]*))`

var (
	tablePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bFROM\s+` + identPattern),
		regexp.MustCompile(`(?i)\bJOIN\s+` + identPattern),
		regexp.MustCompile(`(?i)\bINTO\s+` + identPattern),
		regexp.MustCompile(`(?i)\bUPDATE\s+` + identPattern),
		regexp.MustCompile(`(?i)\bTABLE\s+(?:IF\s+(?:NOT\s+)?EXISTS\s+)?` + identPattern),
	}
	opPattern = regexp.MustCompile(`(?i)^\s*(SELECT|INSERT|UPDATE|DELETE|CREATE|ALTER|DROP)\b`)
)

func (g *sqlGuard) validate(sql string) error {
	// Official trust level bypasses all checks
	if g.trust == TrustLevelOfficial {
		return nil
	}

	// Raw capability bypasses checks
	if g.caps != nil && g.caps.Raw {
		return nil
	}

	// No DB capability declared — deny everything
	if g.caps == nil {
		return fmt.Errorf("no db capability declared")
	}

	// Strip string literals to prevent false matches
	cleaned := stripStringLiterals(sql)

	// Reject multi-statement queries (semicolons)
	if strings.Contains(cleaned, ";") {
		return fmt.Errorf("multi-statement queries are not allowed")
	}

	// Extract operation type
	op := extractOp(cleaned)
	if op == "" {
		return fmt.Errorf("unable to determine SQL operation")
	}

	// Extract all referenced table names
	tables := extractTableNames(cleaned)
	if len(tables) == 0 {
		// For DDL without detectable tables, deny
		if isDDL(op) {
			return fmt.Errorf("cannot determine target table for %s", op)
		}
		// DML without tables (e.g. SELECT 1) — allow if caps exist
		return nil
	}

	// Validate each table
	for _, table := range tables {
		lower := strings.ToLower(table)

		// Own-prefixed tables
		if strings.HasPrefix(lower, strings.ToLower(g.prefix)) {
			if g.caps.Own {
				continue
			}
			return fmt.Errorf("table %s: own-table access not granted", table)
		}

		// DDL on non-own tables is always denied
		if isDDL(op) {
			return fmt.Errorf("DDL on table %s is not allowed", table)
		}

		// Check whitelist
		if !g.isTableOpAllowed(lower, op) {
			return fmt.Errorf("access to table %s with operation %s is not allowed", table, op)
		}
	}

	return nil
}

// isTableOpAllowed checks if the given table+op is in the whitelist.
func (g *sqlGuard) isTableOpAllowed(table, op string) bool {
	lowerOp := strings.ToLower(op)
	for _, rule := range g.caps.Tables {
		if strings.ToLower(rule.Table) == table {
			for _, allowedOp := range rule.Ops {
				if strings.ToLower(allowedOp) == lowerOp {
					return true
				}
			}
			return false
		}
	}
	return false
}

// extractOp returns the SQL operation keyword (SELECT, INSERT, etc.).
func extractOp(sql string) string {
	m := opPattern.FindStringSubmatch(sql)
	if len(m) < 2 {
		return ""
	}
	return strings.ToUpper(m[1])
}

// extractTableNames returns all table names referenced in the SQL.
func extractTableNames(sql string) []string {
	seen := make(map[string]bool)
	var tables []string

	for _, pat := range tablePatterns {
		matches := pat.FindAllStringSubmatch(sql, -1)
		for _, m := range matches {
			// The table name is in one of the capture groups
			name := ""
			for i := 1; i < len(m); i++ {
				if m[i] != "" {
					name = m[i]
					break
				}
			}
			if name != "" && !seen[strings.ToLower(name)] {
				seen[strings.ToLower(name)] = true
				tables = append(tables, name)
			}
		}
	}
	return tables
}

// isDDL returns true if the operation is a DDL statement.
func isDDL(op string) bool {
	switch op {
	case "CREATE", "ALTER", "DROP":
		return true
	}
	return false
}

// stripStringLiterals replaces string literals with empty strings to prevent
// table name patterns inside string values from being matched.
func stripStringLiterals(sql string) string {
	var b strings.Builder
	b.Grow(len(sql))
	i := 0
	for i < len(sql) {
		ch := sql[i]
		if ch == '\'' || ch == '"' {
			quote := ch
			b.WriteByte(quote)
			i++
			for i < len(sql) {
				if sql[i] == quote {
					if i+1 < len(sql) && sql[i+1] == quote {
						// Escaped quote
						i += 2
						continue
					}
					break
				}
				i++
			}
			b.WriteByte(quote)
			if i < len(sql) {
				i++
			}
		} else {
			b.WriteByte(ch)
			i++
		}
	}
	return b.String()
}
