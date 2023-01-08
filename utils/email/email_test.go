package email

import (
	"context"
	"testing"
)

func TestSendEmail(t *testing.T) {
	e := NewEmail("本人", "568089002@qq.com", "mjqdlojunzvnbfii", "smtp.qq.com:587")
	err := e.Send(context.Background(), []string{"568089002@qq.com"}, "html", "主题", `
		<html>
		<body>
		<h3>
		"Test send to email"
		</h3>
		</body>
		</html>
		`)
	t.Log(*e)
	t.Log(err)
}
