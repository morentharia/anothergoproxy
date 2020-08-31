
```bash
go get -u github.com/morentharia/anothergoproxy
```


```bash
anothergoproxy \                                                                                                                                                                                                                                                                          <<<
--proxy-addr :1888 \
--upstream http://localhost:8080 \
--chromedp ws://127.0.0.1:9222/devtools/browser/44a6d3d2-3ce3-47b3-872e-80222e729419 \
--urlmatch '^.*crm.*$'
```

## Dev notes:

```bash
ls *.go | entr -rc  bash -c 'go run proxy.go --addr :1888 --upstream http://localhost:8080 --urlmatch ^.*crm.*$'

ls *.go | entr -rc  bash -c '\
swag i -g api.go; \
go run *.go \
--proxy-addr :1888 \
--upstream http://localhost:8080 \
--chromedp ws://127.0.0.1:9222/devtools/browser/44a6d3d2-3ce3-47b3-872e-80222e729419 \
--urlmatch ^.*crm.*$'


ls *.go | entr -rc  bash -c '\
go run *.go \
--proxy-addr :1888 \
--upstream http://localhost:8080 \
--urlmatch ^.*domclick.*$'
# --chromedp ws://127.0.0.1:9222/devtools/browser/44a6d3d2-3ce3-47b3-872e-80222e729419 \
# swag i -g api.go; \
```
