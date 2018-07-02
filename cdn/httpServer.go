package cdn

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	// "log"
	"net/http"
	"strings"

	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

func httpRouter(m *macaron.Macaron) {

	m.Get("/nodes", ShowNodes)
	m.Post("/nodes", binding.Bind(Node{}), AddNode)
	m.Delete("/nodes", binding.Bind(Node{}), DelNode)

	m.Get("/rules4", ShowRules4)
	m.Post("/rules4", binding.Bind(Rule4{}), AddRule4)
	m.Delete("/rules4", binding.Bind(Rule4{}), DelRule4)

	m.Get("/networks", ShowNetworks)
	m.Post("/networks", binding.Bind(Network{}), AddNetwork)
	m.Delete("/networks", binding.Bind(Network{}), DelNetwork)
}

func httpServer() error {
	m := macaron.Classic()
	m.Use(macaron.Renderer(macaron.RenderOptions{
		IndentJSON: true,
	}))
	httpRouter(m)
	addr := fmt.Sprintf("%s:%s", apiConf.GetString("listen_addr"), apiConf.GetString("listen_port"))
	log.Info("Server is running on [%s]...\n", addr)
	log.Info(http.ListenAndServe(addr, m))
	return nil
}

/*******************************************************************
*
* Network Here!
*
 *******************************************************************/

//AddNetwork : for Network
func AddNetwork(n Network) string {
	if err := n.CreateNetwork(); err != nil {
		log.Errorf(err.Error())
		return err.Error()
	}
	return "ok"
}

//DelNetwork : 暂时还没有写有引用的不能删除，目前全靠docker自己报错
func DelNetwork(n Network) string {
	if err := n.DestoryNetwork(); err != nil {
		log.Errorf(err.Error())
		return err.Error()
	}
	return "ok"
}

//ShowNetworks :
func ShowNetworks(ctx *macaron.Context) {
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
	}

	var list NetworkList

	networkArr := strings.Split(string(byt), "\n")
	for _, one := range networkArr {
		rowArr := strings.Fields(one)
		if len(rowArr) > 2 {
			if rowArr[0] == "NETWORK" {
				continue
			}

			if rowArr[1] == "bridge" || rowArr[1] == "host" || rowArr[1] == "none" {
				continue
			}

			var network Network
			network.NetworkName = rowArr[1]
			list = append(list, &network)
		} else {
			continue
		}
	}

	for _, network := range list {
		invoke := &Invoke{}

		args := []string{
			"network",
			"inspect",
			network.NetworkName,
		}

		if byt, err = invoke.Command("docker", args...); err != nil {
			log.Errorf("cmd: docker %s, err: %s", strings.Join(args, " "), err.Error())
			log.Errorf("%s", string(byt))
		} else {
			type Config struct {
				Subnet  string `json:"Subnet"`
				Gateway string `json:"Gateway"`
			}
			type IPAM struct {
				ConfigArr []Config `json:"Config"`
			}
			type Options struct {
				Parent string `json:"Parent"`
			}
			type NetworkInspect struct {
				Name    string  `json:"Name"`
				Ipam    IPAM    `json:"IPAM"`
				Options Options `json:"Options"`
			}
			var ret []NetworkInspect

			json.Unmarshal(byt, &ret)

			network.NetworkIP = ret[0].Ipam.ConfigArr[0].Subnet
			network.Gateway = ret[0].Ipam.ConfigArr[0].Gateway
			network.Interface = ret[0].Options.Parent
		}
	}

	ctx.JSON(200, &list)
}

/*******************************************************************
*
* Node Here!
*
 *******************************************************************/

//AddNode :
func AddNode(n Node) string {
	invoke := &Invoke{}
	args := []string{"ps"}

	var err error
	var byt []byte
	if byt, err = invoke.Command("docker", args...); err != nil {
		log.Errorf("cmd: docker %s, err: %s", strings.Join(args, " "), err.Error())
		log.Errorf("%s", string(byt))
	}

	nodeArr := strings.Split(string(byt), "\n")
	for _, one := range nodeArr {
		rowArr := strings.Fields(one)
		if len(rowArr) > 2 {
			if rowArr[0] == "CONTAINER" {
				continue
			}

			if rowArr[len(rowArr)-1] == "cdnapi" {
				continue
			}

			if n.NodeName == rowArr[len(rowArr)-1] {
				log.Errorf("name alreay in used.")
				return "name alreay in used."
			}
		} else {
			continue
		}
	}

	if err := n.CreateNode(); err != nil {
		log.Errorf(err.Error())
		return err.Error()
	}
	return "ok"
}

//DelNode :
func DelNode(n Node) string {
	if err := n.DestoryNode(); err != nil {
		return err.Error()
	}
	return "ok"
}

//ShowNodes :
func ShowNodes(ctx *macaron.Context) {
	var list NodeList

	invoke := &Invoke{}
	args := []string{"ps"}

	var err error
	var byt []byte
	if byt, err = invoke.Command("docker", args...); err != nil {
		log.Errorf("cmd: docker %s, err: %s", strings.Join(args, " "), err.Error())
		log.Errorf("%s", string(byt))
	}

	nodeArr := strings.Split(string(byt), "\n")
	for _, one := range nodeArr {
		rowArr := strings.Fields(one)
		if len(rowArr) > 2 {
			if rowArr[0] == "CONTAINER" {
				continue
			}

			if rowArr[len(rowArr)-1] == "cdnapi" {
				continue
			}

			var node Node
			node.NodeName = rowArr[len(rowArr)-1]
			list = append(list, &node)
		} else {
			continue
		}
	}

	for _, node := range list {
		nodefilepath := fmt.Sprintf("%s/node.yaml", node.NodePath())
		nf, _ := os.Open(nodefilepath)
		node.Load(nf)
	}

	ctx.JSON(200, &list)
}

/*******************************************************************
*
* Rule Here!
*
 *******************************************************************/

//AddRule4 :
func AddRule4(r Rule4) string {

	n := &Node{
		NodeName: r.NodeName,
	}

	if err := n.AddRule4(&r); err != nil {
		return err.Error()
	}
	return "ok"
}

//DelRule4 :
func DelRule4(r Rule4) string {
	n := &Node{
		NodeName: r.NodeName,
	}

	if err := n.DelRule4(&r); err != nil {
		return err.Error()
	}
	return "ok"
}

//ShowRules4 :由于不会传参数，就一次性显示
func ShowRules4(ctx *macaron.Context) {
	dir := apiConf.GetString("cdn_create_dir")

	// var list NodeList
	var list []*Rule4
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			return nil
		}
		if !strings.HasSuffix(f.Name(), ".node.d") {
			return nil
		}

		onelist, err := searchRules(fmt.Sprintf("%s/resource.d", path))
		if err != nil {
			log.Errorf(err.Error())
			return err
		}

		list = append(list, onelist...)
		return nil
	})

	if err != nil {
		log.Errorf(err.Error())
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, &list)
}

func searchRules(dir string) ([]*Rule4, error) {
	var list []*Rule4
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			log.Errorf(err.Error())
			return err
		}
		if f.IsDir() {
			return nil
		}
		var rule Rule4
		nf, err := os.Open(path)
		if err != nil {
			log.Errorf(err.Error())
			return err
		}
		if err := rule.Load(nf); err != nil {
			log.Errorf(err.Error())
			return err
		}
		list = append(list, &rule)
		return nil
	})
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return list, nil
}
