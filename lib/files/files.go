package files

import "io/ioutil"

func ReadFile(path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(dat)
}
