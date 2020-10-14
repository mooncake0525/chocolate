package utils

import (
	"os"
	"path"
)

/*
@Author : VictorTu
@Software: GoLand
*/

type fileUtil struct {
}

var FileUtil fileUtil

func (this *fileUtil) CreateFile(fileUri string) (has bool, err error) {
	if has, err = this.IsFileExist(fileUri); err != nil {
		return
	}
	if has {
		return has, nil
	}
	dir := path.Dir(fileUri)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return false, err
	}
	if _, err = os.Create(fileUri); err != nil {
		return false, err
	}
	return false, nil
}

// 判断文件文件夹是否存在
func (this *fileUtil) IsFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	// 我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}
