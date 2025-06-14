package main

import (
	"flag" // 1. Импортируем пакет для работы с флагами
	"io"   // 2. Импортируем пакет io для io.Discard
	"log"
	"os"

	"fyne.io/systray"
)

// 3. Создаем глобальную переменную для нашего "verbose" логгера.
var verboseLog *log.Logger

func main() {
	// 4. Определяем флаг командной строки "-v".
	// flag.Bool(имя, значение по умолчанию, описание)
	verbose := flag.Bool("v", false, "Включить подробное логирование")
	flag.Parse() // Анализируем переданные при запуске флаги

	// 5. Настраиваем наш логгер в зависимости от флага.
	if *verbose {
		// Если флаг -v установлен, логгер будет писать в стандартный вывод (консоль).
		verboseLog = log.New(os.Stdout, "VERBOSE: ", log.LstdFlags)
		verboseLog.Println("Подробное логирование включено.")
	} else {
		// В противном случае, логгер будет писать в "никуда" (io.Discard).
		// Это элегантный способ "отключить" логирование без кучи if-проверок в коде.
		verboseLog = log.New(io.Discard, "", 0)
	}

	verboseLog.Println("Запуск приложения...")
	systray.Run(onReady, onExit)
}

// onReady настраивает иконку и меню в трее.
func onReady() {
	verboseLog.Println("onReady: начало настройки трея.")

	icon, err := os.ReadFile("assets/icon_disconnected.png")
	if err != nil {
		// Фатальные ошибки, которые мешают запуску, логируем всегда.
		log.Fatalf("Не удалось загрузить иконку: %v", err)
	}
	systray.SetIcon(icon)
	systray.SetTitle("AdGuard VPN Status")
	systray.SetTooltip("AdGuard VPN Status")
	verboseLog.Println("onReady: иконка и заголовок установлены.")

	mStatus := systray.AddMenuItem("Статус: Проверка...", "Текущий статус VPN")
	mStatus.Disable()

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Выход", "Закрыть приложение")
	verboseLog.Println("onReady: пункты меню созданы.")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	verboseLog.Println("onReady: настройка трея завершена.")
}

// onExit выполняется перед закрытием приложения.
func onExit() {
	// Сообщение о выходе логируем всегда.
	log.Println("Приложение завершает работу.")
}
