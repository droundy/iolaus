package cook

import (
	"exec"
	"os"
)

func SetRaw() {
	exec.Run("/bin/stty", []string{"/bin/stty","raw"}, os.Environ(),
		exec.PassThrough, exec.PassThrough, exec.PassThrough)
}

func SetCooked() {
	exec.Run("/bin/stty", []string{"/bin/stty","cooked"}, os.Environ(),
		exec.PassThrough, exec.PassThrough, exec.PassThrough)
}
