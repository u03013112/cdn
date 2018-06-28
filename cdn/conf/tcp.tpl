upstream {{.Proto}}-{{.NodeResource}} {
	hash $remote_addr consistent;
    {{range .Site}}
    server {{.SiteIp}}:{{.SitePort}} weight=5 max_fails=3 fail_timeout=30s;
    {{end}}  
}

server {
	listen {{.NodeResource}};
	proxy_connect_timeout 1s;
	proxy_timeout 3s;
	proxy_pass {{.Proto}}-{{.NodeResource}};
}