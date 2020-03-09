package lioss

import "testing"

func TestRemoveBasePath(t *testing.T) {
	testdata := []struct {
		giveBasePath string
		givePath     string
		wontResult   string
	}{
		{"./testdata/project1", "./testdata/project1/LICENSE", "LICENSE"},
		{"./testdata/project1", "testdata/project1/LICENSE", "LICENSE"},
		{"testdata/project1", "./testdata/project1/LICENSE", "LICENSE"},
		{"testdata/project3", "testdata/project3/subproject/license", "subproject/license"},
		{"testdata/project3/subproject", "testdata/project2", "project2"},
	}
	for _, td := range testdata {
		gotResult := removeBasePath(td.giveBasePath, td.givePath)
		if !IsSamePath(gotResult, td.wontResult) {
			t.Errorf("result of removeBasePath(%s, %s) did not match, wont %s, got %s", td.giveBasePath, td.givePath, td.wontResult, gotResult)
		}
	}
}

func TestIsLicenseFile(t *testing.T) {
	testdata := []struct {
		giveName string
		wontFlag bool
	}{
		{"LICENSE", true},
		{"license.txt", true},
		{"License.html", true},
		{"LicenseAnalyzer.java", false},
		{"SomeLicense", false},
	}

	for _, td := range testdata {
		if IsLicenseFile(td.giveName) != td.wontFlag {
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
		{"testdata/project3.jar", []string{"project3/license", "project3/subproject/license"}},
		{"LICENSE", []string{"LICENSE"}},
	}

	for _, td := range testdata {
		project, _ := NewProject(td.basePath)
		defer project.Close()
		if len(project.LicenseIDs()) != len(td.licensePaths) {
			t.Errorf("%s: length of license path did not match, wont %d, got %d", td.basePath, len(td.licensePaths), len(project.LicenseIDs()))
		}
		for i, id := range project.LicenseIDs() {
			if !IsSamePath(td.licensePaths[i], id) {
				t.Errorf("license did not match: wont %s, got %s", td.licensePaths[i], id)
			}
			file, err := project.LicenseFile(id)
			defer file.Close()
			if err != nil {
				t.Errorf("project.LicnseFile(%s) failed: %s", id, err.Error())
			}
		}
	}
}
