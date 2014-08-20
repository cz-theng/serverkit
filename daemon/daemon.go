/**
* Daemon for golang
*/

package daemon 

import (
	"fmt"
	"runtime"
	"strconv"
	"syscall"
	"os/signal"
	"os"
)

func dealSignal(c chan os.Signal){
	signo := <- c
	if signo == syscall.SIGHUP || signo == syscall.SIGINT {
		fmt.Println("sighub or sigint")
	}
}

func Boot(lockFilePath,pidFilePath string){
	var pid  uintptr
	var errno syscall.Errno
	var err error
	daemonCh := make(chan os.Signal, 1)

	// fork
	pid, ret2, errno := syscall.RawSyscall(syscall.SYS_FORK,0,0,0)
	if (errno != 0){
		panic(err)
	}
	
	if (runtime.GOOS == "darwin") && (ret2 == 1){
		pid = 0
	}

	if pid >0 {
		// parent just exit
		os.Exit(0)
	}

	// set sid
	_,err = syscall.Setsid()
	if err != nil {
		panic(err)
	}
	
	// fork twice
/* darwin can't fork twice
	pid, _, errno = syscall.RawSyscall(syscall.SYS_FORK,0,0,0)
	if (errno != 0){
		panic(err)
	}

	if (runtime.GOOS == "darwin") && (ret2 == 1){
		pid = 0
	}
	if pid >0 {
		// parent just exit
		os.Exit(0)
	}
	fmt.Print("Grand child")
*/
	// flock
	lockFd,err:= syscall.Open(lockFilePath,syscall.O_CREAT | syscall.O_WRONLY,syscall.S_IRUSR | syscall.S_IWUSR)
	if err != nil {
		panic(err)
	}

	err = syscall.Flock(lockFd,syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		fmt.Println("The Process is already runing ")
		os.Exit(-1)
	}

	// pid file
	pidFd,err:= syscall.Open(pidFilePath,syscall.O_CREAT | syscall.O_WRONLY,syscall.S_IRUSR | syscall.S_IWUSR)
	if err != nil {
		panic(err)
	}
	_curPid := strconv.Itoa(os.Getpid())
	syscall.Write(pidFd,[]byte(_curPid))
	syscall.Close(pidFd)
	// chdir
	curDir,err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err  = os.Chdir(curDir)
	if err != nil {
		panic(err)
	}
	
	// umask
	syscall.Umask(0x0000)

	// close 
	syscall.Close(syscall.Stdin)
	syscall.Close(syscall.Stdout)
	syscall.Close(syscall.Stderr)


	// redict
	fd,err := syscall.Open("/dev/null",syscall.O_RDWR,0)
	if err!= nil {
		panic(err)
	}
	syscall.Dup2(fd,syscall.Stdin)
	syscall.Dup2(fd,syscall.Stdout)
	syscall.Dup2(fd,syscall.Stderr)


	// signal
	signal.Notify(daemonCh,syscall.SIGHUP,syscall.SIGINT)
	go dealSignal(daemonCh)


	// setuid

	// setgid

	// do process

}
