package cmd

// Import all supported database drivers.
// The active driver is selected by database.default.type in config.yaml.
import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
)
