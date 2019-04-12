package xcat

import (
	"fmt"
	"log"
	"strings"
	"sync"
        "os"
        "os/exec"
        "bytes"
        "regexp"
        "encoding/json"
        "reflect"
        "github.com/jeremywohl/flatten"
        
	"github.com/hashicorp/terraform/helper/hashcode"
	//"github.com/hashicorp/terraform/helper/resource"
	//"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
        //"github.com/jeremywohl/flatten"
)

var systemSyncLock sync.Mutex

func resourceNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceNodeCreate,
		Read:   resourceNodeRead,
		Update: resourceNodeUpdate,
		Delete: resourceNodeDelete,

		Schema: map[string]*schema.Schema{
                        "name": {
                                  Type:     schema.TypeString,
                                  Optional: true,
                                  Computed: true,
                        },
			"obj_info": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"groups": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
                                //Set: resourceHash,
			},

			"device_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, name string) (warn []string, errs []error) {
					s := v.(string)
					validvalues := []string{"switch", "pdu", "rack", "hmc", "server"}
					if !Contains(validvalues, s) {
						errs = append(errs, fmt.Errorf("%s: the valid values: %s", name, strings.Join(validvalues, ",")))
					}
					return
				},
			},

			"device_info": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arch": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, name string) (warn []string, errs []error) {
								s := v.(string)
								validvalues := []string{"ppc64", "ppc64el", "ppc64le", "x86_64", "armv7l", "armel"}
								if !Contains(validvalues, s) {
									errs = append(errs, fmt.Errorf("%s: the valid values: %s", name, strings.Join(validvalues, ",")))
								}
								return
							},
						},
					},
				},
			},

			"role": {
				Type:     schema.TypeString,
				Required: false,
                                Computed: true,
                                /*
				ValidateFunc: func(v interface{}, name string) (warn []string, errs []error) {
					s := v.(string)
					validvalues := []string{"compute", "service"}
					if !Contains(validvalues, s) {
						errs = append(errs, fmt.Errorf("%s: the valid values: %s", name, strings.Join(validvalues, ",")))
					}
					return
				},*/
			},

			"engines": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"netboot_engine": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"engine_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: func(v interface{}, name string) (warn []string, errs []error) {
											s := v.(string)
											validvalues := []string{"pxe", "xnba", "grub2", "yaboot", "petitboot", "onie"}
											if !Contains(validvalues, s) {
												errs = append(errs, fmt.Errorf("%s: the valid values: %s", name, strings.Join(validvalues, ",")))
											}
											return
										},
									},
									"engine_info": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"osimage": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"hardware_mgt_engine": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"engine_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: func(v interface{}, name string) (warn []string, errs []error) {
											s := v.(string)
											validvalues := []string{"openbmc", "ipmi", "hmc", "fsp", "kvm", "mp", "bpa", "ivm", "blade", "pdu", "switch"}
											if !Contains(validvalues, s) {
												errs = append(errs, fmt.Errorf("%s: the valid values: %s", name, strings.Join(validvalues, ",")))
											}
											return
										},
									},
									"engine_info": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bmc": {
													Type:     schema.TypeString,
													Required: true,
												},
												"bmcusername": {
													Type:     schema.TypeString,
													Required: true,
												},
												"bmcpassword": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"network_info": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primarynic": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Required: true,
									},
									"mac": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceNodeCreate(d *schema.ResourceData, meta interface{}) error {
	systemSyncLock.Lock()
	defer systemSyncLock.Unlock()

        //flat, err := flatten.Flatten(d, "", flatten.DotStyle)
	config := meta.(*Config)
        fileName:= "/tmp/log_debug.log"
        logFile,err := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
        if err != nil {
           log.Fatalln("open file error!")
        }
        defer logFile.Close()
        debugLog := log.New(logFile,"[Debug]",log.Llongfile)
        debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)

        //debugLog.Printf("+%v\n",flat)
        debugLog.Printf("checking whether the resource '%s' exists \n",d.Get("name"))

        /*
        debugLog.Printf("type =%s\n",typeof(d.Get("engines")))
        vRaw:=d.Get("engines")
        componentRaw := vRaw.(*schema.Set).List()
        debugLog.Printf("xxxxxxxxx%d VVVvv",len(componentRaw))
        for i, raw := range componentRaw {
            rawMap := raw.(map[string]interface{})
            debugLog.Printf("type =%d,%s\n",i,typeof(rawMap["netboot_engine"]))
            
        }
        */
        node:=d.Get("name").(string)

        retcode,retv:=getattr("engines.netboot_engine.engine_type", d)
        debugLog.Printf("%s xxxxxxx\n",retv)
        if retcode !=0{
           return fmt.Errorf("failed to get key %s","engines.netboot_engine.engine_type")
        }
         


        nodegroups,errcode,errmessage := getnodegroups(node) 
        if errcode!=0 {
            return fmt.Errorf(errmessage )
        }

        debugLog.Printf("nodegroups=%s\n",nodegroups)              

        username:=config.Username
        debugLog.Printf("username=%s\n",username ) 
        errcode,errmessage= occupynode(node, username)
        if errcode!=0 {
            return fmt.Errorf(errmessage )
        }

        
        //debugLog.Printf("engines=%s\n",d.Get("engines").Get("netboot_engine").Get("engine_type").(string)) 
        //debugLog.Printf("engines=%s\n",d.Get("engines").([]interface{}).Get("netboot_engine").([]interface{}).Get("engine_type").([]interface{}).(string)) 
 

        debugLog.SetPrefix("[Info]")

        log.Printf("[INFO] there is a pending resize operation on this pool...")

	return resourceNodeRead(d, meta)
}

func resourceNodeRead(d *schema.ResourceData, meta interface{}) error {
	//config := meta.(*Config)
        fileName:= "/tmp/log_debug.log"
        logFile,err := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
        if err != nil {
           log.Fatalln("open file error!")
        }
        defer logFile.Close()
        debugLog := log.New(logFile,"[Debug]",log.Llongfile)
        debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)



        node:=d.Get("name").(string)
        cmd := exec.Command("xcat-inventory","export","-t","node","-o",node,"--format","json")
        var outbuf, errbuf bytes.Buffer
        cmd.Stdout = &outbuf
        cmd.Stderr = &errbuf
            
        err = cmd.Run()        

        if err != nil {
               debugLog.Printf("Failed to read node resource "+node+" from xcat: "+errbuf.String())
        }


        debugLog.Printf(outbuf.String())        
        debugLog.Printf(errbuf.String())        

        nodejson :=outbuf.String()
         
        nodemap := make(map[string]interface{})
        err = json.Unmarshal([]byte(nodejson), &nodemap)
        debugLog.Printf("%v",nodemap)

        keys := reflect.ValueOf(nodemap["node"]).MapKeys()[0]
        debugLog.Printf("%v",keys)
        nodename:=keys.String()
      
        flattened,err:=flatten.Flatten(nodemap["node"].(map[string]interface{})[nodename].(map[string]interface{}),"",flatten.DotStyle) 
        debugLog.Printf("%v",flattened)
 
        p_obj_info:=&schema.Set{F:resourceHash}

        obj_info := map[string]interface{}{}

        if val,ok := flattened["obj_info.groups"];ok{
            obj_info["groups"]=val
        }

        if val,ok := flattened["obj_info.description"];ok{
            obj_info["description"]=val
        }
        debugLog.Printf("XXXXXXXXXX\n") 
        p_obj_info.Add(obj_info)
        debugLog.Printf("YYYYYYYYYYY\n") 
        //d.Set("obj_info", p_obj_info)
        d.Set("role","XXXXXXXXXX")

        debugLog.Printf("XXXXXXXXXX%v\n",d) 
        
        //nodename:=nodemap["node"] 

	return nil
}

func resourceNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	systemSyncLock.Lock()
	defer systemSyncLock.Unlock()

	//config := meta.(*Config)
        log.Printf("[DEBUG] XXXXXXXXXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@...")

	return resourceNodeRead(d, meta)
}

func resourceNodeDelete(d *schema.ResourceData, meta interface{}) error {

        fileName:= "/tmp/log_debug.log"
        logFile,err := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
        if err != nil {
           log.Fatalln("open file error!")
        }
        defer logFile.Close()
        debugLog := log.New(logFile,"[Debug]",log.Llongfile)
        debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)

	config := meta.(*Config)
        username:=config.Username
        node:=d.Get("name").(string)
        debugLog.Printf("in resourceNodeDelete\n")
        errorcode,errormessage := releasenode(node,username)
        if errorcode!=0 {
            return fmt.Errorf(errormessage )
        }
	return nil
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

func occupynode(node string, user string) (int,string) {
     nodegroups,errcode,errmsg:=getnodegroups(node)
     if errcode !=0 {
         return errcode,errmsg
     }

     nodegrouplist:=strings.Split(nodegroups,",")
     if Contains(nodegrouplist, user) {
         return 0,""
     }


     if user != "root" {
         for _,v := range nodegrouplist {
             if v!="root"{
                 return 1, "node resource "+node+" has already been occupier by user "+user
             }
         }
     }

     cmd := exec.Command("chdef","-t","node","-o",node,"-p","groups="+user)
     var outbuf, errbuf bytes.Buffer
     cmd.Stdout = &outbuf
     cmd.Stderr = &errbuf
         
     err := cmd.Run()        
     if err != nil {
            return 1,"Failed to occupy node resource "+node+" for user "+user +": "+errbuf.String()
     }
     return 0,""
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

     cmd := exec.Command("chdef","-t","node","-o",node,"-m","groups="+user)
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


func resourceHash(v interface{}) int {
	m := v.(map[string]interface{})
        
        mapstr:=""
        for key,value :=range m{
            mapstr=mapstr+","+key+":"+value.(string)
        }
	return hashcode.String(mapstr)
}

