package main

import (
	. "./pathex"
)

func main() {
	rootPath := "D:\\www\\imageco\\wangcai_plateform\\source\\php\\home\\Home\\Lib\\Action\\"
	//	writePath := "d:\\fileList.php"
	container := FileOperation{
		".php",
		"readme",
		"Notice",
		P_SUFFIX_OR_PREFIX,
	}
	ignorer := FileOperation{
		".php",
		"readme",
		"Notice",
		P_NON,
	}
	GetFileList(rootPath, PathFilter{
		container, ignorer,
	})
	//	foreachFileList()
	//	writeFileList(writePath, rootPath)
}
