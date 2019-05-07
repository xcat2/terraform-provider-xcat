xCAT Terraform Provider
==================

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) v0.11.13

Installation
------------

## Download and install Terraform on xCAT MN

```sh
$ wget https://media.github.ibm.com/releases/207181/files/158261?token=AABUypPM6uPxEY5_rpIYtJiFjzxopYNWks5c0Tt7wA%3D%3D -O /usr/bin/terraform
$ chmod +x /usr/bin/terraform
```

## Download and install xCAT Terraform Provider on xCAT MN

```sh
$ wget https://media.github.ibm.com/releases/207181/files/158263?token=AABUyukEerLIW1PPyBj1yrwUdVNf1AxFks5c0TwdwA%3D%3D -O ~/.terraform.d/plugins/terraform-provider-xcat
$ chmod +x ~/.terraform.d/plugins/terraform-provider-xcat 
```

Creat node resource pool on xCAT MN
------------------------------------

```sh
$ chdef <xCAT nodes to be added into the pool> groups=free usercomment=","
```

Label the nodes in the resource pool on xCAT MN
-----------------------------------------------

Label the nodes with IB

```sh
$ chdef <xCAT nodes with IB> usercomment=",ib=1,"
```

Label the nodes with GPU

```sh
$ chdef <xCAT nodes with IB> usercomment=",gpu=1,"
```

Create Terraform working directory
----------------------------------

```sh
$ mkdir -p ~/mycluster/
$ cd ~/mycluster/
$ terraform init
```

Compose the cluster TF files
----------------------------

An example cluster TF files can be found in https://github.ibm.com/yangsbj/terraform-provider-xcat/tree/v0.1/templates/devcluster. Modify the TF files according to your need

Refer https://www.terraform.io/docs/configuration/index.html for the Terraform HCL syntax

Resource operation
------------------
**plan:

```sh
$ cd ~/mycluster/
$ terraform plan
```
 
**resource apply:

```sh
$ terraform apply
```

resource update:

modify the tf file and run
```sh
$ terraform apply
```

resource release:

```sh
$ terraform destroy
```
