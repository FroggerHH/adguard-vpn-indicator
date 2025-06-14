package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"adguard-tray/cli"

	"fyne.io/systray"
)

//go:embed assets/icon_connected.png
var iconConnected []byte

//go:embed assets/icon_disconnected.png
var iconDisconnected []byte

var (
	verboseLog *log.Logger
	adguardCLI *cli.AdGuardCLI

	mStatus     *systray.MenuItem
	mConnect    *systray.MenuItem
	mQuit       *systray.MenuItem
	mDisconnect *systray.MenuItem
)

func main() {
	verbose := flag.Bool("v", false, "Включить подробное логирование")
	flag.Parse()

	if *verbose {
		verboseLog = log.New(os.Stdout, "VERBOSE: ", log.LstdFlags)
		verboseLog.Println("Подробное логирование включено.")
	} else {
		verboseLog = log.New(io.Discard, "", 0)
	}

	adguardCLI = cli.New("adguardvpn-cli", verboseLog)

	verboseLog.Println("Запуск приложения...")
	systray.Run(onReady, onExit)
}

func onReady() {
	verboseLog.Println("onReady: начало настройки трея.")

	systray.SetIcon(iconDisconnected)
	systray.SetTitle("AdGuard VPN Status")
	systray.SetTooltip("AdGuard VPN Status")

	mStatus = systray.AddMenuItem("Статус: Проверка...", "Текущий статус VPN")
	mStatus.Disable()

	systray.AddSeparator()

	mConnect = systray.AddMenuItem("Подключиться (лучшая локация)", "Подключиться к самому быстрому серверу")
	mDisconnect = systray.AddMenuItem("Отключиться", "Отключиться от VPN")
	mDisconnect.Hide()

	systray.AddSeparator()
	mQuit = systray.AddMenuItem("Выход", "Закрыть приложение")

	go handleClicks()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		updateStatusUI()
		for range ticker.C {
			updateStatusUI()
		}
	}()

	verboseLog.Println("onReady: настройка трея завершена.")
}

func handleClicks() {
	for {
		select {
		case <-mConnect.ClickedCh:
			verboseLog.Println("Нажата кнопка 'Подключиться'")
			mConnect.SetTitle("Подключение...")
			mConnect.Disable()
			if err := adguardCLI.Connect(); err != nil {
				log.Printf("Ошибка при подключении: %v", err)
			}
			updateStatusUI()
			mConnect.SetTitle("Подключиться (лучшая локация)")
			mConnect.Enable()

		case <-mDisconnect.ClickedCh:
			verboseLog.Println("Нажата кнопка 'Отключиться'")
			mDisconnect.SetTitle("Отключение...")
			mDisconnect.Disable()
			if err := adguardCLI.Disconnect(); err != nil {
				log.Printf("Ошибка при отключении: %v", err)
			}
			updateStatusUI()
			mDisconnect.SetTitle("Отключиться")
			mDisconnect.Enable()

		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// 4. Добавляем смену иконок в функцию обновления UI.
func updateStatusUI() {
	status, err := adguardCLI.GetStatus()
	if err != nil {
		errMsg := "Ошибка: статус неизвестен"
		mStatus.SetTitle(errMsg)
		log.Printf("%s: %v", errMsg, err)

		systray.SetIcon(iconDisconnected)
		mConnect.Hide()
		mDisconnect.Hide()
		return
	}

	if status.IsConnected {
		mStatus.SetTitle(fmt.Sprintf("Подключено: %s", status.Location))
		systray.SetIcon(iconConnected)
		mConnect.Hide()
		mDisconnect.Show()
	} else {
		mStatus.SetTitle("Отключено")
		systray.SetIcon(iconDisconnected)
		mConnect.Show()
		mDisconnect.Hide()
	}
}

func onExit() {
	log.Println("Приложение завершает работу.")
}
