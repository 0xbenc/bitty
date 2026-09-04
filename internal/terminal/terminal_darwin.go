//go:build darwin

package terminal

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	tcgets     = 0x40487413
	tcsets     = 0x80487414
	tioCGWinsz = 0x40087468
)

type windowSize struct {
	row, col, xpixel, ypixel uint16
}

type State struct{ termios syscall.Termios }

func ioctl(fd, request, arg uintptr) syscall.Errno {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, request, arg)
	return errno
}

func IsTTY(file *os.File) bool {
	var termios syscall.Termios
	return ioctl(file.Fd(), tcgets, uintptr(unsafe.Pointer(&termios))) == 0
}

func MakeRaw(file *os.File) (State, error) {
	var old syscall.Termios
	if errno := ioctl(file.Fd(), tcgets, uintptr(unsafe.Pointer(&old))); errno != 0 {
		return State{}, fmt.Errorf("read terminal mode: %w", errno)
	}
	raw := old
	raw.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	raw.Oflag &^= syscall.OPOST
	raw.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	raw.Cflag &^= syscall.CSIZE | syscall.PARENB
	raw.Cflag |= syscall.CS8
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0
	if errno := ioctl(file.Fd(), tcsets, uintptr(unsafe.Pointer(&raw))); errno != 0 {
		return State{}, fmt.Errorf("set terminal mode: %w", errno)
	}
	return State{termios: old}, nil
}

func Restore(file *os.File, state State) error {
	if errno := ioctl(file.Fd(), tcsets, uintptr(unsafe.Pointer(&state.termios))); errno != 0 {
		return fmt.Errorf("restore terminal mode: %w", errno)
	}
	return nil
}

func Size(file *os.File) (int, int, error) {
	var ws windowSize
	if errno := ioctl(file.Fd(), tioCGWinsz, uintptr(unsafe.Pointer(&ws))); errno != 0 {
		return 0, 0, fmt.Errorf("read terminal size: %w", errno)
	}
	return int(ws.col), int(ws.row), nil
}
