package inspector2util

func (fs Findings) FilterImageHashes(hashesIncl []string) Findings {
	var out Findings
	hashesInclMap := map[string]int{}
	for _, h := range hashesIncl {
		hashesInclMap[h]++
	}
	for _, f := range fs {
		fx := Finding(f)
		imgHashes := fx.ImageHashes()
		for _, h := range imgHashes {
			if _, ok := hashesInclMap[h]; ok {
				out = append(out, f)
				break
			}
		}
	}
	return out
}

func (fs Findings) FilterPOMPropertiesExcl() Findings {
	var out Findings
	for _, f := range fs {
		fx := Finding(f)
		if !fx.FilePathsInclPOMProperties() {
			out = append(out, f)
		}
	}
	return out
}
