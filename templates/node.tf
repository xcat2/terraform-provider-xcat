resource "xcat_node" "devnode" {
  selectors {
    cpucount="128"
    machinetype = "8335-GTC"
  }
  count=5
  osimage="rhels8.0-ppc64le-netboot-compute"
}

resource "xcat_node" "fvtnode" {
  selectors {
    cpucount="128"
    machinetype = "8335-GTC"
  }
  count=6
  osimage="rhels8.0-ppc64le-netboot-compute"
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
