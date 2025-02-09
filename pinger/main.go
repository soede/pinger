package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/go-ping/ping"
)

type Config struct {
	BackendURL   string
	PingInterval time.Duration
	PingTimeout  time.Duration
}

type PingResult struct {
	Status       string    `json:"container_status"`
	Addr         string    `json:"addr"`
	PingDuration float64   `json:"p_duration"`
	PingedAt     time.Time `json:"pinged_at"`
}

func loadConfig() Config {
	backendURL := "http://backend:8080/api/v1/setAll"

	pingInterval := 10 * time.Second
	if v := os.Getenv("PING_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			pingInterval = d
		} else {
			log.Printf("Ошибка парсинга PING_INTERVAL, используем значение по умолчанию: %v", pingInterval)
		}
	}

	pingTimeout := 2 * time.Second
	if v := os.Getenv("PING_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			pingTimeout = d
		} else {
			log.Printf("Ошибка парсинга PING_TIMEOUT, используем значение по умолчанию: %v", pingTimeout)
		}
	}

	return Config{
		BackendURL:   backendURL,
		PingInterval: pingInterval,
		PingTimeout:  pingTimeout,
	}
}

func sendPingResults(ctx context.Context, config Config, results []PingResult) error {
	jsonData, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("не удалось маршализовать данные пинга: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.BackendURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("не удалось создать HTTP запрос: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении HTTP запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный HTTP статус: %d", resp.StatusCode)
	}

	return nil
}

func pingIP(ctx context.Context, ip string, timeout time.Duration) (time.Duration, error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return 0, fmt.Errorf("не удалось создать пингер: %w", err)
	}
	pinger.SetPrivileged(true)

	pinger.Count = 3
	pinger.Timeout = timeout

	if err := pinger.Run(); err != nil {
		return 0, fmt.Errorf("ошибка при пинге %s: %w", ip, err)
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		return 0, fmt.Errorf("нет полученных ответов от %s", ip)
	}

	return stats.AvgRtt, nil
}

func tcpPing(ip string, port int, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к %s:%d: %w", ip, port, err)
	}
	defer conn.Close()
	return nil
}

func main() {
	config := loadConfig()
	log.Printf("Запуск pinger сервиса с конфигурацией: %+v", config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		log.Println("Получен сигнал прерывания, завершаем работу...")
		cancel()
	}()

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Ошибка создания Docker клиента: %v", err)
	}

	ticker := time.NewTicker(config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Завершение основного цикла пинга.")
			return
		case <-ticker.C:
			containers, err := dockerClient.ContainerList(ctx, container.ListOptions{All: true})
			if err != nil {
				log.Printf("Ошибка получения списка контейнеров: %v", err)
				continue
			}

			var wg sync.WaitGroup
			var results []PingResult

			// Для каждого контейнера получаем подробную информацию и IP-адреса.
			for _, container := range containers {
				contJSON, err := dockerClient.ContainerInspect(ctx, container.ID)
				if err != nil {
					log.Printf("Ошибка при получении деталей контейнера %s: %v", container.ID, err)
					continue
				}
				status := contJSON.State.Status

				// Обрабатываем все сети, в которых находится контейнер.
				for netName, network := range contJSON.NetworkSettings.Networks {
					ip := network.IPAddress
					if ip == "" {
						log.Printf("Контейнер %s в сети %s не имеет IP-адреса", container.ID, netName)
						continue
					}

					wg.Add(1)
					go func(ip string) {
						defer wg.Done()

						start := time.Now()
						avgRtt, err := pingIP(ctx, ip, config.PingTimeout)
						if err != nil {
							log.Printf("Ошибка при пинге через ICMP для %s: %v", ip, err)
							if err := tcpPing(ip, 80, config.PingTimeout); err != nil {
								log.Printf("Ошибка при пинге через TCP для %s: %v", ip, err)
								return
							}
						}
						duration := time.Since(start)

						result := PingResult{
							Addr:         ip,
							PingDuration: float64(avgRtt.Milliseconds()),
							PingedAt:     time.Now(),
							Status:       status,
						}

						// Добавляем результат в список
						results = append(results, result)
						log.Printf("Контейнер %s (IP: %s) пингован успешно: среднее время = %v (измерено за %v)", ip, avgRtt, duration)
					}(ip)
				}
			}

			wg.Wait()

			// Отправляем результаты в backend.
			if err := sendPingResults(ctx, config, results); err != nil {
				log.Printf("Ошибка отправки результатов пинга: %v", err)
			}
		}
	}
}
