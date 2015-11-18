/**
* Daemon for golang
 */

package daemon

import (
	"fmt"
	"os"
	"runtime"
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
}

//Daemon run a daemon process
func (ctx *Context) Daemon() (err error) {
	return
}

func Boot(lockFilePath, pidFilePath string, main MainFunc) {
	ctx := &Context{
		lockFile: lockFilePath,
		pidFile:  pidFilePath,
		main:     main,
	}

	ctx.Daemon()

	// flock
	lockFd, err := syscall.Open(lockFilePath, syscall.O_CREAT|syscall.O_WRONLY, syscall.S_IRUSR|syscall.S_IWUSR)
	if err != nil {
		panic(err)
	}

	err = syscall.Flock(lockFd, syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		fmt.Println("The Process is already runing ")
		os.Exit(-1)
	}
	println("lock success")
	// pid file
	pidFd, err := syscall.Open(pidFilePath, syscall.O_CREAT|syscall.O_WRONLY, syscall.S_IRUSR|syscall.S_IWUSR)
	if err != nil {
		panic(err)
	}
	_curPid := strconv.Itoa(os.Getpid())
	syscall.Write(pidFd, []byte(_curPid))
	syscall.Close(pidFd)
	println("pid success")
	// chdir
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(curDir)
	if err != nil {
		panic(err)
	}

	// umask
	//syscall.Umask(0x0000)

	// close
	/*
		syscall.Close(syscall.Stdin)
		syscall.Close(syscall.Stdout)
		syscall.Close(syscall.Stderr)

		// redict
		fd, err := syscall.Open("/dev/null", syscall.O_RDWR, 0)
		if err != nil {
			panic(err)
		}
		syscall.Dup2(fd, syscall.Stdin)
		syscall.Dup2(fd, syscall.Stdout)
		syscall.Dup2(fd, syscall.Stderr)

	*/
	// signal
	//signal.Ignore(syscall.SIGHUP, syscall.SIGINT)
	//go dealSignal(daemonCh)

	// setuid

	// setgid

	// do process
	println("boot success")
}
