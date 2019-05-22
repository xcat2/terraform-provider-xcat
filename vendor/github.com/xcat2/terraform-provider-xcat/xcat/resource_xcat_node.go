package xcat

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	//"os"
	"bytes"
	"os/exec"
	"regexp"
	//"encoding/json"
	"reflect"
	//"github.com/jeremywohl/flatten"
	"github.com/tidwall/gjson"

	//"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var systemSyncLock sync.Mutex

func resourceNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceNodeCreate,
		Read:   resourceNodeRead,
		Update: resourceNodeUpdate,
		Delete: resourceNodeDelete,

		Schema: map[string]*schema.Schema{
			"selectors": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: false,
				ForceNew: true,
				/*
				   ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				                     v := val.([]string)
				                     re:=regexp.MustCompile(`([^=!~><]+)([=!~><]{1,2})([^=!~><]+)`)
				                     ValidSelectors:=make([]string,0,len(SelectorOpMaps))
				                     for k,_:=range SelectorOpMaps{
				                         ValidSelectors=append(ValidSelectors,k)
				                     }
				                     for _,line:=range v {
				                          match:=re.FindAllStringSubmatch(line, -1)
				                          if match != nil {
				                              attr,op,_:=match[0][1],match[0][2],match[0][3]
				                              if availops,ok:=SelectorOpMaps[attr]; ok{
				                                  if !Contains(availops, op) {
				                                      errs = append(errs,fmt.Errorf("invalid operation in selector \"%s\": the valid operation for selector \"%s\": %s",line,attr,availops))
				                                  }
				                              } else {
				                                      errs = append(errs,fmt.Errorf("invalid selector \"%s\": the valid selectors \"%s\"",line,strings.Join(ValidSelectors,",")))
				                              }
				                          }
				                     }
				                     return
				                   },
				*/
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"machinetype": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"arch": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"disksize": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cputype": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cpucount": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gpu": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ib": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mac": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rack": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"room": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"height": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"osimage": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"powerstatus": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, true),
			},
			"sshusername": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sshpassword": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		SchemaVersion: 0,
		//MigrateState: resourceExampleInstanceMigrateState,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

/*
func resourceExampleInstanceMigrateState(v int, inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
    switch v {
    case 0:
        log.Println("[INFO] Found Example Instance State v0; migrating to v1")
        return migrateExampleInstanceStateV0toV1(inst)
    default:
        return inst, fmt.Errorf("Unexpected schema version: %d", v)
    }
}

func migrateExampleInstanceStateV0toV1(inst *terraform.InstanceState) (*terraform.InstanceState, error) {
    if inst.Empty() {
        log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
        return inst, nil
    }

    if !strings.HasSuffix(inst.Attributes["name"], ".") {
        log.Printf("[DEBUG] Attributes before migration: %#v", inst.Attributes)
        inst.Attributes["name"] = inst.Attributes["name"] + "."
        log.Printf("[DEBUG] Attributes after migration: %#v", inst.Attributes)
    }

    return inst, nil
}
*/

func resourceNodeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	selectors := Intf2Map(d.Get("selectors"))
	log.Printf("----------------%v", selectors)

	nodename := d.Get("name")
	if nodename != nil && nodename != "" {
		selectors["name"] = nodename.(string)
	}

	username := config.Username
	systemSyncLock.Lock()
	errcode, out := occupynode(selectors, username)
	systemSyncLock.Unlock()
	if errcode != 0 {
		return fmt.Errorf(out)
	}

	node := out

	osimage := d.Get("osimage")
	provisioned := 0
	if osimage != nil && osimage != "" {
		netbootparam := NetbootParam{
			osimage: osimage.(string),
		}

		errcode, errmsg := Rinstall(node, &netbootparam)
		if errcode != 0 {
			log.Printf("releasenode %s from %s", node, username)
			releasenode(node, username)
			out := "Failed to provision node " + node + ":" + errmsg
			return fmt.Errorf(out)
		}

		log.Printf("%v", resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			errcode, out := GetStatus(node)

			if errcode != 0 {
				log.Printf("Error to get status of node %s: %s", node, out)
				return resource.NonRetryableError(fmt.Errorf("Error to get status of node %s: %s", node, out))
			}

			if out != "booted" {
				log.Printf("Expected instance to be \"booted\" but was in state %s", out)
				return resource.RetryableError(fmt.Errorf("Expected instance to be \"booted\" but was in state %s", out))
			}

			provisioned = 1
			log.Printf("instance %s provisioned!", node)
			return resource.NonRetryableError(fmt.Errorf("instance %s provisioned!", node))
		}))

		if provisioned == 0 {
			releasenode(node, username)
			return fmt.Errorf("node instance %s provision timeout!", node)
		}
	}

	powerstatus := d.Get("powerstatus")
	if powerstatus == "off" {
		errcode, out := Power(node, "off")
		if errcode != 0 {
			releasenode(node, username)
			return fmt.Errorf("Fail to set powerstatus of instance  %s to %s: %s!", node, "off", out)
		}
	} else if powerstatus == "on" {
		errcode, out := Power(node, "on")
		if errcode != 0 {
			releasenode(node, username)
			return fmt.Errorf("Fail to set powerstatus of instance  %s to %s: %s!", node, "on", out)
		}
	}

	d.SetId(node)
	d.Set("name", node)
	return resourceNodeRead(d, meta)
}

func resourceNodeRead(d *schema.ResourceData, meta interface{}) error {
	node := d.Get("name").(string)
	cmd := exec.Command("xcat-inventory", "export", "-t", "node", "-o", node, "--format", "json")
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to read node resource " + node + " from xcat: " + errbuf.String())
	}

	NodeInv2Res(outbuf.String(), d, node)
	return nil
}

func resourceNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	//config := meta.(*Config)
	node := d.Get("name").(string)

	d.Partial(true)
	if d.HasChange("osimage") {
		oldOsimage_v, newOsimage_v := d.GetChange("osimage")
		oldOsimage := oldOsimage_v.(string)
		newOsimage := newOsimage_v.(string)
		log.Printf("%s=========%s", oldOsimage, newOsimage)
		osimage := newOsimage
		if osimage != "" {
			netbootparam := NetbootParam{
				osimage: osimage,
			}

			errcode, errmsg := Rinstall(node, &netbootparam)
			if errcode != 0 {
				out := "Failed to provision node " + node + ":" + errmsg
				return fmt.Errorf(out)
			}

			provisioned := 0

			log.Printf("%v", resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				errcode, out := GetStatus(node)

				if errcode != 0 {
					log.Printf("Error to get status of node %s: %s", node, out)
					return resource.NonRetryableError(fmt.Errorf("Error to get status of node %s: %s", node, out))
				}

				if out != "booted" {
					log.Printf("Expected instance to be \"booted\" but was in state %s", out)
					return resource.RetryableError(fmt.Errorf("Expected instance to be \"booted\" but was in state %s", out))
				}

				provisioned = 1
				log.Printf("instance %s provisioned!", node)
				return resource.NonRetryableError(fmt.Errorf("instance %s provisioned!", node))
			}))

			if provisioned == 0 {
				return fmt.Errorf("node instance %s provision timeout!", node)
			}
			d.SetPartial("osimage")
		}
	}

	if d.HasChange("powerstatus") {
		//oldPowerStatus_v, newPowerStatus_v := d.GetChange("powerstatus")
		_, newPowerStatus_v := d.GetChange("powerstatus")
		//oldPowerStatus:=oldPowerStatus_v.(string)
		newPowerStatus := newPowerStatus_v.(string)

		errcode, out := Power(node, newPowerStatus)
		if errcode != 0 {
			return fmt.Errorf("Fail to set powerstatus of instance  %s to %s: %s!", node, newPowerStatus, out)
		}
		d.SetPartial("powerstatus")
	}

	d.Partial(false)
	return resourceNodeRead(d, meta)
}

func resourceNodeDelete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	username := config.Username
	node := d.Get("name").(string)
	errorcode, errormessage := releasenode(node, username)
	if errorcode != 0 {
		return fmt.Errorf(errormessage)
	}
	return nil
}

func selectnodes(selector map[string]string) []string {
	var cmdslice []string
	log.Printf("selector=%v", selector)
	if node, ok := selector["name"]; ok {
		nodegroups, errcode, _ := getnodegroups(node)
		if errcode != 0 {
			return nil
		}

		if nodegroups != "free" {
			return nil
		}
		cmdslice = []string{"lsdef", "-t", "node", node, "-s"}
		delete(selector, "name")
	} else {
		cmdslice = []string{"lsdef", "-t", "node", "free", "-s"}
	}
	log.Printf("selector=%v", selector)
	for key, value := range selector {
		if _, ok := DictRes2Inv[key]; ok {
			cmdslice = append(cmdslice, "-w", Res2DefAttr(key)+"=="+value)
		} else {
			cmdslice = append(cmdslice, "-w", "usercomment=~,"+Res2DefAttr(key)+"="+value+",")
		}
	}
	log.Printf("cmdslice=%v", cmdslice)
	cmd := exec.Command(cmdslice[0], cmdslice[1:]...)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		return nil
	}

	cmdout := outbuf.String()
	log.Printf("cmdout=%v", cmdout)

	var nodelist []string
	nodelist = nil

	if cmdout != "" {
		var rgx = regexp.MustCompile(`\s*(\b[^(]\S+[^)]\b)\s+\(node\)`)
		rs := rgx.FindAllStringSubmatch(cmdout, -1)
		log.Printf("rs=%v\n", rs)
		for _, mylist := range rs {
			log.Printf("VVVVVVVV %v", mylist)
			nodelist = append(nodelist, mylist[1])
		}
	}
	return nodelist
}

func occupynode(selectors map[string]string, user string) (int, string) {
	nodelist := selectnodes(selectors)
	if nodelist == nil {
		return 1, "cannot find requested  node resources"
	}

	cmd := exec.Command("chdef", "-t", "node", "-o", nodelist[0], "groups="+user)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		return 1, "Failed to occupy node resource " + nodelist[0] + " for user " + user + ": " + errbuf.String()
	}
	return 0, nodelist[0]
}

func getnodegroups(node string) (string, int, string) {
	var nodegroups = ""
	cmd := exec.Command("lsdef", "-t", "node", "-o", node, "-i", "groups")
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		return "", 1, "Failed to apply node resource " + node
	}

	cmdout := outbuf.String()
	if cmdout != "" {
		var rgx = regexp.MustCompile(`groups=(.*)`)
		rs := rgx.FindStringSubmatch(cmdout)
		nodegroups = rs[1]

	}

	return nodegroups, 0, ""
}

func releasenode(node string, user string) (int, string) {
	nodegroups, errcode, errmsg := getnodegroups(node)
	if errcode != 0 {
		return errcode, errmsg
	}

	nodegrouplist := strings.Split(nodegroups, ",")
	if !Contains(nodegrouplist, user) {
		return 0, ""
	}

	cmd := exec.Command("chdef", "-t", "node", "-o", node, "groups=free")
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		return 1, "Failed to release node resource " + node + " from user " + user + ": " + errbuf.String()
	}

	return 0, ""
}

func getattr(key string, d *schema.ResourceData) (int, string) {
	keylist := strings.Split(key, ".")
	key, keylist = keylist[0], keylist[1:]

	vRaw, ok := d.GetOk(key)
	if !ok {
		return -1, "unexist key"
	}

	for {
		if len(keylist) == 0 {
			return 0, vRaw.(string)
		}

		if typeof(vRaw) == "*schema.Set" {
			componentRaw := vRaw.(*schema.Set).List()
			raw := componentRaw[0]
			rawMap := raw.(map[string]interface{})
			key, keylist = keylist[0], keylist[1:]
			vRaw = rawMap[key]
		}
	}

}

func Intf2Map(v interface{}) map[string]string {
	m := v.([]interface{})
	retmap := make(map[string]string)
	re := regexp.MustCompile(`([^=!~><]+)([=!~><]{1,2})([^=!~><]+)`)
	for _, line := range m {
		match := re.FindAllStringSubmatch(line.(string), -1)
		if match != nil {
			attr, _, value := match[0][1], match[0][2], match[0][3]
			retmap[attr] = value
		}

	}
	return retmap
}

/*
func Intf2Map(v interface{}) map[string]string {
	m := v.(map[string]interface{})
        retmap:=make(map[string]string)
        for key,value :=range m{
            retmap[key]=value.(string)
        }
	return retmap
}
*/

var SelectorOpMaps = map[string][]string{
	"disksize":    []string{"=", ">", ">=", "<", "<="},
	"memory":      []string{"=", ">", ">=", "<", "<="},
	"cpucount":    []string{"=", ">", ">=", "<", "<="},
	"cputype":     []string{"=", "!=", "!~", "=~"},
	"machinetype": []string{"="},
	"name":        []string{"="},
	"rack":        []string{"="},
	"unit":        []string{"="},
	"room":        []string{"="},
	"arch":        []string{"="},
	"gpu":         []string{"="},
	"ib":          []string{"="},
}

var DictRes2Inv = map[string]string{
	"machinetype": "device_info.mtm",
	"arch":        "device_info.arch",
	"disksize":    "device_info.disksize",
	"memory":      "device_info.memory",
	"cputype":     "device_info.cputype",
	"cpucount":    "device_info.cpucount",
	"ip":          "network_info.primarynic.ip",
	"mac":         "network_info.primarynic.mac",
	"rack":        "position_info.rack",
	"unit":        "position_info.unit",
	"room":        "position_info.room",
	"height":      "position_info.height",
	"osimage":     "engines.netboot_engine.engine_info.osimage",
}

func Res2DefAttr(resattr string) string {
	if resattr == "machinetype" {
		return "mtm"
	}
	return resattr
}

func NodeInv2Res(myjson string, d *schema.ResourceData, node string) int {
	keys := reflect.ValueOf(DictRes2Inv).MapKeys()
	for _, kres := range keys {
		kinv := DictRes2Inv[kres.String()]
		val := gjson.Get(myjson, "node."+node+"."+kinv).String()
		if val != "" {
			d.Set(kres.String(), val)
		} else {
			d.Set(kres.String(), nil)
		}
	}
	return 0
}

func RunCmd(cmdstr string, args ...string) (error, string, string) {
	cmd := exec.Command(cmdstr, args...)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	return err, outbuf.String(), errbuf.String()
}

type NetbootParam struct {
	osimage     string
	addkcmdline string
}

func Power(node string, action string) (int, string) {
	err, _, errstr := RunCmd("rpower", node, action)
	if err != nil {
		return 1, errstr
	}
	return 0, ""
}

func Rinstall(node string, param *NetbootParam) (int, string) {
	err, _, errstr := RunCmd("makedns", node)
	err, _, errstr = RunCmd("rinstall", node, "osimage="+param.osimage)
	if err != nil {
		return 1, errstr
	}
	return 0, ""
}

func GetStatus(node string) (int, string) {
	err, outstr, errstr := RunCmd("lsdef", "-t", "node", "-o", node, "-i", "status")
	if err != nil {
		return 1, errstr
	}

	var myregex = regexp.MustCompile(`status=(\w+)`)
	match := myregex.FindAllStringSubmatch(outstr, 1)
	if match == nil {
		return 1, "invalid output: " + outstr
	}

	return 0, match[0][1]
}
