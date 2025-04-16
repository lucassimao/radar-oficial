# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands
- Build/Run API: `go run ./cmd/api`
- Run migrations: `go run ./cmd/migrate`
- Create migration: `make create-migration` (enter name when prompted)
- Rollback migration: `make rollback-migration`
- Connect to PostgreSQL: `make psql`

## Code Style Guidelines
- Imports: Group standard library first, then 3rd party packages, then internal packages
- Error handling: Use descriptive error messages with `fmt.Errorf("context: %w", err)`
- Naming: Use CamelCase for exported identifiers, camelCase for local variables
- Types: Use pointers for nullable fields (*string, *time.Time)
- Functions: Constructor functions use New* prefix (e.g., NewDiarioService)
- HTTP: Check errors and return appropriate status codes
- Comments: Add comments for exported functions and complex logic
- Structured logging: Use emoji prefixes for better visibility (✅, ⚠️, ❌)