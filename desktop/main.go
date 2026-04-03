package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/axeprpr/n2nGUI/internal/app"
	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend/*
var assets embed.FS

type statusResponse struct {
	Running bool `json:"running"`
}

type apiErrorResponse struct {
	Error string `json:"error"`
}

func callAPI(api http.Handler, method, path string, out any) error {
	req := httptest.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	if recorder.Code >= 200 && recorder.Code < 300 {
		if out != nil {
			if err := json.NewDecoder(recorder.Body).Decode(out); err != nil {
				return fmt.Errorf("decode response failed: %w", err)
			}
		}
		return nil
	}

	var errResp apiErrorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&errResp); err == nil && errResp.Error != "" {
		return errors.New(errResp.Error)
	}
	return fmt.Errorf("request %s %s failed: %d", method, path, recorder.Code)
}

func main() {
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	baseDir, err = filepath.Abs(filepath.Join(baseDir, ".."))
	if err != nil {
		log.Fatal(err)
	}

	server := app.NewServer(baseDir)
	apiHandler := server.APIHandler()
	var startupCtx context.Context

	go systray.Run(func() {
		systray.SetTitle("n2nGUI")
		systray.SetTooltip("n2nGUI - 运行中")

		statusItem := systray.AddMenuItem("状态: 停止", "当前 n2n 状态")
		statusItem.Disable()
		systray.AddSeparator()
		startItem := systray.AddMenuItem("启动 n2n", "启动 n2n edge")
		stopItem := systray.AddMenuItem("停止 n2n", "停止 n2n edge")
		showItem := systray.AddMenuItem("显示窗口", "显示主窗口")
		quitItem := systray.AddMenuItem("退出", "退出应用")

		updateStatus := func() {
			var status statusResponse
			if err := callAPI(apiHandler, http.MethodGet, "/api/status", &status); err != nil {
				log.Printf("tray status update failed: %v", err)
				statusItem.SetTitle("状态: 异常")
				startItem.Enable()
				stopItem.Enable()
				return
			}

			if status.Running {
				statusItem.SetTitle("状态: 运行中")
				startItem.Disable()
				stopItem.Enable()
			} else {
				statusItem.SetTitle("状态: 停止")
				startItem.Enable()
				stopItem.Disable()
			}
		}

		updateStatus()

		go func() {
			for range startItem.ClickedCh {
				if err := callAPI(apiHandler, http.MethodPost, "/api/control/start", nil); err != nil {
					log.Printf("tray start failed: %v", err)
				}
				updateStatus()
			}
		}()

		go func() {
			for range stopItem.ClickedCh {
				if err := callAPI(apiHandler, http.MethodPost, "/api/control/stop", nil); err != nil {
					log.Printf("tray stop failed: %v", err)
				}
				updateStatus()
			}
		}()

		go func() {
			for range showItem.ClickedCh {
				if startupCtx == nil {
					continue
				}
				runtime.WindowUnminimise(startupCtx)
				runtime.WindowShow(startupCtx)
			}
		}()

		go func() {
			for range quitItem.ClickedCh {
				if startupCtx != nil {
					runtime.Quit(startupCtx)
					return
				}
				systray.Quit()
				os.Exit(0)
			}
		}()
	}, func() {})

	err = wails.Run(&options.App{
		Title:         "n2nGUI",
		Width:         480,
		Height:        640,
		MinWidth:      400,
		MinHeight:     500,
		DisableResize: false,
		Frameless:     false,
		OnStartup: func(ctx context.Context) {
			startupCtx = ctx
		},
		OnShutdown: func(ctx context.Context) {
			systray.Quit()
		},
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: server.APIHandler(),
		},
		BackgroundColour: &options.RGBA{R: 7, G: 17, B: 27, A: 1},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			IsZoomControlEnabled: false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
