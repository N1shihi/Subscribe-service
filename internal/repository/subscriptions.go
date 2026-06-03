package repository

import (
	"Subscribe-service/internal/models"
	"database/sql"
	"fmt"
	"log"
)

type SubsRepository struct {
	db *sql.DB
}

func NewSubsRepository(db *sql.DB) *SubsRepository {
	return &SubsRepository{
		db: db,
	}
}

func (r *SubsRepository) GetSubsById(id int) (*models.Team, error) {
	log.SetPrefix("repository.TeamRepository.GetTeamByName: ")

	var subsExists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM teams WHERE subs_id = $1)",
		name,
	).Scan(&teamExists)
	if err != nil {
		log.Printf("error checking team existence: %v", err)
		return nil, fmt.Errorf("database error")
	}

	if !teamExists {
		return nil, fmt.Errorf("team not found")
	}

	query := `
        SELECT user_id, username, is_active
        FROM users
        WHERE team_name = $1
    `
	rows, err := r.db.Query(query, name)
	if err != nil {
		log.Printf("error querying team users: %v", err)
		return nil, fmt.Errorf("database error")
	}
	defer func() {
		_ = rows.Close()
	}()

	team := &models.Team{
		TeamName: name,
		Members:  []models.TeamMember{},
	}

	for rows.Next() {
		var member models.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			log.Printf("error scanning user row: %v", err)
			return nil, fmt.Errorf("database error")
		}
		team.Members = append(team.Members, member)
	}

	if err = rows.Err(); err != nil {
		log.Printf("rows iteration error: %v", err)
		return nil, fmt.Errorf("database error")
	}

	log.Printf("successfully retrieved team '%s' with %d members", name, len(team.Members))
	return team, nil
}

func (r *SubsRepository) CreateSubs(subs *models.Subscription) error {
	log.SetPrefix("repository.SubsRepository.CreateSubs: ")

	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("begin tx error: %v", err)
		return fmt.Errorf("database error")
	}

	_, err = tx.Exec(`
        INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4,$5)
    `, subs.Service_name, subs.Price, subs.User_id, subs.Start_date, subs.End_date)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("database error")
	}

	if err = tx.Commit(); err != nil {
		log.Printf("commit tx error: %v", err)
		return fmt.Errorf("database error")
	}

	log.Printf("created subscription '%s' ", subs.Service_name)
	return nil
}
