justrun -v -delay=500ms -c 'go build && ./mapi-tigertonic-gorp config.json' -stdin < <(find config.json web api repository utils mapi.go -type "f")
