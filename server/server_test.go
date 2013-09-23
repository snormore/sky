package server

func init() {
	// Standardize servlet count for tests so we don't get different
	// results on different machines.
	defaultServletCount = 16
}
