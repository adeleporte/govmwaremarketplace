# VMware MarketPlace Go Client

This client features a way to list and download files from VMware MarketPlace website.


## Contents

- [VMware MarketPlace Go client](#vmware-marketplace-go-client)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Quick start](#quick-start)


## Installation

To install VMware MarketPlace go client, you need to install Go and set your Go workspace first.

1. The first need [Go](https://golang.org/) installed (**version 1.12+ is required**), then you can use the below Go command to install Gin.

```sh
$ go get -u github.com/adeleporte/govmaremarketplace
```

2. Import it in your code:

```go
import "github.com/adeleporte/govmaremarketplace"
```


## Quick start

```sh
# assume the following codes in example.go file
$ cat example.go
```

```go
package main

import (
    market "github.com/adeleporte/govmaremarketplace"
)



func main() {
    // Create a VMWare Cloud Session
    client, err := market.NewClient("mytoken")
    if err != nil {
        log.Println(err.Error())
        return
    }

    // Get Product List
    pl, err := market.GetProducts(client)
    if err != nil {
        log.Println(err.Error())
        return
    }

    for _, p := range pl.Response.DataList {
        log.Printf("Product Display Name is %+v. Slug is %+v\n", p.DisplayName, p.Slug)
    }

    // Get Fortinet Product Detail
    pd, err := market.GetProductDetail(client, "fortigate-next-generation-firewall-1")
    if err != nil {
        log.Println(err.Error())
        return
    }
    log.Printf("Product Details are: %+v", pd)

    // Get download link info
    dl, err := market.GetDownload(client, DownloadBody{DeploymentFileID: pd.Response.Data.ProductDeploymentFilesList[0].FileID, ProductID: pd.Response.Data.ProductID})
    if err != nil {
        log.Println(err.Error())
        return
    }
    log.Printf("Link is: %+v", dl.Response.PresignedURL)

// Download the file
err = market.Download(client, pd.Response.Data.ProductDeploymentFilesList[0].Name, dl.Response.PresignedURL)
    if err != nil {
        log.Println(err.Error())
    }

}
```

```
# run example.go

adeleporte@adeleporte-a02 govmwaremarketplace % go run .
Downloading... 71 MB complete      
adeleporte@adeleporte-a02 govmwaremarketplace % 

```
