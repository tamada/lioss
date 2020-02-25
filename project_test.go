package lioss

import "testing"

func TestIsLicenseFile(t *testing.T) {
	testdata := []struct {
		giveName string
		wontFlag bool
	}{
		{"LICENSE", true},
		{"license.txt", true},
		{"License.html", true},
		{"LicenseAnalyzer.java", false},
	}

	for _, td := range testdata {
		if isLicenseFile(td.giveName) != td.wontFlag {
			t.Errorf("isLicenseFile(%s) wont %v, but %v", td.giveName, td.wontFlag, !td.wontFlag)
		}
	}
}

func TestFindLicenseFile(t *testing.T) {
	testdata := []struct {
		basePath     string
		licensePaths []string
	}{
		{"testdata/project1", []string{"LICENSE"}},
		{"testdata/project2", []string{"license.txt"}},
		{"testdata/project3", []string{"license", "subproject/license"}},
		{"testdata/project4", []string{}},
	}

	for _, td := range testdata {
		project := NewBasicProject(td.basePath)
		if len(project.LicenseIDs()) != len(td.licensePaths) {
			t.Errorf("%s: length of license path did not match, wont %d, got %d", td.basePath, len(td.licensePaths), len(project.LicenseIDs()))
		}
		for i, id := range project.LicenseIDs() {
			if td.licensePaths[i] != id {
				t.Errorf("license did not match: wont %s, got %s", td.licensePaths[i], id)
			}
		}
	}
}
