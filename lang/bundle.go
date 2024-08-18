package lang

import (
	"embed"
	"errors"
	"github.com/flyhope/kubetea/comm"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"sync"
)

//go:embed active.*.yaml
var localeFS embed.FS

// init bundle i18n function
var bundle = sync.OnceValue(func() *i18n.Bundle {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// load code all language yaml file
	dir, err := localeFS.ReadDir(".")
	if err != nil {
		logrus.Fatal(err)
	}
	for _, file := range dir {
		_, err := b.LoadMessageFileFS(localeFS, file.Name())
		if err != nil {
			logrus.WithFields(logrus.Fields{"file": file.Name()}).Fatal(err)
		}
	}

	// @todo load user all language yaml file

	return b
})

// DefaultLang 获取当前环境的语言
var DefaultLang = sync.OnceValue(func() string {
	lang := comm.ShowKubeteaConfig().Language
	if lang != "" {
		return lang
	}

	// 获取 LANG 环境变量
	lang = os.Getenv("LANG")

	// 如果 LANG 为空，尝试获取 LC_ALL
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}

	// 如果仍然为空，返回默认值
	if lang == "" {
		return ""
	}

	// 移除编码后缀（如 ".UTF-8"）
	if idx := strings.Index(lang, "."); idx != -1 {
		lang = lang[:idx]
	}

	// 将下划线替换为连字符（RFC 2616格式）
	lang = strings.Replace(lang, "_", "-", -1)

	return lang
})

// Txt Render text with i18n
func Txt(lc *i18n.LocalizeConfig) string {
	localizer := i18n.NewLocalizer(bundle(), DefaultLang())
	str, err := localizer.Localize(lc)

	var messageNotFoundErr *i18n.MessageNotFoundErr
	if errors.As(err, &messageNotFoundErr) {
		localizer = i18n.NewLocalizer(bundle(), "en")
		str, err = localizer.Localize(lc)
	}

	if err != nil {
		str = lc.DefaultMessage.ID
		logrus.WithFields(logrus.Fields{"id": lc.MessageID}).Errorln(err)
	}
	return str
}
