package rod_util

import (
	"context"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func DismissDialog(page *rod.Page) {
	page.EvalOnNewDocument(`setInterval(() => {
		window.alert = () => {}; 
		window.confirm = () => {};
		window.prompt = () => {};
		window.onbeforeunload = () => {};
	}, 1);`)
	_, handle := page.HandleDialog()
	go func() {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		tk := time.NewTicker(time.Second)
		for {
			select {
			case <-tk.C:
				handle(&proto.PageHandleJavaScriptDialog{
					Accept:     true,
					PromptText: "",
				})
			case <-ctx.Done():
				return
			}
		}
	}()
}
