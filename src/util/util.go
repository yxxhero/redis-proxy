package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

var readcmd = []string{"get"}

func Cmdanalysis(cmd string) int {
	result := In_array(cmd, readcmd)
	if result == true {
		return 0
	} else {
		return 1
	}
}

func In_array(element string, elementlist []string) bool {
	for _, v := range elementlist {
		if v == element {
			return true
		}
	}
	return false
}

func HostPortToAddress(host string, port uint16) string {
	return host + ":" + strconv.Itoa(int(port))
}

func UrlToHost(url string) string {
	return strings.Split(url, ":")[0]
}

func AbsolutePath(relpath string) string {
	absolutePath, err := filepath.Abs(relpath)
	if err != nil {
		log.Println("current path error:", err)
	}
	return absolutePath
}

func HomePath() string {
	return AbsolutePath(".")
}

func SliceIndex(slice interface{}, element interface{}) int {
	index := -1
	sv := reflect.ValueOf(slice)
	if sv.Kind() != reflect.Slice {
		return index
	}
	ev := reflect.ValueOf(element).Interface()
	length := sv.Len()
	for i := 0; i < length; i++ {
		iv := sv.Index(i).Interface()
		if reflect.DeepEqual(iv, ev) {
			index = i
			break
		}
	}
	return index
}

func Md5String(str string) string {
	hash := md5.New()
	io.WriteString(hash, str)
	return hex.EncodeToString(hash.Sum(nil))
}

func IP4ToInt(ip string) int {
	nums := strings.Split(ip, ".")
	sum := 0
	for i := 0; i < len(nums); i++ {
		n, _ := strconv.Atoi(nums[i])
		sum += n
		sum <<= 8
	}
	return sum >> 8
}
func CheckError(err error) {
        if err != nil {
                panic(err)
        }
}
