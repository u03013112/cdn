package cdn

import (
	"errors"
	"fmt"
	"strconv"
	// "log"
	"os"
	"strings"

	"github.com/fleacloud/app"
	"github.com/spf13/cobra"
)

// var cdnHost = "127.0.0.1:8000"

func cdnJsonPost(method string, url string, resp interface{}, req interface{}) error {
	cdnHost := os.Getenv("CDN_HOST")
	if cdnHost == "" {
		cdnHost = "127.0.0.1:8000"
	}
	api := fmt.Sprintf("http://%s%s", cdnHost, url)
	return JsonPost(strings.ToUpper(method), api, resp, req)
}

func cdnCtlGetRulesCmdRunE(cmd *cobra.Command, args []string) error {
	// var list []*StreamRule

	// if err := cdnJsonPost("GET", "/nodes/rules", &list, nil); err != nil {
	// 	log.Fatalf("err:%s\n", err.Error())
	// 	return err
	// }

	// tbl, _ := prettytable.NewTable([]prettytable.Column{
	// 	{Header: "NodeId", AlignRight: false},
	// 	{Header: "StreamId", AlignRight: false},
	// 	{Header: "Proto", AlignRight: false},
	// 	{Header: "Resource", AlignRight: false},
	// 	{Header: "SiteIp", AlignRight: false},
	// 	{Header: "SitePort", AlignRight: false},
	// }...)
	// tbl.Separator = " "

	// for _, rulelist := range list {
	// 	for _, rule := range rulelist.Site {
	// 		tbl.AddRow([]interface{}{
	// 			rulelist.NodeId,
	// 			rulelist.StreamId,
	// 			rulelist.Proto,
	// 			rulelist.NodeResource,
	// 			rule.SiteIp,
	// 			rule.SitePort,
	// 		}...)
	// 	}
	// }

	// if _, err := tbl.Print(); err != nil {
	// 	log.Errorf(err.Error())
	// }

	return nil
}

var cdnCtlGetRulesCmd = &cobra.Command{
	Use:     "rules",
	Aliases: []string{"r", "ru", "rul", "rule"},
	Short:   "get all rules from cdn devices.",
	Long:    "get all rules from cdn devices.",
	RunE:    cdnCtlGetRulesCmdRunE,
}

//"cdnctl add {{nodeid}} {{proto}} {{resource}} {{siteip}}:{{siteport}}"
func cdnCtlAddRuleCmdRunE(cmd *cobra.Command, args []string) error {
	// if len(args) != 4 {
	// 	return errors.New("Not enough args! mast be 4!")
	// }

	// var rule StreamRule

	// rule.NodeId = args[0]
	// // rule.StreamId = args[0]
	// rule.Proto = args[1]
	// rule.NodeResource = args[2]

	// siteArr := strings.Split(args[3], ",")
	// for _, one := range siteArr {
	// 	arr := strings.Split(one, ":")

	// 	var site SiteInfo
	// 	site.SiteIp = arr[0]
	// 	site.SitePort = arr[1]
	// 	rule.Site = append(rule.Site, &site)
	// }

	// // }

	// if err := cdnJsonPost("POST", "/nodes/rules", nil, &rule); err != nil {
	// 	log.Fatalf("%s\n", err.Error())
	// 	return err
	// }
	// fmt.Println("success")
	return nil
}

//"cdnctl del {{nodeid}} {{proto}} {{resource}}"
func cdnCtlDelRuleCmdRunE(cmd *cobra.Command, args []string) error {
	// if len(args) != 3 && len(args) != 4 {
	// 	return errors.New("Not enough args! mast be 3 or 4!")
	// }

	// var rule StreamRule

	// rule.NodeId = args[0]
	// // rule.StreamId = args[0]
	// rule.Proto = args[1]
	// rule.NodeResource = args[2]

	// if err := cdnJsonPost("DELETE", "/nodes/rules", nil, &rule); err != nil {
	// 	log.Fatalf("%s\n", err.Error())
	// 	return err
	// }
	// fmt.Println("success")
	return nil
}

var cdnCtlAddRuleCmd = &cobra.Command{
	Use:     "rules",
	Aliases: []string{"r", "ru", "rul", "rule"},
	Short:   "cdnctl add rule {{nodeid}} {{proto}} {{resource}} {{siteip}}:{{siteport}}",
	Long:    "cdnctl add rule {{nodeid}} {{proto}} {{resource}} {{siteip}}:{{siteport}}",
	RunE:    cdnCtlAddRuleCmdRunE,
}

var cdnCtlDelRuleCmd = &cobra.Command{
	Use:     "rules",
	Aliases: []string{"r", "ru", "rul", "rule"},
	Short:   "cdnctl del rule {{nodeid}} {{proto}} {{resource}}",
	Long:    "cdnctl del rule {{nodeid}} {{proto}} {{resource}}",
	RunE:    cdnCtlDelRuleCmdRunE,
}

var cdnCtlGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g", "ge"},
	Short:   "get object from cdn devices.",
	Long:    "get object from cdn devices.",
}
var cdnCtlAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a", "ad", "addd", "inster"},
	Short:   "add object to cdn devices.",
	Long:    "add object to cdn devices.",
}

var cdnCtlDelCmd = &cobra.Command{
	Use:     "del",
	Aliases: []string{"d", "de", "delete"},
	Short:   "del object from cdn devices.",
	Long:    "del object from cdn devices.",
}

var cdnCtlRootCmd = &cobra.Command{
	Use:     "cdnctl",
	Aliases: []string{"ctl", "c", "cc"},
	Short:   "env to use \"export CDN_HOST=127.0.0.1:8000\"",
	Long:    "env to use \"export CDN_HOST=127.0.0.1:8000\"",
}

func init() {
	//get
	cdnCtlGetCmd.AddCommand(cdnCtlGetNodesCmd)
	cdnCtlGetCmd.AddCommand(cdnCtlGetRulesCmd)
	cdnCtlGetCmd.AddCommand(cdnCtlGetNetworkCmd)

	//add
	cdnCtlAddCmd.AddCommand(cdnCtlAddRuleCmd)
	cdnCtlAddCmd.AddCommand(cdnCtlAddNodeCmd)
	cdnCtlAddCmd.AddCommand(cdnCtlAddNetworkCmd)

	//del
	cdnCtlDelCmd.AddCommand(cdnCtlDelRuleCmd)
	cdnCtlDelCmd.AddCommand(cdnCtlDelNodeCmd)
	cdnCtlDelCmd.AddCommand(cdnCtlDelNetworkCmd)

	//root
	cdnCtlRootCmd.AddCommand(cdnCtlGetCmd)
	cdnCtlRootCmd.AddCommand(cdnCtlAddCmd)
	cdnCtlRootCmd.AddCommand(cdnCtlDelCmd)

	app.AddCommond(cdnCtlRootCmd)
}

/*
*
*	Network :docker network ctrl
*
 */
//Add
var cdnCtlAddNetworkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"n", "ne", "net", "netw"},
	Short:   "cdnctl add network {{networkName}} {{interface}} {{ipaddress}}/{{maskLen}}@{{gateway}} {{VLAN}}",
	Long:    "cdnctl add network {{networkName}} {{interface}} {{ipaddress}}/{{maskLen}}@{{gateway}} {{VLAN}}",
	RunE:    cdnCtlAddNetworkCmdRunE,
}

func cdnCtlAddNetworkCmdRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 3 && len(args) != 4 {
		fmt.Printf("%d\n", len(args))
		return errors.New("not enough args, mast be 3 or 4")
	}

	var network Network

	network.NetworkName = args[0]
	network.Interface = args[1]

	net := args[2]
	arr := strings.Split(net, "@")
	network.Gateway = arr[1]
	arr = strings.Split(arr[0], "/")
	network.NetworkIP = arr[0]
	network.Mask = arr[1]

	if len(args) == 4 {
		vlan, e := strconv.Atoi(args[3])
		if e == nil {
			network.Vlan = vlan
		} else {
			network.Vlan = 0
		}
	} else {
		network.Vlan = 0
	}

	if err := cdnJsonPost("POST", "/networks", nil, &network); err != nil {
		log.Fatalf("%s\n", err.Error())
		return err
	}
	fmt.Println("success")
	return nil
}

//DELETE
var cdnCtlDelNetworkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"n", "ne", "net", "netw"},
	Short:   "cdnctl del network {{networkName}}",
	Long:    "cdnctl del network {{networkName}}",
	RunE:    cdnCtlDelNetworkCmdRunE,
}

func cdnCtlDelNetworkCmdRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("need networkName")
	}

	var network Network

	network.NetworkName = args[0]

	if err := cdnJsonPost("DELETE", "/networks", nil, &network); err != nil {
		log.Fatalf("%s\n", err.Error())
		return err
	}
	fmt.Println("success")
	return nil
}

//GET
var cdnCtlGetNetworkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"n", "ne", "net", "netw"},
	Short:   "cdnctl get network",
	Long:    "cdnctl get network",
	RunE:    cdnCtlGetNetworksCmdRunE,
}

func cdnCtlGetNetworksCmdRunE(cmd *cobra.Command, args []string) error {
	var list NetworkList

	if err := cdnJsonPost("GET", "/networks", &list, nil); err != nil {
		log.Fatalf("err:%s\n", err.Error())
		return err
	}

	for _, n := range list {
		fmt.Printf("%s NET:%s GW:%s Int:%s\n", n.NetworkName, n.NetworkIP, n.Gateway, n.Interface)
	}

	return nil
}

/*
*
*	Node :docker node ctrl
*
 */

//Add
var cdnCtlAddNodeCmd = &cobra.Command{
	Use:     "node",
	Aliases: []string{"no", "n", "node", "nod"},
	Short:   "cdnctl add node {{nodeName}} {{ip}} {{networkName}}",
	Long:    "cdnctl add node {{nodeName}} {{ip}} {{networkName}}",
	RunE:    cdnCtlAddNodeCmdRunE,
}

func cdnCtlAddNodeCmdRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		fmt.Printf("%d\n", len(args))
		return errors.New("args err")
	}

	var node Node

	node.NodeName = args[0]

	node.IP = args[1]
	node.NetworkName = args[2]

	if err := cdnJsonPost("POST", "/nodes", nil, &node); err != nil {
		log.Fatalf("%s\n", err.Error())
		return err
	}
	fmt.Println("success")
	return nil
}

//DEL
var cdnCtlDelNodeCmd = &cobra.Command{
	Use:     "node",
	Aliases: []string{"no", "n", "node", "nod"},
	Short:   "cdnctl del node {{nodeName}}",
	Long:    "cdnctl del node {{nodeName}}",
	RunE:    cdnCtlDelNodeCmdRunE,
}

func cdnCtlDelNodeCmdRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("args err")
	}

	var node Node

	node.NodeName = args[0]

	if err := cdnJsonPost("DELETE", "/nodes", nil, &node); err != nil {
		log.Fatalf("%s\n", err.Error())
		return err
	}
	fmt.Println("success")
	return nil
}

//GET
var cdnCtlGetNodesCmd = &cobra.Command{
	Use:     "node",
	Aliases: []string{"no", "n", "node", "nod"},
	Short:   "cdnctl get node",
	Long:    "cdnctl get node",
	RunE:    cdnCtlGetNodesCmdRunE,
}

func cdnCtlGetNodesCmdRunE(cmd *cobra.Command, args []string) error {
	var list NodeList

	if err := cdnJsonPost("GET", "/nodes", &list, nil); err != nil {
		log.Fatalf("err:%s\n", err.Error())
		return err
	}

	for _, n := range list {
		fmt.Printf("%s IP:%s Network:%s \n", n.NodeName, n.IP, n.NetworkName)
	}

	return nil
}
