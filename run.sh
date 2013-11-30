justrun -v -delay=500ms -c 'go build && ./mapi-tigertonic-gorp config.json' -stdin < <(find config.json web api mapi.go -type "f")
