package cli

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var ansiRegex = regexp.MustCompile("\x1b\\[[0-9;]*m")

type VpnStatus struct {
	IsConnected bool
	Location    string
}

type AdGuardCLI struct {
	ExecutablePath string
	logger         *log.Logger
}

func New(executablePath string, logger *log.Logger) *AdGuardCLI {
	if executablePath == "" {
		executablePath = "adguardvpn-cli"
	}
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	return &AdGuardCLI{
		ExecutablePath: executablePath,
		logger:         logger,
	}
}

func (c *AdGuardCLI) GetStatus() (*VpnStatus, error) {
	c.logger.Println("Вызов команды:", c.ExecutablePath, "status")
	cmd := exec.Command(c.ExecutablePath, "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"команда 'status' завершилась с ошибкой: %w. Вывод: %s",
			err,
			string(output),
		)
	}
	cleanOutput := ansiRegex.ReplaceAllString(string(output), "")
	c.logger.Printf("Получен очищенный вывод от CLI:\n---\n%s\n---", cleanOutput)
	status := c.parseStatus(cleanOutput)
	c.logger.Printf("Результат парсинга: %+v", status)
	return status, nil
}

// Connect теперь запускает команду в новом окне терминала.
func (c *AdGuardCLI) Connect() error {
	cliArgs := []string{c.ExecutablePath, "connect"}

	// Формируем команду: x-terminal-emulator -e adguardvpn-cli connect ...
	// Флаг -e говорит терминалу выполнить команду и закрыться после ее завершения.
	args := append([]string{"-e"}, cliArgs...)
	cmd := exec.Command("x-terminal-emulator", args...)

	c.logger.Println("Запуск в терминале:", cmd.String())
	// Используем Run() вместо CombinedOutput(), так как нам не нужен вывод,
	// а нужно просто запустить команду и дождаться ее завершения.
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"не удалось запустить команду в терминале: %w",
			err,
		)
	}
	c.logger.Println("Команда 'connect' в терминале завершена.")
	return nil
}

func (c *AdGuardCLI) Disconnect() error {
	c.logger.Println("Вызов команды:", c.ExecutablePath, "disconnect")
	cmd := exec.Command(c.ExecutablePath, "disconnect")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"команда 'disconnect' завершилась с ошибкой: %w. Вывод: %s",
			err,
			string(output),
		)
	}
	c.logger.Println("Команда 'disconnect' выполнена успешно.")
	return nil
}

func (c *AdGuardCLI) parseStatus(output string) *VpnStatus {
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return &VpnStatus{IsConnected: false}
	}
	firstLine := lines[0]
	if strings.HasPrefix(firstLine, "Connected to ") {
		endOfLocation := strings.Index(firstLine, " in ")
		if endOfLocation == -1 {
			return &VpnStatus{IsConnected: false}
		}
		location := firstLine[13:endOfLocation]
		return &VpnStatus{IsConnected: true, Location: location}
	}
	return &VpnStatus{IsConnected: false}
}
