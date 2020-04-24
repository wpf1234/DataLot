package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func StrMd5(str string ) string{
	m := md5.New()
	_,err := io.WriteString(m,str)
	if err != nil {
		panic(err)
	}
	md5String := fmt.Sprintf("%x", md5.Sum(nil))
	return md5String
}
