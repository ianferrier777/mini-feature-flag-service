package flags

// EvaluateFlag checks if a flag is enabled for a specific user and region
func EvaluateFlag(name, userId, region string) (bool, string, bool) {
	flag, found := GetFlag(name)
	if !found {
		return false, "", false
	}

	// If the user is explicitly targeted
	for _, uid := range flag.TargetUsers {
		if uid == userId {
			return true, "User explicitly targeted", true
		}
	}

	// If the region is explicitly targeted
	for _, r := range flag.TargetRegions {
		if r == region {
			return true, "Region explicitly targeted", true
		}
	}

	// Otherwise, fallback to the global flag setting
	return flag.Enabled, "Global flag setting", true
}
