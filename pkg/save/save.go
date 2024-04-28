package save

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Save struct {
	GameName  string
	ProfileID string
}

func (s *Save) getPath() string {
	homeDir, _ := os.UserCacheDir()
	// TODO handle error and figure out a backup directory?
	p := path.Join(homeDir, s.GameName)
	os.MkdirAll(p, os.ModePerm)
	return path.Join(fmt.Sprintf("%s.json", s.getProfileID()))
}

func (s *Save) getProfileID() string {
	if len(s.ProfileID) <= 0 {
		return "local"
	}
	return s.ProfileID
}

func (s *Save) Write(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(s.getPath(), b, os.ModePerm)
}
