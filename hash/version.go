package hash

import "errors"

var ErrUnknownHashVersion = errors.New("unknown hash version")

const (
	version_unknown hashVersion = 0
	version_1       hashVersion = 1
)

type hashVersion byte

func (v hashVersion) isUnknown() bool {
    return byte(v) == byte(version_unknown)
}

func versionFactory(v byte) hashVersion {
	switch v {
    case byte(version_1):
		return version_1
	default:
		return version_unknown
	}
}
