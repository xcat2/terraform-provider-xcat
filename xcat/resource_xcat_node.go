package xcat

import (
	"fmt"
	//"log"
	"strings"
	"sync"

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
			"obj_info": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"groups": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"device_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "server",
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
							Default:  "x86_64",
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
				Required: true,
				Default:  "compute",
				ValidateFunc: func(v interface{}, name string) (warn []string, errs []error) {
					s := v.(string)
					validvalues := []string{"compute", "service"}
					if !Contains(validvalues, s) {
						errs = append(errs, fmt.Errorf("%s: the valid values: %s", name, strings.Join(validvalues, ",")))
					}
					return
				},
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
										Default:  "grub2",
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
							Optional: true,
							Default:  "ipmi",
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

	//config := meta.(*Config)

	return resourceNodeRead(d, meta)
}

func resourceNodeRead(d *schema.ResourceData, meta interface{}) error {
	//config := meta.(*Config)

	return nil
}

func resourceNodeUpdate(d *schema.ResourceData, meta interface{}) error {
	systemSyncLock.Lock()
	defer systemSyncLock.Unlock()

	//config := meta.(*Config)

	return resourceNodeRead(d, meta)
}

func resourceNodeDelete(d *schema.ResourceData, meta interface{}) error {
	//config := meta.(*Config)

	return nil
}
