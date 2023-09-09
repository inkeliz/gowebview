package main

import
(
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"gioui.org/font/gofont"

	"github.com/diyism/goAndView"

	//libs used by app logic:
	"time"
	"io/ioutil"
	"net/http"
	"fmt"
	"strings"
	"strconv"
)

func main()
{	go func()
	{	w := app.NewWindow()
		if err := loop(w); err != nil
		{	log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error
{	th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	var wv  gowebview.WebView

	for
	{	e:=<-w.Events()
		switch e:=e.(type)
		{	case app.ViewEvent:
				if wv==nil
				{	go func()
					{	var err error
						wv,err = gowebview.New(&gowebview.Config{URL:"https://google.com/ncr", WindowConfig:&gowebview.WindowConfig{Title:"Hello World", Window:e.View, VM:app.JavaVM()}})
						if err!=nil {panic(err)}
						defer wv.Destroy()
						//go wv.Wakelock()
                        //go wv.Locktask()
						//go checkTimestamp(wv)
						wv.Run()
					}()
				}
			case *system.CommandEvent:
				if e.Type==system.CommandBack
				{	log.Println("===================back button hibernate==============", wv)
					e.Cancel = true
					//wv.SetVisibility(gowebview.VisibilityMinimized)
					//time.AfterFunc(0*time.Second, func(){	wv.Vibrate()},)
					//time.AfterFunc(3*time.Second, func(){	wv.Vibrate()},)
					//time.AfterFunc(6*time.Second, func(){	wv.Hibernate()},)
					go wv.Hibernate()
				}
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				l := material.H1(th, "Hello, gowebview")
				maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
				l.Color = maroon
				l.Alignment = text.Middle
				l.Layout(gtx)
				e.Frame(gtx.Ops)
		}
	}
}