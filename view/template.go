package view

import (
	"bytes"
	"github.com/charmbracelet/bubbles/table"
	"github.com/flyhope/kubetea/comm"
	"github.com/flyhope/kubetea/lang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"html/template"
	"sync"
	"time"
)

// 模板函数定义
var defaultTemplate = template.New("kubetea").Funcs(template.FuncMap{
	"CountByPodLabelKeyValue": countByPodLabelKeyValue,
	"PodPhaseView":            podPhaseView,
	"PodReadyView":            podReadyView,
	"ContainerStateView":      containerStateView,
	"BoolView":                boolView,
	"FormatTime":              formatTime,
	"txt":                     txt,
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

// txt get string by i18n
func txt(name string) string {
	return lang.Txt(&i18n.LocalizeConfig{
		MessageID: name,
	})
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

// GetBody 根据KEY获取一个模板
func (t *templateStore) GetBody(name comm.ConfigTemplateName) []*template.Template {
	storeField := "body-" + name
	t.lock.RLock()
	data, ok := t.store[storeField]
	t.lock.RUnlock()
	if !ok {
		var err error
		configTemplate := comm.ShowKubeteaConfig().Template[name]
		if data, err = NewTmplParseSlice(configTemplate.Body); err != nil {
			logrus.WithFields(logrus.Fields{"name": name}).Fatal(err)
		}
		t.lock.Lock()
		t.store[storeField] = data
		t.lock.Unlock()
	}

	return data
}

// GetColumn 根据KEY获取一个模板
func (t *templateStore) GetColumn(name comm.ConfigTemplateName) []*template.Template {
	storeField := "column-" + name
	t.lock.RLock()
	data, ok := t.store[storeField]
	t.lock.RUnlock()
	if !ok {
		configTemplate := comm.ShowKubeteaConfig().Template[name]

		data = make([]*template.Template, 0, len(configTemplate.Column))
		for _, column := range configTemplate.Column {
			newTmpl, errParse := NewTmplParse(column.Title)
			if errParse != nil {
				logrus.WithFields(logrus.Fields{"name": name, "column": column}).Fatal(errParse)
			}
			data = append(data, newTmpl)
		}
		t.lock.Lock()
		t.store[storeField] = data
		t.lock.Unlock()
	}

	return data
}

// TemplateRenderBody 根据名称批量渲染模板
func TemplateRenderBody(name comm.ConfigTemplateName, data any) []string {
	tmpls := TemplateStore.GetBody(name)
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

// TemplateRenderColumn render table column
func TemplateRenderColumn(name comm.ConfigTemplateName) []table.Column {
	tmpls := TemplateStore.GetColumn(name)

	columnConfig := comm.ShowKubeteaConfig().ShowTemplateColumn(name)
	result := make([]table.Column, 0, len(columnConfig))
	for idx, column := range columnConfig {
		tableColumn := table.Column{
			Width: column.Width,
		}
		var buf bytes.Buffer
		if errExecute := tmpls[idx].Execute(&buf, nil); errExecute != nil {
			logrus.Warnln(errExecute)
			tableColumn.Title = column.Title
		} else {
			tableColumn.Title = buf.String()
		}
		result = append(result, tableColumn)
	}

	return result
}
