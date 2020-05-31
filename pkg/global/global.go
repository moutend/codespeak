package global

import (
	"os/user"
	"path/filepath"
)

var (
	CodespeakPath       string
	CodespeakBufferPath string
	CodespeakAudioPath  string
)

func init() {
	u, err := user.Current()

	if err != nil {
		panic(err)
	}

	CodespeakPath = filepath.Join(u.HomeDir, ".codespeak")
	CodespeakBufferPath = filepath.Join(CodespeakPath, "text")
	CodespeakAudioPath = filepath.Join(CodespeakPath, "sound")
}
