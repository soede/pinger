# Pinger
Pinger – это сервис, который получает адреса docker-контейнеров и пингует их.

## Как запустить?
Убедитесь, что у вас установлены Docker и Docker Compose и введите:
```bash
git clone https://github.com/soede/pinger.git
cd pinger
docker compose up
```

## UI скриншоты
![Home page](./docs/media/img.png "Home page")

## Сервисы
1. **Сервер на Go (Golang 1.23+)**
   - Использует библиотеку net/http для обработки запросов. 
   - Написал сервер с новым http.NewServeMux(), поэтому нужна версия golang 1.23+.
   - Работает на порту 8080.
2. **Пингер**
    - Написан на Go
    - Использует библиотеку ping.
3. **PostgreSQL**
4. **Pinger-ui (Фронтенд)**
   - Написан на React (TypeScript) и собирается с помощью Vite. Взял UI из https://ant.design/
   - После сборки статические файлы копируются в контейнер NGINX
   - Работает по адресу 8080
5. **Migrate**
   - Запускает миграции при первом запуске
