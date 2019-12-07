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
		basePath    string
		licensePath string
	}{
		{"testdata/project1", "testdata/project1/LICENSE"},
		{"testdata/project2", "testdata/project2/license.txt"},
		{"testdata/project3", "testdata/project3/license"},
		{"testdata/project4", ""},
	}

	for _, td := range testdata {
		project := NewBasicProject(td.basePath)
		if project.licenseFile != td.licensePath {
			t.Errorf("license path did not match, wont %s, got %s", td.licensePath, project.licenseFile)
		}
	}
}
