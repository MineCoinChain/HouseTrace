package application

import (
	// "github.com/astaxie/beego"
	"web/models/fabricsetup"
	"fmt"
	"os"
)

var App *Application

type Application struct {
	//维护fabric-sdk对象
	FabricSetup *fabricsetup.FabricSetup
}

func NewApplication() (*Application, error){
	setup := fabricsetup.NewFabricSetup()

	//通道创建，添加，链码安装，初始化
	err := setup.Init()
	if err != nil  {
		return nil, err
	}

	App := Application{
		FabricSetup:setup,
	}
	return &App, nil
}

func init() {
	app, err := NewApplication()
	if err !=nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}

	App = app
	fmt.Println("Application init successfully!")
}