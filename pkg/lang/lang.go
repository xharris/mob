package lang

import (
	"encoding/json"
	"log"
	"log/slog"
	"path"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func Init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile(path.Join("pkg", "lang", "en.json"))
}

func Get(id string) string {
	if bundle == nil {
		log.Fatal("i18n bundle is not initialized")
	}
	loc := i18n.NewLocalizer(bundle, "en")
	str, err := loc.LocalizeMessage(&i18n.Message{ID: id})
	if err != nil {
		slog.Warn("error while getting lang string", "err", err)
		// str = id
	}
	if str == " " {
		str = ""
	}
	return str
}
