// Package sql embeds the database schema files so they can be executed
// automatically at startup without requiring any manual SQL import step.
package sql

import _ "embed"

// SQLite is the full SQLite schema (CREATE TABLE IF NOT EXISTS …).
//
//go:embed init.sql
var SQLite string

// PostgreSQL is the full PostgreSQL schema (CREATE TABLE IF NOT EXISTS …).
//
//go:embed init_postgres.sql
var PostgreSQL string
