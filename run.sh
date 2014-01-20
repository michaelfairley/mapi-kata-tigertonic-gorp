justrun -v -delay=500ms -c 'go build && ./mapi-kata-tigertonic-gorp config.json' -stdin < <(find config.json web api db utils mapi.go -type "f")
