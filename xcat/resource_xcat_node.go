package xcat

import (
	"fmt"
	"log"
	"strings"
	"sync"
        //"os"
        "os/exec"
        "bytes"
        "regexp"
        //"encoding/json"
        "reflect"
        //"github.com/jeremywohl/flatten"
        "github.com/thedevsaddam/gojsonq"
        
	//"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
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
                                  Type:     schema.TypeMap,
                                  Optional: true,
                                  Computed: false,
                                  ForceNew: true,
                        },
                        "name": {
                                  Type:     schema.TypeString,
                                  Optional: true,
                                  Computed: true,
                                  ForceNew: true,
                        },
                        "mtm": {
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
                        "serial": {
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
                        "firmware": {
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
                                  Type:     schema.TypeString,
                                  Optional: true,
                                  Computed: true,
                        },
                        "zone": {
                                  Type:     schema.TypeString,
                                  Optional: true,
                                  Computed: true,
                        },
		},
	}
}

func resourceNodeCreate(d *schema.ResourceData, meta interface{}) error {
	//systemSyncLock.Lock()
	//defer systemSyncLock.Unlock()

	config := meta.(*Config)
        selectors:=Intf2Map(d.Get("selectors"))
        log.Printf("----------------%v",selectors)

        nodename:=d.Get("name")
        if nodename != nil && nodename != ""{
            selectors["name"]=nodename.(string)
        }

        username:=config.Username
	systemSyncLock.Lock()
        errcode,out:= occupynode(selectors, username)
	systemSyncLock.Unlock()
        if errcode!=0 {
            return fmt.Errorf(out)
        }

        
        

        node:=out

        osimage:=d.Get("osimage")
        if osimage!=nil && osimage!= ""{
            netbootparam:=NetbootParam{
                osimage:osimage.(string),
            } 

            errcode,errmsg:=ProvisionNode(node,&netbootparam)
            if errcode!=0 {
                log.Printf("releasenode %s from %s",node,username)
                releasenode(node,username)
                out:="Failed to provision node "+node+":"+errmsg
                return fmt.Errorf(out)
            }
        }
 
        d.SetId(node)
        d.Set("name",node)
        log.Printf("[INFO] there is a pending resize operation on this pool...")
	return resourceNodeRead(d, meta)
}

func resourceNodeRead(d *schema.ResourceData, meta interface{}) error {
        node:=d.Get("name").(string)
        cmd := exec.Command("xcat-inventory","export","-t","node","-o",node,"--format","json")
        var outbuf, errbuf bytes.Buffer
        cmd.Stdout = &outbuf
        cmd.Stderr = &errbuf
            
        err := cmd.Run()        
        if err != nil {
               log.Printf("Failed to read node resource "+node+" from xcat: "+errbuf.String())
        }

        mynodejson :=gojsonq.New().JSONString(outbuf.String())
        NodeInv2Res(mynodejson, d,node)
	return nil
}

func resourceNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	//systemSyncLock.Lock()
	//defer systemSyncLock.Unlock()

	//config := meta.(*Config)
        node:=d.Get("name").(string)

        if d.HasChange("osimage") {
            oldOsimage_v, newOsimage_v := d.GetChange("osimage")
            oldOsimage:=oldOsimage_v.(string)
            newOsimage:=newOsimage_v.(string)
            log.Printf("%s=========%s",oldOsimage,newOsimage)
            osimage:=newOsimage
            if osimage!= ""{
                netbootparam:=NetbootParam{
                    osimage:osimage,
                } 

                errcode,errmsg:=ProvisionNode(node,&netbootparam)
                if errcode!=0 {
                    out:="Failed to provision node "+node+":"+errmsg
                    return fmt.Errorf(out)
                }
            }
        }

	return resourceNodeRead(d, meta)
}

func resourceNodeDelete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
        username:=config.Username
        node:=d.Get("name").(string)
        errorcode,errormessage := releasenode(node,username)
        if errorcode!=0 {
            return fmt.Errorf(errormessage )
        }
	return nil
}


func selectnodes(selector map[string]string) ([]string) {
        var cmdslice []string
        log.Printf("selector=%v",selector)
        if node,ok:=selector["name"]; ok{
            nodegroups,errcode,_:=getnodegroups(node)
            if errcode !=0 {
                return nil
            }
            
            if nodegroups != "free" {
                return nil
            } 
            cmdslice=[]string{"lsdef", "-t", "node",node, "-s"}
            delete(selector,"name")
        } else {
            cmdslice=[]string{"lsdef", "-t", "node","free", "-s"}
        }
        log.Printf("selector=%v",selector)
        for key,value :=range selector {
             cmdslice = append(cmdslice,"-w",Res2DefAttr(key)+"=="+value)
        }
        log.Printf("cmdslice=%v",cmdslice)
        cmd := exec.Command(cmdslice[0],cmdslice[1:]...)
        var outbuf, errbuf bytes.Buffer
        cmd.Stdout = &outbuf
        cmd.Stderr = &errbuf

        err := cmd.Run()        
        if err != nil {
               return nil
        }
       
        cmdout:=outbuf.String()
        log.Printf("%v",cmdout)
        
        var nodelist []string
        nodelist=nil

        if cmdout != "" {
           var rgx=regexp.MustCompile(`\s*(\b[^(]\S+[^)]\b)\s+\(node\)`) 
           rs := rgx.FindAllStringSubmatch(cmdout,-1)
           log.Printf("rs=%v\n",rs)
           for _,mylist := range rs {
               log.Printf("VVVVVVVV %v",mylist)
               nodelist=append(nodelist,mylist[1])
           }
        }
        return nodelist
}

func occupynode(selectors map[string]string, user string) (int,string) {
     nodelist:=selectnodes(selectors)
     if nodelist ==nil {
         return 1, "cannot find requested  node resources"
     }


     cmd := exec.Command("chdef","-t","node","-o",nodelist[0],"groups="+user)
     var outbuf, errbuf bytes.Buffer
     cmd.Stdout = &outbuf
     cmd.Stderr = &errbuf
         
     err := cmd.Run()        
     if err != nil {
            return 1,"Failed to occupy node resource "+nodelist[0]+" for user "+user +": "+errbuf.String()
     }
     return 0,nodelist[0]
}


func getnodegroups(node string) (string,int,string) {
        var nodegroups=""
        cmd := exec.Command("lsdef", "-t", "node","-o",node,"-i","groups")
        var outbuf, errbuf bytes.Buffer
        cmd.Stdout = &outbuf
        cmd.Stderr = &errbuf

        err := cmd.Run()        
        if err != nil {
               return "",1,"Failed to apply node resource "+node
        }
       
        cmdout:=outbuf.String()
        if cmdout != "" {
           var rgx=regexp.MustCompile(`groups=(.*)`) 
           rs := rgx.FindStringSubmatch(cmdout)
           nodegroups=rs[1]

        }

        return nodegroups,0,""
}

func releasenode(node string,user string)(int,string){
     nodegroups,errcode,errmsg:=getnodegroups(node)
     if errcode !=0 {
         return errcode,errmsg
     }

     nodegrouplist:=strings.Split(nodegroups,",")
     if !Contains(nodegrouplist, user) {
         return 0,""
     }

     cmd := exec.Command("chdef","-t","node","-o",node,"groups=free")
     var outbuf, errbuf bytes.Buffer
     cmd.Stdout = &outbuf
     cmd.Stderr = &errbuf
         
     err := cmd.Run()        
     if err != nil {
            return 1,"Failed to release node resource "+node+" from user "+user +": "+errbuf.String()
     }

     return 0,""
} 

func getattr(key string, d *schema.ResourceData) (int,string) {
     keylist:=strings.Split(key,".")
     key, keylist = keylist[0], keylist[1:]
    
     vRaw, ok := d.GetOk(key) 
     if !ok {
         return -1,"unexist key"
     } 


 
     for {
        if len(keylist)==0{
            return 0,vRaw.(string);
        }


        if typeof(vRaw) == "*schema.Set" {
           componentRaw := vRaw.(*schema.Set).List()
           raw:=componentRaw[0]
           rawMap := raw.(map[string]interface{})
           key, keylist = keylist[0], keylist[1:]
           vRaw=rawMap[key]
        } 
     }

}


func Intf2Map(v interface{}) map[string]string {
	m := v.(map[string]interface{})
        retmap:=make(map[string]string)
        for key,value :=range m{
            retmap[key]=value.(string)
        }
	return retmap
}



var DictRes2Inv = map[string]string{
    "mtm" : "device_info.mtm",
    "arch": "device_info.arch",
    "disksize":"device_info.disksize",
    "memory":"device_info.memory",
    "cputype":"device_info.cputype",
    "cpucount":"device_info.cpucount",
    "serial":"device_info.serial",
    "ip":"network_info.primarynic.ip",
    "mac":"network_info.primarynic.mac",
    "rack":"position_info.rack",
    "unit":"position_info.unit",
    "room":"position_info.room",
    "height":"position_info.height",
    "osimage":"engines.netboot_engine.engine_info.osimage",
    "zone":"security_info.zonename",
}


func Res2DefAttr(resattr string) string{
    if resattr == "machinetype" {
        return "mtm"
    }
    return resattr
}

func NodeInv2Res(myjson *gojsonq.JSONQ, d *schema.ResourceData,node string) int {
    keys := reflect.ValueOf(DictRes2Inv).MapKeys()
    for _, kres := range keys {
        kinv:=DictRes2Inv[kres.String()]
        val:=myjson.Reset().From("node."+node).Find(kinv)
        if val != nil {
           d.Set(kres.String(),val)
        } else {
           d.Set(kres.String(),nil)
        }
    }
    return 0
}



func RunCmd(cmdstr string,args ...string) (error,string,string) {
     cmd := exec.Command(cmdstr,args...)
     var outbuf, errbuf bytes.Buffer
     cmd.Stdout = &outbuf
     cmd.Stderr = &errbuf
     err := cmd.Run()        
     return err,outbuf.String(),errbuf.String()
}

type NetbootParam struct {
     osimage string
     addkcmdline string
}
func ProvisionNode(node string, param *NetbootParam) (int,string) {
     err,outstr,errstr:=RunCmd("makedns",node)
     err,outstr,errstr=RunCmd("rinstall",node,"osimage="+param.osimage)
     if err!=nil{
         return 1,errstr
     }    

     var myregex = regexp.MustCompile("status=booted")
     for {
         err,outstr,errstr=RunCmd("lsdef","-t","node","-o",node,"-i","status")
         if myregex.MatchString(outstr) { 
             return 0,""
         }
     }
     return 1,""
}
