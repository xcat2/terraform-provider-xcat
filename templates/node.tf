resource "xcat_node" "cn1" {
  name         = "cn1"
  obj_info {
     groups = "all"
     comments = "node for test" 
  }

  device_type = "server"
  device_info {
     arch = "ppc64le"
  }

  role = "compute"

  engines {
     netboot_engine {
        engine_type = "grub2"
        engine_info {
           osimage = "osimage1"
        }
     }
     hardware_mgt_engine {
        engine_type = "openbmc"
        engine_info {
           bmc = "bmc1"
           bmcusername = "USERID"
           bmcpassword = "PASSW0RD"
        }
     }
  }

  network_info {
     primarynic {
        ip = "10.3.5.102"
        mac="aa:bb:cc:dd:ee:ff"
     }
  }

}

  
