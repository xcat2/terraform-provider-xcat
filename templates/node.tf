resource "xcat_node" "node1" {
  name="mid08tor03cn01"
  #name="c910f03c05k27"
  
  machinetype = "server"
  arch = "ppc64le"
  disksize="200G"
  memory="200G"
  cpucount="4"
  
  osimage="rhels8.0-ppc64le-netboot-compute"

}

  
resource "xcat_node" "newnode" {
  #name="mid08tor03cn01"
  name="node000${count.index}"
  count=5
}
