package lioss

import "testing"

func TestNewZipFile(t *testing.T) {
	testdata := []struct {
		path        string
		successFlag bool
	}{
		{"testdata/project3.jar", true},
		{"testdata/project4", true},
		{"testdata/liossdb.json", false},
	}

	for _, td := range testdata {
		project, err := NewProject(td.path)
		if (err == nil) != td.successFlag {
			t.Errorf("%s: success flag did not match, wont %v: %v", td.path, td.successFlag, err)
		}
		if project != nil {
			defer project.Close()
		}
	}
}

func TestNoLicenseIDs(t *testing.T) {
	project := &ZipProject{path: "testdata/project3"}
	ids := project.LicenseIDs()
	if len(ids) != 0 {
		t.Errorf("no zip file, wont 0, got %d", len(ids))
	}
}

func TestLicenseIDsOfZipProject(t *testing.T) {
	project, _ := NewProject("testdata/project3.jar")
	defer project.Close()
	ids := project.LicenseIDs()
	if len(ids) != 2 {
		t.Errorf("testdata/project3.jar: size of license ids did not match, wont %d, got %d", 2, len(ids))
	}
	for _, id := range ids {
		file, err := project.LicenseFile(id)
		if err != nil {
			t.Errorf("%s: license file open error: %s", id, err.Error())
		}
		if file != nil {
			defer file.Close()
		}
	}
	_, err := project.LicenseFile("not/existing/file")
	if err == nil {
		t.Errorf("%s: found, wont not found", "not/existing/file")
	}
}
