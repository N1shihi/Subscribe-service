package repository

import (
	"Subscribe-service/internal/models"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type SubsRepository struct {
	db *sql.DB
}

func NewSubsRepository(db *sql.DB) *SubsRepository {
	return &SubsRepository{db: db}
}

func nullableString(v *string) any {
	if v == nil {
		return nil
	}
	return *v
}

func scanSubscription(row interface {
	Scan(dest ...any) error
}) (*models.Subscription, error) {
	var sub models.Subscription
	var endDate sql.NullString

	err := row.Scan(
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&endDate,
	)
	if err != nil {
		return nil, err
	}

	if endDate.Valid {
		sub.EndDate = &endDate.String
	}

	return &sub, nil
}

func (r *SubsRepository) CreateSubs(subs *models.Subscription) error {
	if _, err := r.db.Exec(`
		INSERT INTO subscriptions (
			service_name,
			price,
			user_id,
			start_date,
			end_date
		)
		VALUES ($1, $2, $3, $4, $5)
	`, subs.ServiceName, subs.Price, subs.UserID, subs.StartDate, nullableString(subs.EndDate)); err != nil {

		log.Printf("repository create subscription failed user %s service %s error %v",
			subs.UserID.String(), subs.ServiceName, err)

		return fmt.Errorf("database error")
	}

	return nil
}

func (r *SubsRepository) GetSubs(userID uuid.UUID, serviceName string) (*models.Subscription, error) {
	row := r.db.QueryRow(`
		SELECT service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1 AND service_name = $2
	`, userID, serviceName)

	sub, err := scanSubscription(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subscription not found")
		}

		log.Printf("repository get subscription error user %s service %s error %v",
			userID.String(), serviceName, err)

		return nil, fmt.Errorf("database error")
	}

	return sub, nil
}

func (r *SubsRepository) GetAllSubs() ([]models.Subscription, error) {
	rows, err := r.db.Query(`
		SELECT service_name, price, user_id, start_date, end_date
		FROM subscriptions
		ORDER BY user_id, service_name
	`)
	if err != nil {
		log.Printf("repository get all subscriptions error %v", err)
		return nil, fmt.Errorf("database error")
	}
	defer rows.Close()

	subs := make([]models.Subscription, 0)

	for rows.Next() {
		var sub models.Subscription
		var endDate sql.NullString

		if err := rows.Scan(
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&endDate,
		); err != nil {

			log.Printf("repository scan subscription error %v", err)
			return nil, fmt.Errorf("database error")
		}

		if endDate.Valid {
			sub.EndDate = &endDate.String
		}

		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		log.Printf("repository rows iteration error %v", err)
		return nil, fmt.Errorf("database error")
	}

	return subs, nil
}

func (r *SubsRepository) UpdateSubs(subs *models.Subscription) error {
	res, err := r.db.Exec(`
		UPDATE subscriptions
		SET price = $1,
			start_date = $2,
			end_date = $3
		WHERE user_id = $4 AND service_name = $5
	`, subs.Price, subs.StartDate, nullableString(subs.EndDate), subs.UserID, subs.ServiceName)

	if err != nil {
		log.Printf("repository update subscription error user %s service %s error %v",
			subs.UserID.String(), subs.ServiceName, err)

		return fmt.Errorf("database error")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("repository update rows affected error %v", err)
		return fmt.Errorf("database error")
	}

	if affected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *SubsRepository) DeleteSubs(userID uuid.UUID, serviceName string) error {
	res, err := r.db.Exec(`
		DELETE FROM subscriptions
		WHERE user_id = $1 AND service_name = $2
	`, userID, serviceName)

	if err != nil {
		log.Printf("repository delete subscription error user %s service %s error %v",
			userID.String(), serviceName, err)

		return fmt.Errorf("database error")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("repository delete rows affected error %v", err)
		return fmt.Errorf("database error")
	}

	if affected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func parseMonthYear(v string) (time.Time, error) {
	return time.Parse("01-2006", v)
}

func (r *SubsRepository) GetSubsSum(req models.AggregateRequest) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE 1=1
	`

	args := make([]any, 0)

	if req.UserID != "" {
		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id")
		}

		args = append(args, userID)
		query += fmt.Sprintf(" AND user_id = $%d", len(args))
	}

	if req.ServiceName != "" {
		args = append(args, req.ServiceName)
		query += fmt.Sprintf(" AND service_name = $%d", len(args))
	}

	if req.StartDate != "" {
		fromDate, err := parseMonthYear(req.StartDate)
		if err != nil {
			return 0, fmt.Errorf("invalid start date")
		}

		args = append(args, fromDate)
		query += fmt.Sprintf(
			" AND (to_date(end_date, 'MM-YYYY') IS NULL OR to_date(end_date, 'MM-YYYY') >= $%d)",
			len(args),
		)
	}

	if req.EndDate != "" {
		toDate, err := parseMonthYear(req.EndDate)
		if err != nil {
			return 0, fmt.Errorf("invalid end date")
		}

		args = append(args, toDate)
		query += fmt.Sprintf(" AND to_date(start_date, 'MM-YYYY') <= $%d", len(args))
	}

	var total int

	if err := r.db.QueryRow(query, args...).Scan(&total); err != nil {
		log.Printf("repository aggregate query error %v", err)
		return 0, fmt.Errorf("database error")
	}

	return total, nil
}
