// Package ports defines application boundary interfaces (hexagonal ports).
package ports

// Logger is a minimal structured logging facade used across packages.
type Logger interface {
	Sync() error
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Debug(msg string, keysAndValues ...any)
}
