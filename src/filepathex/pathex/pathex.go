package pathex

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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

var nonFileFilter = FileFilter{
	"",
	"",
	"",
	P_NON,
}

var nonDirecotyFilter = DirecotyFilter{
	make([]string, 0),
	P_NON,
}

var currPathFilter PathFilter

var fileList map[string][]string = make(map[string][]string)
var direcotyListDefault int = 5
var direcotyList []string = make([]string, direcotyListDefault)

func WalkFuncImpl(path string, info os.FileInfo, err error) error {
	if err != nil { //get error when walking directory
		fmt.Printf("%s:get panic %s", path, err.Error())
	} else {
		//filterResult := filterPath(path, info) //filter
		filterResult := true
		fmt.Println(path, ":", filterResult) // for feature filter path function
		if !info.IsDir() && filterResult {   //不是目录且包含指定后缀
			funcName := ReadSpecialFile(path)
			fileList[path] = funcName
		}
	}
	return err
}

func filterPath(path string, info os.FileInfo, pathFilter PathFilter) bool {
	result := true
	if !info.IsDir() {
		include := pathFilter.FileInclude
		exclude := pathFilter.FileExlude
		if exclude.Operate != P_NON {
			result = filterFileViaExclude(path, exclude, info)
		}
		if result && include.Operate != P_NON && include.Operate != P_ALL {
			result = filterFileViaInclude(path, include, info)
		}
	} else {
		include := pathFilter.DirctoryInclude
		exclude := pathFilter.DirctoryExclude
		if exclude.Operate != P_NON {
			result = filterDirecotyViaInclude(path, exclude, info)
		}
		if result && include.Operate != P_NON && include.Operate != P_ALL {
			result = filterDirecotyViaExclude(path, include, info)
		}
		result = true // todo for dirctory filter
	}

	fmt.Println("filterPath:", path, "=====", result)
	return result
}

var filterContainRegexp *regexp.Regexp = nil
var filterIgnoreRegexp *regexp.Regexp = nil

func initFileContainRegexp(operation interface{}) *regexp.Regexp {
	if filterContainRegexp == nil {
		if fileFilterOperation, ok := operation.(FileFilter); ok {
			fileOperation := fileFilterOperation.FileContain
			filterContainRegexp = regexp.MustCompile(fileOperation)
		} else if direcotyFilterOperation, ok := operation.(DirecotyFilter); ok {
			fmt.Println(direcotyFilterOperation)
			direcotyOperation := "" //todo
			filterContainRegexp = regexp.MustCompile(direcotyOperation)
		}
	}
	return filterContainRegexp
}
func initFileIgnoreRegexp(operation interface{}) {
	if filterIgnoreRegexp == nil {
		if filterOperation, ok := operation.(FileFilter); ok {
			fileOperation := filterOperation.FileContain
			fmt.Println("initFileIgnoreRegexp  fileOperation:", initFileIgnoreRegexp)
			filterIgnoreRegexp = regexp.MustCompile(fileOperation)
		}
	}
}

func filterFileViaOperation(path string, filterOperation FileFilter, filterRegexp *regexp.Regexp, info os.FileInfo) bool {
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

func filterFileViaInclude(path string, include FileFilter, info os.FileInfo) bool {
	if include.Operate == P_NON {
		return true
	}
	initFileContainRegexp(include)
	return filterFileViaOperation(path, include, filterContainRegexp, info)
}

func filterFileViaExclude(path string, exclude FileFilter, info os.FileInfo) bool {
	if exclude.Operate == P_NON {
		return false
	}
	initFileIgnoreRegexp(exclude)
	result := filterFileViaOperation(path, exclude, filterIgnoreRegexp, info)
	return !result
}

func filterDirecotyViaOperation(path string, filter DirecotyFilter, filterRegexp *regexp.Regexp, info os.FileInfo) bool {
	result := true
	operate := filter.Operate
	if operate != P_ALL {
		suffix := ""  //todo
		prefix := ""  //todo
		contain := "" //todo
		hasSuffix := false
		hasPrefix := false
		hasContain := false
		name := info.Name()
		if operate == P_SUFFIX || operate == P_SUFFIX_AND_PREFIX || operate == P_SUFFIX_OR_PREFIX {
			hasSuffix = strings.HasSuffix(name, suffix)
		}
		if operate == P_PREFIX || operate == P_SUFFIX_AND_PREFIX || operate == P_SUFFIX_OR_PREFIX {
			hasPrefix = strings.HasPrefix(name, prefix)
		}
		if operate == P_CONTAIN {
			hasContain = strings.Contains(name, contain)
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
			result = filterRegexp.Match([]byte(name))
		}
	}

	return result
}

func filterDirecotyViaInclude(path string, include DirecotyFilter, info os.FileInfo) bool {
	if include.Operate == P_NON {
		return true
	}
	initFileContainRegexp(include)
	return filterDirecotyViaOperation(path, include, filterContainRegexp, info)
}

func filterDirecotyViaExclude(path string, exclude DirecotyFilter, info os.FileInfo) bool {
	if exclude.Operate == P_NON {
		return false
	}
	initFileIgnoreRegexp(exclude)
	result := filterDirecotyViaOperation(path, exclude, filterIgnoreRegexp, info)
	return !result
}

func GetFileList(path string, pathFilter PathFilter) {
	currPathFilter = pathFilter
	Walk(path, pathFilter, WalkFuncImpl)
}

func Walk(root string, pathFilter PathFilter, walkFn filepath.WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	return walk(root, info, pathFilter, walkFn)
}

func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

// walk recursively descends path, calling w.
func walk(path string, info os.FileInfo, pathFilter PathFilter, walkFn filepath.WalkFunc) error {
	err := walkFn(path, info, nil)

	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	filterResult := filterPath(path, info, pathFilter) //filter

	if filterResult == false { //filtered
		return filepath.SkipDir
	}

	direcotyListLen := len(direcotyList)
	infoName := info.Name()
	fmt.Println("infoName === ", infoName)
	if direcotyListLen < direcotyListDefault {
		direcotyList[direcotyListLen] = infoName
	} else {
		direcotyList = append(direcotyList, infoName)
	}

	names, err := readDirNames(path)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := os.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			pathFilterResult := filterPath(path, fileInfo, pathFilter) //filter
			if pathFilterResult {
				err = walk(filename, fileInfo, pathFilter, walkFn)
			}

			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

func GetFileListViaStartWith(rootPath string, condition string) {
	pathFilter := PathFilter{
		FileFilter{
			condition,
			"",
			"",
			P_PREFIX,
		},
		nonFileFilter,
		nonDirecotyFilter,
		nonDirecotyFilter,
	}

	GetFileList(rootPath, pathFilter)
}

func GetFileListViaEndWith(rootPath string, condition string) {
	pathFilter := PathFilter{
		FileFilter{
			condition,
			"",
			"",
			P_SUFFIX,
		},
		nonFileFilter,
		nonDirecotyFilter,
		nonDirecotyFilter,
	}

	GetFileList(rootPath, pathFilter)
}

func GetFileListViaContain(rootPath string, condition string, isRegex bool) {
	var pathFilter PathFilter
	if isRegex {
		pathFilter = PathFilter{
			FileFilter{
				condition,
				"",
				"",
				P_CONTAIN_REGEXP,
			},
			nonFileFilter,
			nonDirecotyFilter,
			nonDirecotyFilter,
		}
	} else {
		pathFilter = PathFilter{
			FileFilter{
				condition,
				"",
				"",
				P_CONTAIN,
			},
			nonFileFilter,
			nonDirecotyFilter,
			nonDirecotyFilter,
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

func WriteFileList(fileName string, rootPath string) {
	var file *os.File
	var err1 error
	if IsFileExists(fileName) { //文件已存在
		file, err1 = os.OpenFile(fileName, os.O_WRONLY, 0777)
	} else {
		file, err1 = os.Create(fileName)
	}

	if err1 == nil { //获取文件成果
		currentContent := "<?php\r\n return array(\r\n"
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

		for kk, vv := range direcotyList {
			fmt.Println("dir === ", kk, "------", vv)
		}
	} else {
		fmt.Println("打开文件", fileName, "失败，error:", err1.Error())
	}
	defer file.Close()
}
