//go:build !(darwin || freebsd || illumos || linux || netbsd || openbsd || solaris)

package lookup

func execute(string, int) ([]byte, string, error) {
	return nil, "", errUnsupported
}
