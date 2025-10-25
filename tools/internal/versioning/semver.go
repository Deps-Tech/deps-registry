package versioning

import (
	"sort"

	"github.com/Masterminds/semver/v3"
)

func Compare(v1, v2 string) int {
	ver1, err1 := semver.NewVersion(v1)
	ver2, err2 := semver.NewVersion(v2)

	if err1 != nil && err2 != nil {
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	}

	if err1 != nil {
		return -1
	}
	if err2 != nil {
		return 1
	}

	return ver1.Compare(ver2)
}

func GetLatest(versions []string) string {
	if len(versions) == 0 {
		return ""
	}

	var validVersions []*semver.Version
	var invalidVersions []string

	for _, v := range versions {
		ver, err := semver.NewVersion(v)
		if err != nil {
			invalidVersions = append(invalidVersions, v)
		} else {
			validVersions = append(validVersions, ver)
		}
	}

	if len(validVersions) == 0 {
		sort.Strings(invalidVersions)
		return invalidVersions[len(invalidVersions)-1]
	}

	sort.Sort(semver.Collection(validVersions))
	return validVersions[len(validVersions)-1].String()
}

func Sort(versions []string) []string {
	if len(versions) == 0 {
		return versions
	}

	var validVersions []*semver.Version
	var invalidVersions []string
	versionMap := make(map[string]string)

	for _, v := range versions {
		ver, err := semver.NewVersion(v)
		if err != nil {
			invalidVersions = append(invalidVersions, v)
		} else {
			validVersions = append(validVersions, ver)
			versionMap[ver.String()] = v
		}
	}

	sort.Sort(semver.Collection(validVersions))
	sort.Strings(invalidVersions)

	result := make([]string, 0, len(versions))
	for _, v := range validVersions {
		if original, ok := versionMap[v.String()]; ok {
			result = append(result, original)
		} else {
			result = append(result, v.String())
		}
	}
	result = append(result, invalidVersions...)

	return result
}

func IsValid(version string) bool {
	_, err := semver.NewVersion(version)
	return err == nil
}

