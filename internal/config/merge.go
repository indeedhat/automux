package config

// mergeSessions takes two session instances and overrides fields witin the target session
// with the non nil values from the override session
func mergeSessions(target, override Session) Session {
	if override.Directory != "" {
		target.Directory = override.Directory
	}
	if override.ConfigPath != nil {
		target.ConfigPath = override.ConfigPath
	}
	if override.SessionId != "" {
		target.SessionId = override.SessionId
	}
	if override.AttachExisting != nil {
		target.AttachExisting = override.AttachExisting
	}

	target.Windows = mergeWindows(target.Windows, override.Windows)
	return target
}

// mergeWindows merges two winow slices
//
// Merges are based on window titles
//   - When a window is found in each with the same title all target fields will be
//     replaced by any non nil fields within the override slices window
//   - windows found in only the target slice will be untouched
//   - windows found in only the override slice will be appendend
func mergeWindows(target, override []Window) []Window {
	var (
		extras []Window
		merged = target[:]
	)

	for _, window := range override {
		var found bool

		for i, final := range target {
			if window.Title != target[i].Title {
				continue
			}

			if window.Exec != nil {
				final.Exec = window.Exec
			}
			if window.Focus != nil {
				final.Focus = window.Focus
			}
			if window.Directory != nil {
				final.Directory = window.Directory
			}

			final.Splits = mergeSplits(final.Splits, window.Splits)

			merged[i] = final
			found = true
		}

		if !found {
			extras = append(extras, window)
		}
	}

	return append(merged, extras...)
}

// mergeSplits merges two split slices
//
// Merges are done based on split slice indeces
//   - When a slice in the overrides shares an index with target it will have any non nil
//     fields replaced in the target split by the override one
//   - Extra slices will be appenden
func mergeSplits(target, override []Split) []Split {
	var (
		extras []Split
		merged = target[:]
	)

	for i, split := range override {
		if i >= len(target) {
			extras = append(extras, split)
			continue
		}

		final := &target[i]

		if split.Exec != nil {
			(*final).Exec = split.Exec
		}
		if split.Vertical != nil {
			(*final).Vertical = split.Vertical
		}
		if split.Size != nil {
			(*final).Size = split.Size
		}
		if split.Focus != nil {
			(*final).Focus = split.Focus
		}
		if split.Directory != nil {
			(*final).Directory = split.Directory
		}
	}

	return append(merged, extras...)
}
