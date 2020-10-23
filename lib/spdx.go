package lib

import (
	"fmt"
	"os"
	"strings"

	"github.com/denisbrodbeck/striphtmltags"
	"gopkg.in/xmlpath.v1"
)

/*
LicenseMeta shows meta information of license of SPDX.
*/
type LicenseMeta struct {
	Names       *Names
	OsiApproved bool
	Deprecated  bool
	Urls        []string
}

/*
Names shows names of license (short and full name)
*/
type Names struct {
	ShortName string
	FullName  string
}

func (lm *LicenseMeta) String() string {
	return fmt.Sprintf("%s (%s), OSI: %v, Deprecated: %v", lm.Names.ShortName, lm.Names.FullName, lm.OsiApproved, lm.Deprecated)
}

func createMeta(root *xmlpath.Node) *LicenseMeta {
	meta := new(LicenseMeta)
	meta.OsiApproved = isTrue(root, "/SPDXLicenseCollection/license/@isOsiApproved")
	meta.Deprecated = isTrue(root, "/SPDXLicenseCollection/license/@isDeprecated")
	meta.Urls = stringSlice(root, "/SPDXLicenseCollection/license/crossRefs/crossRef")
	meta.Names = new(Names)
	meta.Names.ShortName = findString(root, "/SPDXLicenseCollection/license/@licenseId")
	meta.Names.FullName = findString(root, "/SPDXLicenseCollection/license/@name")
	// fmt.Printf("\"%s\",\"%s\",%v,%v,%v\n", meta.Names.ShortName, meta.Names.FullName, meta.OsiApproved, meta.Deprecated, strings.Join(meta.Urls, ","))
	return meta
}

/*
ReadSPDX reads license data from SPDX xml file and returns meta information of license and lisense terms.
*/
func ReadSPDX(path string) (*LicenseMeta, string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer reader.Close()
	root, err := xmlpath.Parse(reader)
	if err != nil {
		return nil, "", err
	}
	meta := createMeta(root)
	data := stripHTML(findString(root, "/SPDXLicenseCollection/license/text"))

	return meta, Normalize([]byte(data)), nil
}

func stringSlice(root *xmlpath.Node, xpath string) []string {
	path := xmlpath.MustCompile(xpath)
	results := []string{}
	iter := path.Iter(root)
	for iter.Next() {
		results = append(results, strings.TrimSpace(iter.Node().String()))
	}
	return results
}

func isTrue(root *xmlpath.Node, xpath string) bool {
	val := findString(root, xpath)
	if strings.ToLower(val) == "true" {
		return true
	}
	return false
}

func findString(root *xmlpath.Node, xpath string) string {
	path := xmlpath.MustCompile(xpath)
	str, _ := path.String(root)
	return strings.TrimSpace(str)
}

func stripHTML(text string) string {
	return striphtmltags.StripTags(text)
}

/*
Normalize normalizes given string (remove return code, tab, and adjacent spaces)
*/
func Normalize(dataArray []byte) string {
	data := strings.ReplaceAll(string(dataArray), "\r", " ")
	data = strings.ReplaceAll(data, "\n", " ")
	data = strings.ReplaceAll(data, "\t", " ")
	for strings.Index(data, "  ") >= 0 {
		data = strings.ReplaceAll(data, "  ", " ")
	}
	return strings.TrimSpace(data)
}
