package main

import (
	. "./pathex"
	//"regexp"
	//"fmt"
)

func main() {
	rootPath := "E:\\www\\imageco\\wangcaio2o\\wangcai_plateform\\source\\php\\home\\Home\\Lib\\Action"

	include := FilterOperation{
		".php",
		"readme",
		"Notice",
		P_SUFFIX_OR_PREFIX,
	}
	exclude := FilterOperation{
		".php",
		"readme",
		"(?i)i[a-z]+n",
		P_CONTAIN_REGEXP,
	}
	GetFileList(rootPath, PathFilter{
		include, exclude,
	})
	//	writePath := "d:\\fileList.php"
	//	foreachFileList()
	//	writeFileList(writePath, rootPath)
}
