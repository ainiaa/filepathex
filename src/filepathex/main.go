package main

import (
	. "./pathex"
	//"regexp"
	//"fmt"
)

func main() {

	rootPath := "D:\\www\\imageco\\wangcai_plateform\\source\\php\\home\\Home\\Lib\\Action\\"

	include := FileFilter{
		".php",
		"readme",
		".class.php",
		P_SUFFIX,
	}
	exclude := FileFilter{
		".php",
		"readme",
		"(?i)i[a-z]+n",
		P_NON,
	}

	nonDirecotyFilter := DirecotyFilter{
		make([]string, 0),
		P_NON,
	}
	GetFileList(rootPath, PathFilter{
		include, exclude, nonDirecotyFilter, nonDirecotyFilter,
	})
	writePath := "d:\\fileList.php"
	//	foreachFileList()
	WriteFileList(writePath, rootPath)
}
