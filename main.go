package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"adguard-tray/cli"

	"fyne.io/systray"
)

var (
	verboseLog *log.Logger
	adguardCLI *cli.AdGuardCLI

	// 1. Объявляем пункты меню как глобальные переменные
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

	icon, err := os.ReadFile("assets/icon_disconnected.png")
	if err != nil {
		log.Fatalf("Не удалось загрузить иконку: %v", err)
	}
	systray.SetIcon(icon)
	systray.SetTitle("AdGuard VPN Status")
	systray.SetTooltip("AdGuard VPN Status")

	// 2. Создаем все пункты меню при запуске
	mStatus = systray.AddMenuItem("Статус: Проверка...", "Текущий статус VPN")
	mStatus.Disable()

	systray.AddSeparator()

	mConnect = systray.AddMenuItem("Подключиться (лучшая локация)", "Подключиться к самому быстрому серверу")
	mDisconnect = systray.AddMenuItem("Отключиться", "Отключиться от VPN")

	// Изначально скрываем кнопку "Отключиться"
	mDisconnect.Hide()

	systray.AddSeparator()
	mQuit = systray.AddMenuItem("Выход", "Закрыть приложение")

	// 3. Запускаем обработчики кликов
	go handleClicks()

	// Запускаем цикл обновления статуса
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		updateStatusUI() // Первый запуск
		for range ticker.C {
			updateStatusUI()
		}
	}()

	verboseLog.Println("onReady: настройка трея завершена.")
}

// 4. Новая функция для обработки всех кликов в одном месте
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
			updateStatusUI() // Немедленно обновить UI
			mConnect.SetTitle("Подключиться (лучшая локация)")
			mConnect.Enable()

		case <-mDisconnect.ClickedCh:
			verboseLog.Println("Нажата кнопка 'Отключиться'")
			mDisconnect.SetTitle("Отключение...")
			mDisconnect.Disable()
			if err := adguardCLI.Disconnect(); err != nil {
				log.Printf("Ошибка при отключении: %v", err)
			}
			updateStatusUI() // Немедленно обновить UI
			mDisconnect.SetTitle("Отключиться")
			mDisconnect.Enable()

		case <-mQuit.ClickedCh:
			systray.Quit()
			return // Выход из цикла и горутины
		}
	}
}

// 5. Переименованная и обновленная функция для обновления UI
func updateStatusUI() {
	status, err := adguardCLI.GetStatus()
	if err != nil {
		errMsg := "Ошибка: статус неизвестен"
		mStatus.SetTitle(errMsg)
		log.Printf("%s: %v", errMsg, err)
		// В случае ошибки скрываем обе кнопки для безопасности
		mConnect.Hide()
		mDisconnect.Hide()
		return
	}

	if status.IsConnected {
		mStatus.SetTitle(fmt.Sprintf("Подключено: %s", status.Location))
		mConnect.Hide()
		mDisconnect.Show()
	} else {
		mStatus.SetTitle("Отключено")
		mConnect.Show()
		mDisconnect.Hide()
	}
}

func onExit() {
	log.Println("Приложение завершает работу.")
}
