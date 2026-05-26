package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sports-registration/models"

	_ "github.com/lib/pq"
)

type DatabaseRepository struct {
	db *sql.DB
}

func NewDatabaseRepository() (*DatabaseRepository, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "kinetic_user"
	}
	password := os.Getenv("DB_PASS")
	if password == "" {
		password = "kinetic_pass"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "kinetic_db"
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("✅ Подключение к PostgreSQL успешно")
	return &DatabaseRepository{db: db}, nil
}

func (r *DatabaseRepository) Close() error {
	return r.db.Close()
}

// GetAllServices - получение всех активных услуг (ORM не требуется, простой SQL)
func (r *DatabaseRepository) GetAllServices() ([]models.Service, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description, status, image_url, service_type, 
		       event_date, location, price, created_at 
		FROM services 
		WHERE status = 'active' 
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		var eventDate sql.NullString
		var imageURL sql.NullString

		err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Status, &imageURL,
			&s.ServiceType, &eventDate, &s.Location, &s.Price, &s.CreatedAt)
		if err != nil {
			return nil, err
		}

		if eventDate.Valid {
			s.EventDate = &eventDate.String
		}
		if imageURL.Valid {
			s.ImageURL = &imageURL.String
		}

		services = append(services, s)
	}

	return services, rows.Err()
}

// GetServiceByID - получение услуги по ID
func (r *DatabaseRepository) GetServiceByID(id int) (*models.Service, error) {
	row := r.db.QueryRow(`
		SELECT id, name, description, status, image_url, service_type, 
		       event_date, location, price, created_at 
		FROM services 
		WHERE id = $1 AND status = 'active'
	`, id)

	var s models.Service
	var eventDate sql.NullString
	var imageURL sql.NullString

	err := row.Scan(&s.ID, &s.Name, &s.Description, &s.Status, &imageURL,
		&s.ServiceType, &eventDate, &s.Location, &s.Price, &s.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if eventDate.Valid {
		s.EventDate = &eventDate.String
	}
	if imageURL.Valid {
		s.ImageURL = &imageURL.String
	}

	return &s, nil
}

// GetOrCreateDraftApplication - получить или создать заявку в статусе черновик для пользователя
// Использует ORM-подобный подход через sqlx-style queries
func (r *DatabaseRepository) GetOrCreateDraftApplication(userID int) (*models.Application, error) {
	// Сначала пытаемся найти существующий черновик
	row := r.db.QueryRow(`
		SELECT id, user_id, status, created_at, formed_at, completed_at, 
		       moderator_id, team_name, total_amount 
		FROM applications 
		WHERE user_id = $1 AND status = 'draft'
	`, userID)

	var app models.Application
	var formedAt, completedAt sql.NullString
	var moderatorID sql.NullInt64

	err := row.Scan(&app.ID, &app.UserID, &app.Status, &app.CreatedAt,
		&formedAt, &completedAt, &moderatorID, &app.TeamName, &app.TotalAmount)

	if err == nil {
		// Черновик найден, загружаем услуги
		services, err := r.GetApplicationServices(app.ID)
		if err != nil {
			return nil, err
		}
		app.Services = services
		return &app, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Черновик не найден, создаем новый
	teamName := fmt.Sprintf("Команда пользователя %d", userID)
	err = r.db.QueryRow(`
		INSERT INTO applications (user_id, status, team_name, created_at) 
		VALUES ($1, 'draft', $2, NOW()) 
		RETURNING id, user_id, status, created_at, formed_at, completed_at, 
		          moderator_id, team_name, total_amount
	`, userID, teamName).Scan(&app.ID, &app.UserID, &app.Status, &app.CreatedAt,
		&formedAt, &completedAt, &moderatorID, &app.TeamName, &app.TotalAmount)

	if err != nil {
		return nil, err
	}

	app.Services = []models.AppService{}
	log.Printf("📝 Создана новая заявка-черновик #%d для пользователя %d", app.ID, userID)
	return &app, nil
}

// GetApplicationServices - получить услуги заявки (M:M связь)
func (r *DatabaseRepository) GetApplicationServices(applicationID int) ([]models.AppService, error) {
	rows, err := r.db.Query(`
		SELECT asp.application_id, asp.service_id, asp.quantity, asp.sort_order, 
		       asp.is_primary, asp.added_at,
		       s.id, s.name, s.description, s.status, s.image_url, s.service_type, 
		       s.event_date, s.location, s.price, s.created_at
		FROM application_services asp
		JOIN services s ON asp.service_id = s.id
		WHERE asp.application_id = $1
		ORDER BY asp.sort_order
	`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.AppService
	for rows.Next() {
		var as models.AppService
		var s models.Service
		var eventDate, imageURL sql.NullString

		err := rows.Scan(&as.ApplicationID, &as.ServiceID, &as.Quantity, &as.SortOrder,
			&as.IsPrimary, &as.AddedAt,
			&s.ID, &s.Name, &s.Description, &s.Status, &imageURL,
			&s.ServiceType, &eventDate, &s.Location, &s.Price, &s.CreatedAt)
		if err != nil {
			return nil, err
		}

		if eventDate.Valid {
			s.EventDate = &eventDate.String
		}
		if imageURL.Valid {
			s.ImageURL = &imageURL.String
		}

		as.Service = s
		services = append(services, as)
	}

	return services, rows.Err()
}

// AddServiceToApplication - добавить услугу в заявку (через ORM-подобный подход)
func (r *DatabaseRepository) AddServiceToApplication(applicationID, serviceID, quantity, sortOrder int, isPrimary bool) error {
	// Проверяем, существует ли уже такая связь
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM application_services WHERE application_id = $1 AND service_id = $2)
	`, applicationID, serviceID).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		// Обновляем количество
		_, err = r.db.Exec(`
			UPDATE application_services 
			SET quantity = quantity + $1 
			WHERE application_id = $2 AND service_id = $3
		`, quantity, applicationID, serviceID)
		return err
	}

	// Вставляем новую запись
	_, err = r.db.Exec(`
		INSERT INTO application_services (application_id, service_id, quantity, sort_order, is_primary, added_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`, applicationID, serviceID, quantity, sortOrder, isPrimary)

	return err
}

// DeleteApplication - логическое удаление заявки через прямой SQL UPDATE (без ORM)
func (r *DatabaseRepository) DeleteApplication(applicationID int) error {
	log.Printf("🗑️ Логическое удаление заявки #%d через SQL UPDATE", applicationID)

	result, err := r.db.Exec(`
		UPDATE applications 
		SET status = 'deleted' 
		WHERE id = $1 AND status = 'draft'
	`, applicationID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("заявка #%d не найдена или уже удалена", applicationID)
	}

	log.Printf("✅ Заявка #%d успешно помечена как удаленная", applicationID)
	return nil
}

// GetApplicationByID - получить заявку по ID с услугами
func (r *DatabaseRepository) GetApplicationByID(id int) (*models.Application, error) {
	row := r.db.QueryRow(`
		SELECT id, user_id, status, created_at, formed_at, completed_at, 
		       moderator_id, team_name, total_amount 
		FROM applications 
		WHERE id = $1 AND status != 'deleted'
	`, id)

	var app models.Application
	var formedAt, completedAt sql.NullString
	var moderatorID sql.NullInt64

	err := row.Scan(&app.ID, &app.UserID, &app.Status, &app.CreatedAt,
		&formedAt, &completedAt, &moderatorID, &app.TeamName, &app.TotalAmount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	services, err := r.GetApplicationServices(app.ID)
	if err != nil {
		return nil, err
	}
	app.Services = services

	return &app, nil
}

// CalculateTotalAmount - рассчитать общую сумму заявки (дополнительное поле)
func (r *DatabaseRepository) CalculateTotalAmount(applicationID int) (float64, error) {
	var total float64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(s.price * asp.quantity), 0)
		FROM application_services asp
		JOIN services s ON asp.service_id = s.id
		WHERE asp.application_id = $1
	`, applicationID).Scan(&total)

	return total, err
}

// CompleteApplication - завершение заявки (расчет total_amount и обновление статуса)
func (r *DatabaseRepository) CompleteApplication(applicationID int, moderatorID int) error {
	total, err := r.CalculateTotalAmount(applicationID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		UPDATE applications 
		SET status = 'completed', 
		    completed_at = NOW(), 
		    moderator_id = $1,
		    total_amount = $2
		WHERE id = $3 AND status = 'formed'
	`, moderatorID, total, applicationID)

	return err
}
