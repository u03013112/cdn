upstream {{.Proto}}-{{.NodeResource}} {
	hash $remote_addr consistent;
    {{range .Site}}
    server {{.SiteIp}}:{{.SitePort}};
    {{end}}  
}

server {
	listen {{.NodeResource}} udp;
	proxy_pass {{.Proto}}-{{.NodeResource}};
}