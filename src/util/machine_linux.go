package machine

import (
	"syscall"
)

var Name = func() string {
	var un syscall.Utsname
	syscall.Uname(&un)
	return charsToString(un.Nodename) + " (" + charsToString(un.Sysname) +
		" " + charsToString(un.Release) + ")"
}()

func charsToString(x [65]int8) string {
	n := ""
	for _,v := range x {
		if v == 0 { return n }
		n += string(byte(v))
	}
	return n
}
