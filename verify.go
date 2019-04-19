package seed

import "os"

func Verify(source *VideoSource) (e error) {
	_, e = os.Stat(source.PosterPath)
	if os.IsNotExist(e) {
		return e
	}

	for _, value := range source.Files {
		_, e := os.Stat(value)
		if os.IsNotExist(e) {
			return e
		}
	}
	return nil
}
