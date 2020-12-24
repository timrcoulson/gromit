package data

import "io/ioutil"

func Set(key string, value string) error  {
	d1 := []byte(value)
	return ioutil.WriteFile("/etc/grommit/" + key, d1, 0644)
}

func Get(key string) (string, error)  {
	dat, err := ioutil.ReadFile("/etc/grommit/" + key)
	return string(dat), err
}
