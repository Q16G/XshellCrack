package ssh

import (
	"errors"
	"gorun/config"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

var xshFile []string

func checkpath(path string) (string, error) { //对数据进行清理
	r := regexp.MustCompile("HKEY_CURRENT_USER")
	s := r.MatchString(path)
	if s {
		return path, nil
	} else {
		return "1", errors.New("not find")
	}
}
func findpath(path string) []string { //找到路径
	c := exec.Command("cmd.exe", "/C", "reg query "+path+"\\ProfileDialog /v Recent")
	b, err := c.CombinedOutput()
	if err != nil {
		err2 := errors.New("no find xshell path")
		config.Checkerror(err2)
	}
	b2, _ := simplifiedchinese.GBK.NewDecoder().Bytes(b)
	r := regexp.MustCompile(`(  )+|()`) //去除空格还有问题，待完善
	s := r.ReplaceAllString(string(b2), "")
	s2 := strings.Split(s, "\n")
	return s2
}
func findsession(path string) (string, error) { //定位session路径
	b := strings.Contains(path, "RecentREG_SZC:")
	if b {
		return path, nil
	}
	return "", errors.New("not find")
}
func checkXshell() []string { //检测是否存在xshell
	c := exec.Command("cmd.exe", "/C", "reg query HKEY_CURRENT_USER\\SOFTWARE\\NetSarang\\Xshell")
	b, err := c.CombinedOutput()
	if err != nil {
		err2 := errors.New("no find xshell path")
		config.Checkerror(err2)
	}
	b2, _ := simplifiedchinese.GBK.NewDecoder().Bytes(b)
	r := regexp.MustCompile(`( )+|()`)
	s := r.ReplaceAllString(string(b2), "")
	s2 := strings.Split(s, "\n")
	return s2
}
func retupath(path string) string {
	r := regexp.MustCompile(`RecentREG_SZ(.*)\\`)
	s := r.FindStringSubmatch(path)
	return s[1]
}
func GetUserSid() (localname, SID string) { //获取到本机的localname和SID
	c := exec.Command("cmd.exe", "/C", "whoami /user /fo csv /nh")
	b, _ := c.CombinedOutput()
	b2, _ := simplifiedchinese.GBK.NewDecoder().Bytes(b)
	r := regexp.MustCompile(`"|\n`)
	s := r.ReplaceAllString(string(b2), "")
	s2 := strings.Split(s, ",")
	// fmt.Printf("s2: %v\n", s2[1])
	i := strings.Index(s2[0], "\\")
	return s2[0][i+1:], s2[1]
}
func GetFileName(pathstr string) []string {
	AllFileList, err := ioutil.ReadDir(pathstr)
	if err != nil {
		config.Checkerror(err)
	}
	for i := range AllFileList {
		FileExt := path.Ext(AllFileList[i].Name())
		if FileExt == ".xsh" {
			xshFile = append(xshFile, pathstr+"\\"+AllFileList[i].Name())
		}
	}
	return xshFile
}
func ReplaceStringByRegex(str, rule, replace string) (string, error) {
	reg, err := regexp.Compile(rule)
	if reg == nil || err != nil {
		return "", errors.New("正则MustCompile错误:" + err.Error())
	}
	return reg.ReplaceAllString(str, replace), nil
}
func getHost_name_encryptpw(str string) (map[string]string, error) {
	i := 0
	var strmap = make(map[string]string, 1024)
	r := regexp.MustCompile("Host=(.*)|Password=(.*)|UserName=(.*)|Version=(.*)")
	s := r.FindAllString(str, -1)
	// fmt.Printf("s: %v\n", s)
	for _, v := range s {
		i++
		s2, _ := ReplaceStringByRegex(v, "Host=|Password=|UserName=|Version=", "")
		if s2 == "\r" {
			return strmap, errors.New("not found")
			break
		}
		if i == 1 {
			strmap["version"] = s2
		}
		if i == 2 {
			strmap["ip"] = s2
		}
		if i == 3 {
			strmap["encryptpw"] = s2
		}
		if i == 4 {
			strmap["name"] = s2
		}
	}
	// fmt.Printf("strmap: %v\n", strmap)
	return strmap, nil
}
func conver(xshbyte []byte) string {
	var data []byte
	for i := 0; i < len(xshbyte); i++ {
		if xshbyte[i] == 0 {
			continue
		}
		data = append(data, xshbyte[i])
	}
	return string(data)
}
func GetxshFile_context(str string) (string, error) { //如果只是传入FIle路径,就执行getcontext
	f, err := os.Open(str)
	if err != nil {
		return "", errors.New("not find")
	}
	b, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		return "", errors.New("not find")
	}
	s := conver(b)
	return s, nil
}
func Run() []map[string]string {
	var resultmap []map[string]string //切片中存放着map
	var path []string
	s := checkXshell()
	for _, v := range s {
		s2, err := checkpath(v)
		if err != nil {
			continue
		}
		s3 := findpath(s2)
		for _, v := range s3 {
			s4, err2 := findsession(v)
			if err2 == nil {
				s5 := retupath(s4)
				GetFileName(s5)
				path = append(path, s5)
			}
		}
	}
	// fmt.Printf("xshFile: %v\n", xshFile)
	for _, v := range xshFile {
		// fmt.Printf("v: %v\n", v)
		s2, err := GetxshFile_context(v)
		if err == nil {
			m, err2 := getHost_name_encryptpw(s2)
			if err2 == nil {
				// fmt.Printf("m: %v\n", m)
				resultmap = append(resultmap, m)
			}
		}
	}
	return resultmap
}
func RunFile(file string) {
	s, err := GetxshFile_context(file)
	if err != nil {
		config.Errlog_print(err.Error())
	}
	m, err2 := getHost_name_encryptpw(s)
	if err2 != nil {
		config.Errlog_print(err2.Error())
	}
	InitXSh(m)
}
