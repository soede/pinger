# Указываем базовый образ
FROM node:lts

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и package-lock.json для установки зависимостей
COPY package.json package-lock.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем остальные файлы в рабочую директорию
COPY . .

# Сборка приложения
RUN npm run build

# Указываем базовый образ для NGINX
FROM nginx:alpine

# Копируем файлы сборки в папку NGINX
COPY --from=0 /app/dist /usr/share/nginx/html

# Копируем конфигурационный файл NGINX
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

# Указываем команду для запуска NGINX
CMD ["nginx", "-g", "daemon off;"]

# Открываем порт 80 для доступа к приложению
EXPOSE 80
