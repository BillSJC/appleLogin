# appleLogin

A tool for `Sign In with Apple` REST API written in `Golang`

![logo here](logo.jpg)

### What it dose?

- Full support of JWT Authorization for Apple Developer
- Quick API to finish [Sign In with Apple REST API](https://developer.apple.com/documentation/signinwithapplerestapi) on Apple Document
- Tool to create callBackURL

### Install:

> go get github.com/BillSJC/appleLogin

### Test:

> go test

### Quick Start

```go
package main

import (
	"fmt"
	"github.com/BillSJC/appleLogin"
)

func main(){
	a := appleLogin.InitAppleConfig("123ABC456D",   //Team ID
		"com.example.pkg",  //Client ID (Service ID)
		"https://www.example.com/callback", //Callback URL
		"your Apple Key ID")    //Key ID
		
	//import cert
	err := a.LoadP8CertByFile("path to your p8 cert file")  //path to cert file
	//or you can load cert from a string
	err = a.LoadP8CertByByte([]byte("set your cert string here"))
	
	if err != nil {
		panic(err)
	}
	
	//create callback URL
	callbackURL := a.CreateCallbackURL("state here")
	fmt.Println(callbackURL)
	
	// ... some code to get Apple`s AuthorizationCode
	code := "xxxx"
	token,err := a.GetAppleToken(code,3600)
	if err != nil {
		panic(err)
	}
	fmt.Println(token)

}

```

### Functions

#### InitAppleConfig

#### LoadP8CertByFile

#### LoadP8CertByByte

#### CreateCallbackURL

#### GetAppleToken