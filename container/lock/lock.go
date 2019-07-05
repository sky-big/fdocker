package lock

import (
	"fmt"
	"os"
	"syscall"

	log "common/clog"
)

type Locker interface {
	Lock() error
	UnLock() error
}

type fileLocker struct {
	path string
	f    *os.File
}

func NewFileLocker(path string) (Locker, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Blog.Warningf("new file locker open file error : %v", err)
		return nil, err
	}

	return &fileLocker{
		path: path,
		f:    f,
	}, nil
}

func (l *fileLocker) Lock() error {
	err := syscall.Flock(int(l.f.Fd()), syscall.LOCK_EX)
	if err != nil {
		return fmt.Errorf("cannot flock directory %s - %s", l.path, err)
	}
	return nil
}

func (l *fileLocker) UnLock() error {
	defer l.f.Close()
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}
