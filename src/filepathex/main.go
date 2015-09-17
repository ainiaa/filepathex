package main

import (
	. "./pathex"
)

func main() {
	rootPath := "D:\\www\\imageco\\wangcai_plateform\\source\\php\\home\\Home\\Lib\\Action\\"
	//	writePath := "d:\\fileList.php"

	GetFileList(rootPath, PathFilter{
		".php",
		"readme",
		"Notice",
		FileOperation{P_SUFFIX_OR_PREFIX, P_NON},
	})
	//	foreachFileList()
	//	writeFileList(writePath, rootPath)
}
