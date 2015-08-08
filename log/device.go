package log

import (
	"fmt"
	"os"
	"sync"
	"bytes"
	"path/filepath"
)


type Devicer interface {
	Write(buf []byte)(n int,err error)
}

type Device struct {
	fp *os.File
	mtx sync.Mutex
}

type ConsoleDevice struct {
	Device
}

func NewConsoleDevice() (*ConsoleDevice,error){
	cd := new(ConsoleDevice)
	cd.fp = os.Stdout
	return cd,nil
}

func (cd *ConsoleDevice) Write(buf []byte) (n int, err error) {
	cd.mtx.Lock()
	defer cd.mtx.Unlock()
	n,err = cd.fp.Write(buf)
	return 
}

type FileDevice struct {
	Device
	fileName string
	fileSize uint64
	logLen   uint64
}

func NewFileDevice(fileName string) (fd *FileDevice, err error){ 
	fd = new(FileDevice)
	fd.fileName = fileName
	fd.fp, err = openFile(fd.fileName)
	return fd,err
}

func (fd *FileDevice) SetFileSize(size uint64) {
	fd.fileSize = size
}

func (fd *FileDevice) SetFileName(fileName string) {
	fd.fileName = fileName
}


func (fd *FileDevice) Write(buf []byte) (n int, err error) {
	bufLen := uint64(len(buf))
	if bufLen+fd.logLen <= fd.fileSize {
		n,err = fd.fp.Write(buf)
		fd.logLen += uint64(n)
		return
	} else {
		remainBuf := buf[fd.fileSize-fd.logLen:]
		n,err = fd.fp.Write(buf[:fd.fileSize-fd.logLen])
		if err != nil {
			return
		}
		fd.fp.Sync()
		fd.fp.Close()
		fd.fp,err = openFile(fd.fileName)
		if err != nil {
			return 
		}
		fd.logLen = 0
		n,err = fd.fp.Write(remainBuf)
		if err != nil {
			return
		}
		fd.logLen += uint64(n)
		return
	}
}

func fileNotExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil{
		if (os.IsNotExist(err)) {
			return true
		}
		fmt.Println("stat error:",err)
		return false
	}
	return false
}

func openFile(fileName string) (fp *os.File, err error) {
	err = os.MkdirAll(filepath.Dir(fileName),os.ModePerm)
	if err != nil {
		return 
	}
	logPath := fileName
	for i:=1;;i++{
		ret := fileNotExist(logPath)
		if ret  {
			flag := os.O_WRONLY|os.O_CREATE|os.O_APPEND
			fp,err = os.OpenFile(logPath, flag, 0666)
		    break
		} else {
			var buf bytes.Buffer
			buf.WriteString(fmt.Sprintf("%s.%d",fileName,i))
			logPath = buf.String()
		}
	}
	return
}


