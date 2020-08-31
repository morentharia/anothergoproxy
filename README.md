## Another proxy
for MacOS -> https://github.com/natkuhn/Chrome-debug


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
--urlmatch ^.*yandex.*$'
# --chromedp ws://127.0.0.1:9222/devtools/browser/44a6d3d2-3ce3-47b3-872e-80222e729419 \
# swag i -g api.go; \



{find . -name "*.go"; find . -name "*.js"}  | grep -v vendor | entr -rc  bash -c '\
go run *.go \
--proxy-addr :1888 \
--upstream http://localhost:8080 \
--urlmatch ^.*firing.*$ \
--pagematch ^.*firing.*$ \
--chromedp ws://localhost:9222/devtools/browser/4a805b7a-f336-4aec-ad4e-ebf49b6cae69 \
'


find . -name "*.js" | grep -v vendor | entr -rc  bash -c 'date; go generate'

```


## ----------------------
```bash
http --follow -v GET http://localhost:3333/infoPages 
http -v POST http://localhost:3333/navigatePage url="https://public-firing-range.appspot.com/address/location.hash/documentwrite#kjkjddRRRRR" targetId=172F3118FBF10C4349DFA26E100E0CFF waitSec=0


http --follow -v POST http://localhost:3333/navigatePage url="https://public-firing-range.appspot.com/address/location.hash/documentwrite#kjkjddRRRRRddddddddddddddddd<dflkj>alert(1)</script>" targetId=172F3118FBF10C4349DFA26E100E0CFF waitSec=1 && 
bat -pp --color=always /tmp/output/page/page_public-firing-range.appspot.com___address__location.hash__documentwrite_172F3118FBF10C4349DFA26E100E0CFF_body.html | grep -C8 RRRRR
```

