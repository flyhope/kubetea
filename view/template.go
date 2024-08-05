package view

import (
	"bytes"
	"github.com/flyhope/kubetea/comm"
	"github.com/sirupsen/logrus"
	"html/template"
	"sync"
	"time"
)

// 模板函数定义
var defaultTemplate = template.New("kubetea").Funcs(template.FuncMap{
	"PodPhaseView":       podPhaseView,
	"PodReadyView":       podReadyView,
	"ContainerStateView": containerStateView,
	"BoolView":           boolView,
	"FormatTime":         formatTime,
})

// formatTime 格式化模板时间
func formatTime(t time.Time) string {
	return t.Format(time.DateTime)
}

// 获取字符串输出的Bool值
func boolView(val bool) string {
	if val {
		return "✔️"
	}
	return "❌️"
}

// NewTmplParse 创建模板
func NewTmplParse(text string) (*template.Template, error) {
	newTmpl, err := defaultTemplate.Clone()
	if err != nil {
		return nil, err
	}
	return newTmpl.Parse(text)
}

// NewTmplParseSlice 创建一组模板
func NewTmplParseSlice(texts []string) ([]*template.Template, error) {
	templates := make([]*template.Template, 0, len(texts))
	for _, text := range texts {
		newTmpl, err := NewTmplParse(text)
		if err != nil {
			return nil, err
		}
		templates = append(templates, newTmpl)
	}
	return templates, nil
}

type templateStore struct {
	lock  *sync.RWMutex
	store map[comm.ConfigTemplateName][]*template.Template
}

var TemplateStore = &templateStore{
	lock:  new(sync.RWMutex),
	store: make(map[comm.ConfigTemplateName][]*template.Template),
}

// Get 根据KEY获取一个模板
func (t *templateStore) Get(name comm.ConfigTemplateName) []*template.Template {
	t.lock.RLock()
	data, ok := t.store[name]
	t.lock.RUnlock()
	if !ok {
		var err error
		configTemplate := comm.ShowKubeteaConfig().Template[name]
		if data, err = NewTmplParseSlice(configTemplate.Body); err != nil {
			logrus.WithFields(logrus.Fields{"name": name}).Fatal(err)
		}
		t.lock.Lock()
		t.store[name] = data
		t.lock.Unlock()
	}

	return data
}

// TemplateRender 根据名称批量渲染模板
func TemplateRender(name comm.ConfigTemplateName, data any) []string {
	tmpls := TemplateStore.Get(name)
	result := make([]string, 0, len(tmpls))
	for _, tmpl := range tmpls {
		var buf bytes.Buffer
		if errExecute := tmpl.Execute(&buf, data); errExecute != nil {
			logrus.Warnln(errExecute)
		}
		result = append(result, buf.String())
	}
	return result
}
