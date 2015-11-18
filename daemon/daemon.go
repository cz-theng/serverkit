/**
* Daemon for golang
 */

package daemon

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"syscall"
)

type MainFunc func()

//Context is a Dameon Context
type Context struct {
	lockFile string
	pidFile  string
	process  *os.Process
	main     MainFunc
	pid      int
}

func (ctx *Context) fatherDo(procName string) (err error) {
	// flock
	lockFd, err := syscall.Open(ctx.lockFile, syscall.O_CREAT|syscall.O_WRONLY, syscall.S_IRUSR|syscall.S_IWUSR)
	if err != nil {
		return
	}
	syscall.Write(lockFd, []byte(procName))
	err = syscall.Flock(lockFd, syscall.LOCK_EX|syscall.LOCK_NB)
	//syscall.Close(lockFd)
	if err != nil {
		fmt.Println("The Process is already runing ")
		return
	}

	stdin, err := os.Open(os.DevNull)
	if err != nil {
		return err
	}
	stdout, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	stderr, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	procattr := os.ProcAttr{
		Dir:   cwd,
		Env:   append(os.Environ(), procName+"_daemon=true"),
		Files: []*os.File{stdin, stdout, stderr},
	}
	childProc, err := os.StartProcess(os.Args[0], os.Args, &procattr)
	if err != nil {
		return
	}

	ctx.pid = childProc.Pid
	// pid file
	pidFd, err := syscall.Open(ctx.pidFile, syscall.O_CREAT|syscall.O_WRONLY, syscall.S_IRUSR|syscall.S_IWUSR)
	if err != nil {
		return
	}
	curPid := strconv.Itoa(ctx.pid)
	syscall.Write(pidFd, []byte(curPid))
	syscall.Close(pidFd)

	err = childProc.Release()
	if err != nil {
		return
	}
	return
}

func (ctx *Context) childDo() {
	ctx.main()
}

//Daemon run a daemon process
func (ctx *Context) Daemon() (err error) {
	cmd := os.Args[0]
	procName := path.Base(cmd)
	if procName == "" {
		fmt.Println("procName is nil ")
		err = errors.New("ProcName is nil")
		return
	}
	isDaemon := os.Getenv(procName + "_daemon")
	if isDaemon == "true" {
		println("child")
		ctx.childDo()
		return
	}
	err = ctx.fatherDo(procName)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func Boot(lockFilePath, pidFilePath string, main MainFunc) {
	ctx := &Context{
		lockFile: lockFilePath,
		pidFile:  pidFilePath,
		main:     main,
	}

	ctx.Daemon()

}
