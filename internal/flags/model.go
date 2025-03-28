package flags

type Flag struct {
	Name           string
	Enabled        bool     // Global flag toggle
	TargetUsers    []string // Whitelist of user IDs
	TargetRegions  []string // Whitelist of allowed regions
}