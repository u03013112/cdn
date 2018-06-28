package cdn

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	// "log"
	"os"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

// func (node *Node) GlobalResourcePath() string {
// 	return apiConf.GetString("cdn_resource_dir")
// }

// func (node *Node) createDefaultNginxConf() error {
// 	default_ng := `
// server {
//     listen       80 default_server;
//     server_name  _;
//     access_log off;
//     return 301 https://$host$request_uri;
// }
// `

// 	f, err := os.OpenFile(fmt.Sprintf("%s/default.conf", node.NginxHttpConfPath()), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
// 	if err != nil {

// 		return err
// 	}
// 	defer f.Close()

// 	f.WriteString(default_ng)
// 	return nil
// }

func (node *Node) Dump(w io.Writer) error {
	byt, _ := yaml.Marshal(node)

	w.Write(byt)
	return nil
}

func (node *Node) reload() error {
	invoke := &Invoke{}
	if res, err := invoke.Command("docker", "exec", "-d", node.NodeName, "nginx", "-s", "reload"); err != nil {
		log.Errorf("docker exec error: %s\n", string(res))
		return err
	}
	return nil
}
func (node *Node) dumpRuleTo(r *StreamRule, fname string) error {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	r.Dump(f)

	return nil
}

func (node *Node) dumpRule(r *StreamRule) error {

	node.dumpRuleTo(r, fmt.Sprintf("%s/%s", node.NodeResourcePath(), r.Index()))
	return nil
}

func (node *Node) removeRule(r *StreamRule) error {

	os.RemoveAll(fmt.Sprintf("%s/%s", node.NodeResourcePath(), r.Index()))
	return nil
}

func (node *Node) AddRule(r *StreamRule) error {

	if err := r.Create(node.NginxConfPath()); err != nil {
		return err
	}
	if err := node.reload(); err != nil {
		return err
	}
	if err := node.dumpRule(r); err != nil {
		return err
	}
	return nil
}

func (node *Node) DelRule(r *StreamRule) error {

	if err := r.Destory(node.NginxConfPath()); err != nil {
		return err
	}
	if err := node.reload(); err != nil {
		return err
	}

	if err := node.removeRule(r); err != nil {
		return err
	}
	return nil
}

type CdnNet struct {
	NodeIp  string `json:"node_ip"`
	Mask    string `json:"mask"`
	Gateway string `json:"gateway"`
}

// 10.0.0.2/24@10.0.0.1
// ip/mask@gateway
func (node *CdnNet) Parse(net string) error {
	return nil
}

func (node *CdnNet) String() string {
	return fmt.Sprintf("%s/%s@%s", node.NodeIp, node.Mask, node.Gateway)
}

type SiteInfo struct {
	SiteIp   string `json:"site_ip"`
	SitePort string `json:"site_port"`
}
type StreamRule struct {
	NodeName     string      `json:"node_id"`
	StreamId     string      `json:"stream_id"` //default : [NodeName_Proto_NodeResource]
	StreamName   string      `json:"stream_name"`
	Proto        string      `json:"proto"` //http, tcp, udp
	NodeResource string      `json:"node_resource"`
	Site         []*SiteInfo `json:"site"`
	// SiteIp       string `json:"site_ip"`
	// SitePort     string `json:"site_port"`
}

func (node *StreamRule) GetName() string {
	return fmt.Sprintf("%s-%s", node.Proto, node.NodeResource)
}

func (node *StreamRule) readTpl(tplfilename string) string {

	f, err := os.Open(tplfilename)
	if err != nil {
		return ""
	}

	defer f.Close()
	byt, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}
	return string(byt)
}

func (node *StreamRule) GetTpl() string {
	var tplfeild string
	switch node.Proto {
	case "http":
		tplfeild = "http_tpl"
	case "tcp":
		tplfeild = "tcp_tpl"
	case "udp":
		tplfeild = "udp_tpl"
	}
	tplfilename := apiConf.GetString(tplfeild)
	return node.readTpl(tplfilename)
}

func (node *StreamRule) getConfPathName(dir string) string {
	return fmt.Sprintf("%s/%s.d/%s.conf", dir, node.Proto, node.NodeResource)
}

func (node *StreamRule) Create(dir string) error {
	tpl := node.GetTpl()
	fmt.Printf("%+v\n", node)
	t, err := template.New(node.GetName()).Parse(tpl)
	if err != nil {
		return err
	}

	fname := node.getConfPathName(dir)
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, node)
}

func (node *StreamRule) Destory(dir string) error {
	fname := node.getConfPathName(dir)
	os.Remove(fname)
	return nil
}

func (node *StreamRule) Load(r io.Reader) error {
	byt, _ := ioutil.ReadAll(r)
	return yaml.Unmarshal(byt, node)
}

func (node *StreamRule) Dump(w io.Writer) error {
	byt, _ := yaml.Marshal(node)
	w.Write(byt)
	return nil
}

func (node *StreamRule) Index() string {
	if node.StreamId != "" {
		return node.StreamId
	}

	return fmt.Sprintf("%s_%s_%s", node.NodeName, node.Proto, node.NodeResource)
}

/*******************************************************************
*
* Network Here!
*
 *******************************************************************/

//Network :
type Network struct {
	NetworkName string `json:"networkname"`
	Interface   string `json:"interface"` //eth0 or eth0.100
	NetworkIP   string `json:"networkip"`
	Mask        string `json:"mask"` //http, tcp, udp
	Gateway     string `json:"gateway"`
	Vlan        int    `json:"vlan"` //vlan Num 2~4096
}

//CreateNetwork : create a docker network 4 docker run --net=network
func (network *Network) CreateNetwork() error {
	invoke := &Invoke{}

	var interfaceName string
	if network.Vlan == 0 {
		interfaceName = fmt.Sprintf("parent=%s", network.Interface)
	} else {
		interfaceName = fmt.Sprintf("parent=%s.%d", network.Interface, network.Vlan)
	}
	args := []string{
		"network",
		"create",
		"-d",
		"macvlan",
		fmt.Sprintf("--subnet=%s/%s", network.NetworkIP, network.Mask),
		fmt.Sprintf("--gateway=%s", network.Gateway),
		"-o",
		interfaceName,
		fmt.Sprintf("%s", network.NetworkName),
	}

	if byt, err := invoke.Command("docker", args...); err != nil {
		log.Errorf("cmd: docker %s, err: %s", strings.Join(args, " "), err.Error())
		log.Errorf("%s", string(byt))
		return err
	}

	return nil
}

//DestoryNetwork : 暂时没有检查引用，全靠docker报错
func (network *Network) DestoryNetwork() error {
	invoke := &Invoke{}

	args := []string{
		"network",
		"rm",
		fmt.Sprintf("%s", network.NetworkName),
	}

	if byt, err := invoke.Command("docker", args...); err != nil {
		log.Errorf("cmd: docker %s, err: %s", strings.Join(args, " "), err.Error())
		log.Errorf("%s", string(byt))
		return err
	}
	return nil
}

//NetworkList :传输Network的列表，以一定填满
type NetworkList []*Network

//ShowNetworkAll :
func (network *Network) ShowNetworkAll() (string, error) {
	invoke := &Invoke{}

	args := []string{
		"network",
		"ls",
	}

	var err error
	var byt []byte
	if byt, err = invoke.Command("docker", args...); err != nil {
		log.Errorf("cmd: docker %s, err: %s", strings.Join(args, " "), err.Error())
		log.Errorf("%s", string(byt))
		// return string(byt), err
	}
	return string(byt), err
}

/*******************************************************************
*
* Node Here!
*
 *******************************************************************/

//NodeList :传输Node的列表
type NodeList []*Node

//Node :docker node
type Node struct {
	NodeName    string `json:"NodeName"`
	IP          string `json:"IP"`
	NetworkName string `json:"NetworkName"`
}

//CreateNode :
func (node *Node) CreateNode() error {
	if err := node.CreateNodeDir(); err != nil {
		log.Errorf("id: %s , node: %+v, err: %s", node.NodeName, node, err.Error())
		return err
	}

	if err := node.CreateNodeContainer(); err != nil {
		log.Errorf("id: %s , node: %+v, err: %s", node.NodeName, node, err.Error())
		node.DestoryNodeDir()
		return err
	}

	if err := node.CreateNodeFile(); err != nil {
		log.Errorf("id: %s , node: %+v, err: %s", node.NodeName, node, err.Error())
		node.DestoryContainer()
		node.DestoryNodeDir()
		return err
	}
	return nil
}

//NodePath :Node文件存储路径
func (node *Node) NodePath() string {
	return fmt.Sprintf("%s/%s.node.d", apiConf.GetString("cdn_create_dir"), node.NodeName)
}

//NodeResourcePath :Resource路径，暂时保留
func (node *Node) NodeResourcePath() string {
	return fmt.Sprintf("%s/%s.node.d/resource.d", apiConf.GetString("cdn_create_dir"), node.NodeName)
}

//NginxConfPath :
func (node *Node) NginxConfPath() string {
	return fmt.Sprintf("%s/%s.node.d/nginx.conf.d", apiConf.GetString("cdn_create_dir"), node.NodeName)
}

//NginxHTTPConfPath :
func (node *Node) NginxHTTPConfPath() string {
	return fmt.Sprintf("%s/%s.node.d/nginx.conf.d/http.d", apiConf.GetString("cdn_create_dir"), node.NodeName)
}

//NginxTCPConfPath :
func (node *Node) NginxTCPConfPath() string {
	return fmt.Sprintf("%s/%s.node.d/nginx.conf.d/tcp.d", apiConf.GetString("cdn_create_dir"), node.NodeName)
}

//NginxUDPConfPath :
func (node *Node) NginxUDPConfPath() string {
	return fmt.Sprintf("%s/%s.node.d/nginx.conf.d/udp.d", apiConf.GetString("cdn_create_dir"), node.NodeName)
}

//CreateNodeDir :为每个Node建立目录，保存一些信息
func (node *Node) CreateNodeDir() error {
	if err := os.MkdirAll(node.NodePath(), os.ModePerm); err != nil {
		log.Errorf(err.Error())
		return err
	}

	if err := os.MkdirAll(node.NodeResourcePath(), os.ModePerm); err != nil {
		log.Errorf(err.Error())
		return err
	}

	if err := os.MkdirAll(node.NginxConfPath(), os.ModePerm); err != nil {

		log.Errorf(err.Error())
		return err
	}

	if err := os.MkdirAll(node.NginxHTTPConfPath(), os.ModePerm); err != nil {
		log.Errorf(err.Error())
		return err
	}

	if err := os.MkdirAll(node.NginxTCPConfPath(), os.ModePerm); err != nil {
		log.Errorf(err.Error())
		return err
	}
	if err := os.MkdirAll(node.NginxUDPConfPath(), os.ModePerm); err != nil {
		log.Errorf(err.Error())
		return err
	}
	return nil
}

//CreateNodeContainer :
func (node *Node) CreateNodeContainer() error {
	invoke := &Invoke{}
	args := []string{
		"run",
		"-d",
		"--restart=always",
		fmt.Sprintf("--net=%s", node.NetworkName),
		fmt.Sprintf("--ip=%s", node.IP),
		"-v",
		fmt.Sprintf("%s:/etc/nginx/conf.d", node.NginxConfPath()),
		fmt.Sprintf("--name=%s", node.NodeName),
		apiConf.GetString("cdn_image"),
	}

	if byt, err := invoke.Command("docker", args...); err != nil {
		log.Errorf("cmd: docker %s, err: %s", strings.Join(args, " "), err.Error())
		log.Errorf("%s", string(byt))
		return err
	}
	return nil
}

//DestoryNodeDir :
func (node *Node) DestoryNodeDir() error {
	return os.RemoveAll(node.NodePath())
}

//CreateNodeFile :
func (node *Node) CreateNodeFile() error {
	var error error
	var byt []byte
	if byt, error = yaml.Marshal(node); error == nil {
		f, err := os.OpenFile(fmt.Sprintf("%s/node.yaml", node.NodePath()), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(string(byt))
		return err
	}
	return error
}

//DestoryContainer :
func (node *Node) DestoryContainer() error {
	invoke := &Invoke{}

	if res, err := invoke.Command("docker", "stop", node.NodeName); err != nil {
		log.Errorf("[ docker stop %s ]error: %s\n", node.NodeName, string(res))
		return err
	}
	if res, err := invoke.Command("docker", "rm", node.NodeName); err != nil {
		log.Errorf("docker rm error: %s\n", string(res))
		return err
	}
	return nil
}

//DestoryNode :
func (node *Node) DestoryNode() error {
	if err := node.DestoryContainer(); err != nil {
		return err
	}
	if err := node.DestoryNodeDir(); err != nil {
		return err
	}
	return nil
}

//Load :
func (node *Node) Load(r io.Reader) error {
	byt, _ := ioutil.ReadAll(r)
	fmt.Print(string(byt))
	return yaml.Unmarshal(byt, &node)
}
