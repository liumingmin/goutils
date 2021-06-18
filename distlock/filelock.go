package distlock

/*
import (
	"fmt"
	"os"
	"syscall"
)

type FileLock struct {
	dir   string
	f     *os.File
	block bool
}

func NewFileLock(dir string, block bool) *FileLock {
	return &FileLock{
		dir:   dir,
		block: block,
	}
}

func (l *FileLock) Lock() error {
	f, err := os.Open(l.dir)
	if err != nil {
		return err
	}
	l.f = f

	opt := syscall.LOCK_EX
	if !l.block {
		opt = opt | syscall.LOCK_NB
	}

	err = syscall.Flock(int(f.Fd()), opt)
	if err != nil {
		return fmt.Errorf("cannot flock directory %s - %s", l.dir, err)
	}
	return nil
}

func (l *FileLock) RLock() error {
	f, err := os.Open(l.dir)
	if err != nil {
		return err
	}
	l.f = f

	opt := syscall.LOCK_SH
	if !l.block {
		opt = opt | syscall.LOCK_NB
	}

	err = syscall.Flock(int(f.Fd()), opt)
	if err != nil {
		return fmt.Errorf("cannot flock directory %s - %s", l.dir, err)
	}
	return nil
}

func (l *FileLock) Unlock() error {
	defer l.f.Close()
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}
*/
