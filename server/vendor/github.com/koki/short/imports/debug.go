package imports

// TrimToDepth returns true if it trimmed away any imports.
func TrimToDepth(module *Module, depth int) bool {
	trimmed := false
	if depth <= 0 {
		trimmed = len(module.Imports) > 0
		module.Imports = nil
		return trimmed
	}

	for _, imprt := range module.Imports {
		trimmed = trimmed || TrimToDepth(imprt.Module, depth-1)
	}

	return trimmed
}
