package pathex

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	P_NON               = -1   //-1
	P_ALL               = iota //0
	P_SUFFIX                   //1
	P_PREFIX                   //2
	P_CONTAIN                  //3
	P_SUFFIX_OR_PREFIX         //4
	P_SUFFIX_AND_PREFIX        //5
	P_CONTAIN_REGEXP           //6
)

type (
	FileFilter struct {
		FilePrefix  string //file prefix
		FileContain string //file contain
		FileSuffix  string //file suffix
		Operate     int    // -1: non 0:all 1:FileSuffix 2:FilePrefix 3:FileContain 4:FileSuffix or FilePrefix 5: FileSuffix and FilePrefix
	}
	DirecotyFilter struct {
		DirecotyPath []string
		Operate      int
	}
	PathFilter struct {
		FileInclude     FileFilter
		FileExlude      FileFilter
		DirctoryInclude DirecotyFilter
		DirctoryExclude DirecotyFilter
	}
)

var currPathFilter PathFilter

var fileList map[string][]string = make(map[string][]string)

func WalkFuncImpl(path string, info os.FileInfo, err error) error {
	if err != nil { //get error when walking directory
		fmt.Printf("%s:get panic %s", path, err.Error())
	} else {
		filterResult := filterPath(path, info) //filter
		fmt.Println(path, ":", filterResult)   // for feture filter path function
		if !info.IsDir() && filterResult {     //不是目录且包含制定后缀
			fileNameList := ReadSpecialFile(path)
			fileList[path] = fileNameList
		}
	}
	return err
}

func filterPath(path string, info os.FileInfo) bool {
	result := true
	include := currPathFilter.FileInclude
	exclude := currPathFilter.FileExlude
	if exclude.Operate != P_NON {
		result = filterPathViaExclude(path, exclude, info)
	}
	if result && include.Operate != P_NON && include.Operate != P_ALL {
		result = filterPathViaInclude(path, include, info)
	}
	return result
}

var filterContainRegexp *regexp.Regexp = nil
var filterIgnoreRegexp *regexp.Regexp = nil

func initFileContainRegexp(filterOperation FileFilter) {
	if filterContainRegexp == nil {
		fileOperation := filterOperation.FileContain
		filterContainRegexp = regexp.MustCompile(fileOperation)
	}
}
func initFileIgnoreRegexp(filterOperation FileFilter) {
	if filterIgnoreRegexp == nil {
		fileOperation := filterOperation.FileContain
		fmt.Println("initFileIgnoreRegexp  fileOperation:", initFileIgnoreRegexp)
		filterIgnoreRegexp = regexp.MustCompile(fileOperation)
	}
}

func filterPathViaOperation(path string, filterOperation FileFilter, filterRegexp *regexp.Regexp, info os.FileInfo) bool {
	result := true
	operate := filterOperation.Operate
	if operate != P_ALL {
		FileSuffix := filterOperation.FileSuffix
		FilePrefix := filterOperation.FilePrefix
		FileContain := filterOperation.FileContain
		hasSuffix := false
		hasPrefix := false
		hasContain := false
		fileName := info.Name()
		if operate == P_SUFFIX || operate == P_SUFFIX_AND_PREFIX || operate == P_SUFFIX_OR_PREFIX {
			hasSuffix = strings.HasSuffix(fileName, FileSuffix)
		}
		if operate == P_PREFIX || operate == P_SUFFIX_AND_PREFIX || operate == P_SUFFIX_OR_PREFIX {
			hasPrefix = strings.HasPrefix(fileName, FilePrefix)
		}
		if operate == P_CONTAIN {
			hasContain = strings.Contains(fileName, FileContain)
		}

		if operate == P_SUFFIX {
			result = hasSuffix
		} else if operate == P_PREFIX {
			result = hasPrefix
		} else if operate == P_SUFFIX_AND_PREFIX {
			result = hasSuffix && hasPrefix
		} else if operate == P_SUFFIX_OR_PREFIX {
			result = hasSuffix || hasPrefix
		} else if operate == P_CONTAIN {
			result = hasContain
		} else if operate == P_CONTAIN_REGEXP {
			result = filterRegexp.Match([]byte(fileName))
		}
	}

	return result
}

func filterPathViaInclude(path string, include FileFilter, info os.FileInfo) bool {
	if include.Operate == P_NON {
		return true
	}
	initFileContainRegexp(include)
	return filterPathViaOperation(path, include, filterContainRegexp, info)
}

func filterPathViaExclude(path string, exclude FileFilter, info os.FileInfo) bool {
	if exclude.Operate == P_NON {
		return false
	}
	initFileIgnoreRegexp(exclude)
	result := filterPathViaOperation(path, exclude, filterIgnoreRegexp, info)
	return !result
}

func GetFileList(path string, pathFilter PathFilter) {
	currPathFilter = pathFilter
	filepath.Walk(path, WalkFuncImpl)
}

func GetFileListViaStartWith(rootPath string, condition string) {
	pathFilter := pathFilter{
		FilterOperation{
			condition,
			"",
			"",
			P_PREFIX,
		},
	}

	GetFileList(rootPath, pathFilter)
}

func GetFileListViaEndWith(rootPath string, condition string) {
	pathFilter := pathFilter{
		FilterOperation{
			condition,
			"",
			"",
			P_SUFFIX,
		},
	}

	GetFileList(rootPath, pathFilter)
}

func GetFileListViaContain(rootPath string, condition string, isRegex bool) {
	if isRegex {
		pathFilter := pathFilter{
			FilterOperation{
				condition,
				"",
				"",
				P_CONTAIN_REGEXP,
			},
		}
	} else {
		pathFilter := pathFilter{
			FilterOperation{
				condition,
				"",
				"",
				P_CONTAIN,
			},
		}

	}
	GetFileList(rootPath, pathFilter)
}

func ReadSpecialFile(path string) []string {
	lenCount := 5
	funcName := make([]string, lenCount)
	file, err := os.Open(path)
	functionFormat := "\\s+function\\s+(?P<funcName>\\w+)\\s*\\("
	if err == nil {
		bufferReader := bufio.NewReader(file)
		i := 0
		for {
			content, _, err := bufferReader.ReadLine()
			if err == nil { //读取成功
				reg := regexp.MustCompile(functionFormat)
				c := string(content)
				s := reg.FindStringSubmatch(c)

				if s != nil && len(s) == 2 {
					if i < lenCount {
						funcName[i] = s[1]
					} else {
						funcName = append(funcName, s[1])
					}
					i++
					//fmt.Printf("%d:%s 找到function %s\r\n", i, path, s[1])
				}
			} else if err == io.EOF { //读取到文件结尾
				break
			} else { //读取的时候出现错误
				fmt.Printf("%d:%s get error:%s\r\n", i, path, err.Error())
			}

		}
	} else {
		fmt.Println("读取文件%s 失败:%s", path, err.Error())
	}

	defer file.Close()
	return funcName
}

func foreachFileList() {
	for k, v := range fileList {
		fmt.Println(k, ":", v)
		for kk, vv := range v {
			fmt.Println(k, ":", kk, ":", vv)
		}
		break
	}
}

// check is file exists
func IsFileExists(fileName string) (isExists bool) {
	isExists = true
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		isExists = false
	}
	return isExists
}

func writeFileList(fileName string, rootPath string) {
	var file *os.File
	var err1 error
	if IsFileExists(fileName) { //文件已存在
		file, err1 = os.OpenFile(fileName, os.O_WRONLY, 0777)
	} else {
		file, err1 = os.Create(fileName)
	}

	if err1 == nil { //获取文件成果
		currentContent := "<?php\r\n $fileList = array(\r\n"
		//fmt.Println("rootPath:",rootPath)
		suffix := ".class.php"
		for k, v := range fileList {
			fileKey := strings.Replace(k, rootPath, "", -1)
			fileKey = strings.TrimSuffix(fileKey, suffix)
			currentContent += "    '" + fileKey + "' => array(\r\n"
			for _, vv := range v {
				if strings.Trim(vv, "") != "" {
					currentContent += "        '" + vv + "' => 1,\r\n"
				}
			}
			currentContent += "\r\n    ),\r\n"
		}
		currentContent += ");"
		file.WriteString(currentContent + "\r\n")

	} else {
		fmt.Println("打开文件", fileName, "失败，error:", err1.Error())
	}
	defer file.Close()
}
