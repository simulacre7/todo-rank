package scan

// ScanOptions holds all configuration for a scan operation.
type ScanOptions struct {
	Root     string
	Ignore   []string
	Format   string
	OutPath  string
	MinScore int
	Tags     []string
}
