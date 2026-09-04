//go:build linux

package terminal

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	tcgets     = 0x5401
	tcsets     = 0x5402
	tioCGWinsz = 0x5413
)

type windowSize struct {
	row, col, xpixel, ypixel uint16
}

type State struct{ termios syscall.Termios }

func IsTTY(file *os.File) bool {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, file.Fd(), tcgets, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return errno == 0
}

func MakeRaw(file *os.File) (State, error) {
	var old syscall.Termios
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, file.Fd(), tcgets, uintptr(unsafe.Pointer(&old)), 0, 0, 0); errno != 0 {
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
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, file.Fd(), tcsets, uintptr(unsafe.Pointer(&raw)), 0, 0, 0); errno != 0 {
		return State{}, fmt.Errorf("set terminal mode: %w", errno)
	}
	return State{termios: old}, nil
}

func Restore(file *os.File, state State) error {
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, file.Fd(), tcsets, uintptr(unsafe.Pointer(&state.termios)), 0, 0, 0); errno != 0 {
		return fmt.Errorf("restore terminal mode: %w", errno)
	}
	return nil
}

func Size(file *os.File) (int, int, error) {
	var ws windowSize
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, file.Fd(), tioCGWinsz, uintptr(unsafe.Pointer(&ws)), 0, 0, 0); errno != 0 {
		return 0, 0, fmt.Errorf("read terminal size: %w", errno)
	}
	return int(ws.col), int(ws.row), nil
}
