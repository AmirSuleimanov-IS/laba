-- Инициализация базы данных Kinetic Sports Registration

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица услуг (спортивные события/услуги)
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'deleted')),
    image_url VARCHAR(500),
    service_type VARCHAR(50) NOT NULL,
    event_date DATE,
    location VARCHAR(255),
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица заявок
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'deleted', 'formed', 'completed', 'rejected')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    formed_at TIMESTAMP,
    completed_at TIMESTAMP,
    moderator_id INTEGER REFERENCES users(id),
    team_name VARCHAR(255),
    total_amount DECIMAL(10, 2) DEFAULT 0,
    UNIQUE(user_id, status) -- Ограничение: один черновик на пользователя (частично реализуется через приложение)
);

-- Таблица связи заявки-услуги (M:M)
CREATE TABLE IF NOT EXISTS application_services (
    application_id INTEGER NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    service_id INTEGER NOT NULL REFERENCES services(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (application_id, service_id),
    UNIQUE(application_id, service_id)
);

-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_applications_user_status ON applications(user_id, status);
CREATE INDEX IF NOT EXISTS idx_applications_status ON applications(status);
CREATE INDEX IF NOT EXISTS idx_services_status ON services(status);
CREATE INDEX IF NOT EXISTS idx_application_services_application ON application_services(application_id);
CREATE INDEX IF NOT EXISTS idx_application_services_service ON application_services(service_id);

-- Вставка тестовых данных

-- Пользователи
INSERT INTO users (username, email, full_name) VALUES
('admin', 'admin@kinetic.com', 'Администратор Системы'),
('athlete1', 'volkov@kinetic.com', 'Алексей Волков'),
('athlete2', 'sokolov@kinetic.com', 'Александр Соколов'),
('athlete3', 'morozova@kinetic.com', 'Екатерина Морозова');

-- Услуги (события)
INSERT INTO services (name, description, status, image_url, service_type, event_date, location, price) VALUES
('Московский марафон 2026', 'Ежегодный городской марафон с дистанциями 42км, 21км, 10км', 'active', 'http://localhost:9000/sports/1.png', 'Марафон', '2026-05-15', 'Москва', 2500.00),
('Ironman Sochi', 'Полный триатлон: плавание 3.8км, велогонка 180км, бег 42км', 'active', 'http://localhost:9000/sports/2.png', 'Триатлон', '2026-06-20', 'Сочи', 15000.00),
('Ultra Trail Caucasus', 'Горный ультрамарафон 100км с набором высоты 6000м', 'active', 'http://localhost:9000/sports/3.png', 'Трейл', '2026-07-10', 'Красная Поляна', 8000.00),
('Велогонка Тур Сочи', 'Шоссейная велогонка 120км по побережью', 'active', 'http://localhost:9000/sports/4.png', 'Велогонка', '2026-08-05', 'Сочи', 3500.00),
('Зимний забег', 'Лыжная гонка 50км классическим стилем', 'active', 'http://localhost:9000/sports/5.png', 'Лыжная гонка', '2026-02-20', 'Красная Поляна', 2000.00),
('Спринт триатлон', 'Укороченная дистанция триатлона для новичков', 'deleted', NULL, 'Триатлон', '2025-12-01', 'Москва', 5000.00);

-- Заявка в статусе черновик для пользователя 1
INSERT INTO applications (user_id, status, team_name) VALUES
(1, 'draft', 'Kinetic Elite Draft');

-- Добавление услуги в черновик
INSERT INTO application_services (application_id, service_id, quantity, sort_order, is_primary) VALUES
(1, 1, 1, 1, TRUE);

COMMENT ON TABLE users IS 'Пользователи системы';
COMMENT ON TABLE services IS 'Спортивные услуги и события';
COMMENT ON TABLE applications IS 'Заявки пользователей на участие';
COMMENT ON TABLE application_services IS 'Связь заявок и услуг (многие-ко-многим)';

COMMENT ON COLUMN applications.status IS 'Статусы: draft(черновик), deleted(удалён), formed(сформирован), completed(завершён), rejected(отклонён)';
COMMENT ON COLUMN applications.total_amount IS 'Рассчитывается при завершении заявки как сумма всех услуг';
