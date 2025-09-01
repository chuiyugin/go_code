package urltool

import (
	"errors"
	"net/url"
	"path"
)

func GetbasePath(lurl string) (string, error) {
	myUrl, err := url.Parse(lurl)
	if err != nil {
		return "", err
	}
	if len(myUrl.Host) == 0 {
		return "", errors.New("no host in lurl")
	}
	basePath := path.Base(myUrl.Path)
	return basePath, nil
}
