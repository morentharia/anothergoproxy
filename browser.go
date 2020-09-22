package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/k0kubun/pp"
	"github.com/morentharia/anothergoproxy/js"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Browser struct {
	*rod.Browser
	pageURLMatch *regexp.Regexp
}

func NewBrowser() (*Browser, error) {
	var err error
	b := &Browser{}
	logrus.WithField("controlURL", options.ControlURL).Info("connect to chromedp")

	b.Browser = rod.New().ControlURL(options.ControlURL)
	if err = b.Browser.Connect(); err != nil {
		return nil, errors.WithStack(err)
	}

	b.pageURLMatch, err = regexp.Compile(options.PageMatch)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, p := range b.MustPages() {
		u, err := b.pageURL(p)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		b.DefaultViewport(&proto.EmulationSetDeviceMetricsOverride{Width: 1800, Height: 1200})
		if b.pageURLMatch.MatchString(u.String()) {
			if _, err := p.EvalOnNewDocument(js.Bypass); err != nil {
				return nil, err
			}
			if _, err := p.EvalOnNewDocument(js.Underscore); err != nil {
				return nil, err
			}

			initJS := strings.ReplaceAll(js.Init, "{{{ANOTHERPROXY_API_URL}}}", options.RestAddr)
			if _, err := p.EvalOnNewDocument(initJS); err != nil {
				return nil, err
			}
			if err := b.reloadPage(p, 2); err != nil {
				return nil, err
			}
			pp.Println("--H-----------------------")
		}
	}
	go func() {
		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-ticker.C:
				for _, p := range b.MustPages() {
					u, err := b.pageURL(p)
					if err != nil {
						continue
					}
					if b.pageURLMatch.MatchString(u.String()) {
						if err := b.StorePage(p); err != nil {
							logrus.WithError(err).Error("store page")
						}
					}
				}
			}

		}
	}()

	return b, nil
}

// TODO:remove me
func (b *Browser) ReloadPageByURL() error {
	b.DefaultViewport(&proto.EmulationSetDeviceMetricsOverride{Width: 1800, Height: 1200})
	for _, p := range b.MustPages() {
		u, err := b.pageURL(p)
		if err != nil {
			return errors.WithStack(err)
		}
		if b.pageURLMatch.MatchString(u.String()) {
			logrus.WithField("url", u.String()).Info("reload")

			// b.reloadPage(p)
			b.StorePage(p)
		}
	}
	return nil
}

func (b *Browser) Navigate(targetID string, pageURL string, waitSec int) error {
	// TODO: b.PageFromTarget  !!
	for _, p := range b.MustPages() {
		if p.MustInfo().TargetID == proto.TargetTargetID(targetID) {
			wait := p.WaitRequestIdle(time.Second*time.Duration(waitSec), []string{}, []string{})
			if err := p.Navigate(pageURL); err != nil {
				return err
			}
			b.reloadPage(p, waitSec)
			wait()
			// html := innerHTML(p)
			// pp.Println(html)
			if err := b.StorePage(p); err != nil {
				return err
			}

			return nil
		}
	}
	return errors.Errorf("targetId == %s not exists", targetID)
}

func (b *Browser) PagesInfo() []*proto.TargetTargetInfo {
	res := make([]*proto.TargetTargetInfo, 0, 16)
	for _, p := range b.MustPages() {
		res = append(res, p.MustInfo())
	}
	return res
}

func (b *Browser) StorePage(p *rod.Page) error {
	html := innerHTML(p)

	u, err := b.pageURL(p)
	if err != nil {
		return errors.WithStack(err)
	}
	filename := b.PageMetaFilename(p, u)
	jsonBytes, err := json.MarshalIndent(
		struct {
			*url.URL
			PageURL      string
			TargetID     string
			BodyFilename string
		}{
			u,
			p.MustInfo().URL,
			string(p.MustInfo().TargetID),
			b.PageBodyFilename(p, u),
		},
		"", "  ",
	)
	err = ioutil.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		// logrus.WithError(err).Error("write file")
		return errors.WithStack(err)
	}
	logrus.WithField("page filename", filename).Info("write")

	filename = b.PageBodyFilename(p, u)
	err = ioutil.WriteFile(filename, []byte(html), 0644)
	if err != nil {
		// logrus.WithError(err).Error("write file")
		return errors.WithStack(err)
	}
	logrus.WithField("page filename", filename).Info("write")
	return nil
}

func (b *Browser) PageBodyFilename(p *rod.Page, u *url.URL) string {
	return filepath.Join(options.PagePath(), fmt.Sprintf(
		"page_%s_%s_%s_body.html",
		u.Hostname(),
		strings.ReplaceAll(u.Path, "/", "__"),
		p.MustInfo().TargetID,
	))
}

func (b *Browser) PageMetaFilename(p *rod.Page, u *url.URL) string {
	return filepath.Join(options.PagePath(), fmt.Sprintf(
		"page_%s_%s_%s_meta.json",
		u.Hostname(),
		strings.ReplaceAll(u.Path, "/", "_"),
		p.MustInfo().TargetID,
	))
}

func (b *Browser) pageURL(p *rod.Page) (*url.URL, error) {
	// urlStr := p.MustEval("window.location.href").Result.String()
	return url.Parse(p.MustInfo().URL)
}

//TODO: err
func (b *Browser) reloadPage(p *rod.Page, waitSec int) error {
	logrus.WithField("url", p.MustEval("window.location.href").Result.String()).Info("reload")
	wait := p.WaitRequestIdle(time.Second*time.Duration(waitSec), []string{}, []string{})
	_, err := p.Eval("location.reload(true)")
	if err != nil {
		logrus.WithError(err).Error("eval")
		return err
	}
	wait()
	return nil
}

func innerHTML(p *rod.Page) string {
	return p.MustEval("document.documentElement.innerHTML").Result.String()
}
