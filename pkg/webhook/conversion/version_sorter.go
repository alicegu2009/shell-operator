package conversion

import (
	"regexp"
	"strings"
)

type SortableSliceOfVersions []string

func (a SortableSliceOfVersions) Len() int      { return len(a) }
func (a SortableSliceOfVersions) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortableSliceOfVersions) Less(i, j int) bool {
	return IsGreaterThan(a[i], a[j])
}

func IsGreaterThan(apiVer0, apiVer1 string) bool {
	// Remove group prefix
	idx := strings.IndexRune(apiVer0, '/')
	if idx >= 0 {
		apiVer0 = apiVer0[idx+1:]
	}
	idx = strings.IndexRune(apiVer1, '/')
	if idx >= 0 {
		apiVer1 = apiVer1[idx+1:]
	}

	vi, vj := strings.TrimLeft(apiVer0, "v"), strings.TrimLeft(apiVer1, "v")
	major := regexp.MustCompile("^[0-9]+")
	viMajor, vjMajor := major.FindString(vi), major.FindString(vj)
	viRemaining, vjRemaining := strings.TrimLeft(vi, viMajor), strings.TrimLeft(vj, vjMajor)
	switch {
	case viMajor != vjMajor:
		// more mature versions are greater
		return viMajor < vjMajor
	case len(viRemaining) == 0 && len(vjRemaining) == 0:
		return viMajor < vjMajor
	case len(viRemaining) == 0 && len(vjRemaining) != 0:
		// stable version is greater than unstable version
		return false
	case len(viRemaining) != 0 && len(vjRemaining) == 0:
		// stable version is greater than unstable version
		return true
	default:
		// neither are stable versions, compare remainings
		return viRemaining < vjRemaining
	}

	//	// assuming at most we have one alpha or one beta version, so if vi contains "alpha", it's the lesser one.
	//	return strings.Contains(viRemaining, "alpha")
}
