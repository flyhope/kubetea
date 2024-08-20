package view

import "github.com/nicksnyder/go-i18n/v2/i18n"

var langUpdateTime = &i18n.Message{
	ID:    "data_update_time",
	Other: "Data update time: {{.UpdateTime}}",
}

var langTotalWithNumber = &i18n.Message{
	ID:    "total_with_number",
	Other: "Total: {{.number}}",
}
