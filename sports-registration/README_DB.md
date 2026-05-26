# Kinetic Sports Registration - Инструкция по запуску с PostgreSQL

## 📋 Описание

Приложение для управления спортивными командами и заявками на события с использованием:
- **Go** (backend)
- **PostgreSQL** (база данных)
- **Minio** (хранилище файлов)
- **Adminer** (управление БД)

## 🗄️ Структура базы данных

### Таблицы (4 таблицы, каскадное удаление запрещено):

1. **users** - пользователи системы
   - id, username, email, full_name, created_at

2. **services** - услуги/события
   - id, name, description, status (active/deleted), image_url, service_type, event_date, location, price, created_at

3. **applications** - заявки
   - id, user_id, status (draft/deleted/formed/completed/rejected), created_at, formed_at, completed_at, moderator_id, team_name, total_amount

4. **application_services** - связь заявки-услуги (M:M)
   - application_id, service_id (составной ключ), quantity, sort_order, is_primary, added_at

### Статусы заявок (5 статусов):
- `draft` - черновик
- `deleted` - удалён
- `formed` - сформирован
- `completed` - завершён
- `rejected` - отклонён

## 🚀 Запуск приложения

### Вариант 1: Docker Compose (рекомендуется)

```bash
cd /workspace/sports-registration

# Запуск всех сервисов
docker compose up -d

# Проверка логов
docker compose logs -f app
```

**Адреса:**
- Приложение: http://localhost:8080
- Adminer (БД): http://localhost:8081
- Minio Console: http://localhost:9001 (minioadmin/minioadmin)

### Вариант 2: Локальный запуск (без Docker)

#### 1. Установка PostgreSQL

```bash
# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib

# macOS (Homebrew)
brew install postgresql@15
brew services start postgresql@15
```

#### 2. Создание базы данных

```bash
sudo -u postgres psql -c "CREATE USER kinetic_user WITH PASSWORD 'kinetic_pass';"
sudo -u postgres psql -c "CREATE DATABASE kinetic_db OWNER kinetic_user;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE kinetic_db TO kinetic_user;"
```

#### 3. Инициализация схемы БД

```bash
PGPASSWORD=kinetic_pass psql -h localhost -U kinetic_user -d kinetic_db -f init.sql
```

Или через Adminer:
1. Откройте http://localhost:8081
2. Система: PostgreSQL, Сервер: localhost, Пользователь: kinetic_user, Пароль: kinetic_pass, База: kinetic_db
3. Выполните SQL из файла `init.sql` во вкладке "SQL command"

#### 4. Запуск приложения

```bash
cd /workspace/sports-registration

# Установка переменных окружения
export DB_HOST=localhost
export DB_USER=kinetic_user
export DB_PASS=kinetic_pass
export DB_NAME=kinetic_db

# Сборка и запуск
go build -o main .
./main
```

## 🔧 Настройка Minio

1. Откройте http://localhost:9001
2. Логин: `minioadmin`, Пароль: `minioadmin`
3. Создайте bucket `sports`
4. Установите Access Policy → Public
5. Загрузите файлы:
   - `1.png`, `2.png`, `3.png`, `4.png`, `5.png` - изображения услуг
   - `video.mp4` - видео для главной страницы

## 📱 HTTP Методы (5 методов)

### GET запросы (3):
1. **GET /services** - получение и поиск услуг
2. **GET /application** - просмотр текущей заявки (корзины)
3. **GET /** - главная страница

### POST запросы (2):
1. **POST /application/add** - добавление услуги в заявку (через ORM)
   - Параметр: `service_id`
   
2. **POST /application/delete** - удаление заявки (через SQL UPDATE без ORM)
   - Параметр: `application_id`
   - Выполняет: `UPDATE applications SET status = 'deleted' WHERE id = $1`

## ✅ Требования к БД выполнены

- ✅ 4 таблицы (users, services, applications, application_services)
- ✅ Каскадное удаление запрещено (ON DELETE CASCADE только для application_services)
- ✅ 5+ статусов заявок (draft, deleted, formed, completed, rejected)
- ✅ Один черновик на пользователя
- ✅ M:M связь с доп. полями (quantity, sort_order, is_primary)
- ✅ Поле total_amount рассчитывается при завершении
- ✅ Логическое удаление через UPDATE (без ORM)
- ✅ Добавление через ORM-методы

## 🔍 Проверка работы

1. Откройте http://localhost:8080/services
2. При первом посещении будет сообщение о необходимости создания заявки
3. Добавьте услугу в заявку (кнопка "Добавить в заявку")
4. Перейдите в заявку (http://localhost:8080/application)
5. Удалите заявку через кнопку "Удалить заявку" (выполняется SQL UPDATE)
6. После удаления заявка не отображается

## 📊 Данные для входа в Adminer

- **Система**: PostgreSQL
- **Сервер**: db (в Docker) или localhost (локально)
- **Пользователь**: kinetic_user
- **Пароль**: kinetic_pass
- **База данных**: kinetic_db

## 🛠️ Переменные окружения

| Переменная | Значение по умолчанию | Описание |
|------------|----------------------|----------|
| DB_HOST | localhost | Хост PostgreSQL |
| DB_USER | kinetic_user | Пользователь БД |
| DB_PASS | kinetic_pass | Пароль БД |
| DB_NAME | kinetic_db | Имя базы данных |

## 💡 Примечания

- Без подключенной PostgreSQL приложение работает в режиме "только статические данные"
- У каждого пользователя может быть только одна заявка в статусе `draft`
- Удаленные заявки (`status = 'deleted'`) не отображаются в интерфейсе
- Поле `total_amount` рассчитывается автоматически при завершении заявки
