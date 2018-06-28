upstream {{.Proto}}-{{.NodeResource}} {
    hash $remote_addr consistent;
    {{range .Site}}
    server {{.SiteIp}}:{{.SitePort}}; 
    {{end}}                                                
	keepalive 2000;
}
server {
    listen       80;                                                         
    server_name  {{.NodeResource}};                                               
    client_max_body_size 1024M;

    location / {
        proxy_pass http://{{.Proto}}-{{.NodeResource}}/;
        proxy_set_header Host $host:$server_port;
    }
}