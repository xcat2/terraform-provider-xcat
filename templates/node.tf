resource "xcat_node" "devnode" {
  selectors {
    //cpucount="128"
    //machinetype = "8335-GTC"
    arch="ppc64le"
  }
  count=1
  osimage="rhels7.4-ppc64le-netboot-compute"
}

resource "xcat_node" "fvtnode" {
  selectors {
    //cpucount="128"
    //machinetype = "8335-GTC"
    arch="ppc64le"
  }
  count=3
  osimage="rhels7.4-ppc64le-netboot-compute"
  //osimage="rhels8.1-ppc64le-netboot-compute"
}
/*
resource "xcat_node" "node1" {
  name="mid08tor03cn01"
}
*/


output "devnodes" {
  value=[ 
      "${xcat_node.devnode.*.name}"
  ]
}

output "fvtnodes" {
  value=[ 
      "${xcat_node.fvtnode.*.name}"
  ]
}
