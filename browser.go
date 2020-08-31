package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sirupsen/logrus"
)

type Browser struct {
	*rod.Browser
}

func NewBrowser() (*Browser, error) {
	browser := rod.New().ControlURL(options.ControlURL)
	if err := browser.Connect(); err != nil {
		return nil, err
	}
	return &Browser{browser}, nil
}

// TODO:
func (b *Browser) reloadPageByURL() error {
	b.DefaultViewport(&proto.EmulationSetDeviceMetricsOverride{Width: 1800, Height: 1200})
	for _, v := range b.MustPages() {
		url := v.MustEval("window.location.href").Result.String()
		logrus.WithField("url", url).Info("")
		// TODO: regexp
		// TODO: snapshot timestamp
		if strings.Contains(url, "https://crm.stage-dc.ru/pypo/crm/v1/topline") {
			reloadPage(v)
			fmt.Printf("%s\n", innerHTML(v))
		}
	}
	return nil
}

func reloadPage(p *rod.Page) {
	logrus.WithField("url", p.MustEval("window.location.href").Result.String()).Info("reload")
	wait := p.WaitRequestIdle(time.Second*5, []string{}, []string{})
	p.Eval("location.reload(true)")
	wait()
}

func innerHTML(p *rod.Page) string {
	return p.MustEval("document.documentElement.innerHTML").Result.String()
}
