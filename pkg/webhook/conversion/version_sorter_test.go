package conversion

import (
	"sort"
	"testing"
)

func Test_Versions_sorter(t *testing.T) {
	versions := []string{
		"v2beta1",
		"v1",
		"v1beta3",
		"v1alpha3",
		"stable.example.com/v1alpha4",
		"example.com/v2alpha3",
		"v1beta2",
		"v1alpha1",
	}

	expected := []string{
		"v1alpha1",
		"v1alpha3",
		"stable.example.com/v1alpha4",
		"v1beta2",
		"v1beta3",
		"v1",
		"example.com/v2alpha3",
		"v2beta1",
	}

	sort.Sort(SortableSliceOfVersions(versions))

	for i, ver := range versions {
		if ver != expected[i] {
			t.Fatalf("%d element should be %s, got %s", i, expected[i], ver)
		}
	}
}
