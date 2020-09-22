//go:generate go run ./generate
package main

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	// docs is generated by Swag CLI, you have to import it.

	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	// gin-swagger middleware

	"github.com/urfave/cli/v2"
	// rotlog "github.com/judwhite/logrjack"
)

var proxy *Proxy
var browser *Browser
var rotlog *logrus.Logger

type Options struct {
	ProxyAddr        string `json:"proxy_addr"`
	RestAddr         string `json:"rest_addr"`
	UpstreamProxyURL string `json:"upstream_proxy_url"`
	ControlURL       string `json:"control_url"`
	URLMatch         string `json:"urlmatch"`
	PageMatch        string `json:"pagematch"`
	Verbose          bool   `json:"verbose"`
	OutputPath       string `json:"output_path"`
}

var options Options

func (o Options) CachePath() string {
	return filepath.Join(options.OutputPath, "cache")
}

func (o Options) PagePath() string {
	return filepath.Join(options.OutputPath, "page")
}
func (o Options) LogsPath() string {
	return filepath.Join(options.OutputPath, "logs")
}

var flags []cli.Flag

func init() {
	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "proxy-addr",
			Value:       "127.0.0.1:8080",
			Usage:       "proxy listen address",
			Destination: &options.ProxyAddr,
		},
		&cli.StringFlag{
			Name:        "rest-addr",
			Value:       "http://localhost:3333",
			Usage:       "proxy listen address",
			Destination: &options.RestAddr,
		},
		&cli.StringFlag{
			Name:        "upstream",
			Value:       "",
			Usage:       "upstream HTTP Proxy URL (example: http://127.0.0.1:8080)",
			Destination: &options.UpstreamProxyURL,
		},
		&cli.StringFlag{
			Name:        "chromedp",
			Value:       "",
			Usage:       "chrome controlURL (example: ws://127.0.0.1:9222/devtools/browser/44a6d3d2-3ce3-47b3-872e-80222e729419)",
			Destination: &options.ControlURL,
		},
		&cli.StringFlag{
			Name:        "urlmatch",
			Value:       "^.*$",
			Usage:       "urls to trace (regexp pattern)",
			Destination: &options.URLMatch,
		},
		&cli.StringFlag{
			Name:        "pagematch",
			Value:       "^.*$",
			Usage:       "urls to trace (regexp pattern)",
			Destination: &options.PageMatch,
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Value:       false,
			Usage:       "setting verbose to true will log information on each request sent to the proxy",
			Destination: &options.Verbose,
		},
		&cli.StringFlag{
			Name:        "output path",
			Value:       "/tmp/output",
			Usage:       "path to flash and load cache values",
			Destination: &options.OutputPath,
		},
	}
}

func main() {
	app := &cli.App{
		Name:  "anothergoproxy",
		Flags: flags,
		Action: func(c *cli.Context) error {
			var err error

			logrus.Printf("Config: %s", pp.Sprint(options))

			for _, pathName := range []string{
				options.OutputPath, options.CachePath(), options.PagePath(), options.LogsPath(),
			} {
				if _, err = os.Stat(pathName); os.IsNotExist(err) {
					err = os.Mkdir(pathName, 0700)
					if err != nil {
						logrus.WithError(err).Errorf("mkdir(\"%s\")", pathName)
						return err
					}
				}
			}

			// DEBUG TODO:
			if true {
				if browser, err = NewBrowser(); err != nil {
					return err
				}
			}
			proxy, err = NewProxy()
			if err != nil {
				// logrus.Printf("%v", err)
				logrus.WithError(err).Errorf("NewProxy")
				return err
			}
			//TODO

			go http.ListenAndServe(options.ProxyAddr, proxy)

			api, err := NewApi()
			if err != nil {
				logrus.WithError(err).Errorf("NewProxy")
				return err
			}
			go func() {
				restURL, err := url.Parse(options.RestAddr)
				if err != nil {
					logrus.WithError(err).Errorf("parse RestAddr")
					return
				}
				if err := api.Run(restURL.Host); err != nil {
					logrus.WithError(err).Error("api.Run")
				}
			}()

			rotlog = logrus.New()
			rotlog.SetFormatter(&logrus.JSONFormatter{})
			rotlog.SetOutput(&lumberjack.Logger{
				Filename:   path.Join(options.LogsPath(), "events.log"),
				MaxSize:    1, // megabytes
				MaxBackups: 20,
				MaxAge:     28,   //days
				Compress:   true, // disabled by default
			})

			// TODO
			select {}
			// return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
