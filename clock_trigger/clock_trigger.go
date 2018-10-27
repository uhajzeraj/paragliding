package clock_trigger

import (
	"net/http"
	"time"
)

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func triggerWebhook() {
	http.Get("https://imt2681-paragliding.herokuapp.com/paragliding/admin/api/webhook")
}

func main() {
	doEvery(20*time.Second, triggerWebhook)
}
