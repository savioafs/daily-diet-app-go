package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"savioafs/daily-diet-app-go/internal/entity"
)

type MealRepositoryPG struct {
	DB *sql.DB
}

func NewMealRepositoryPG(db *sql.DB) *MealRepositoryPG {
	return &MealRepositoryPG{DB: db}
}

func (r *MealRepositoryPG) Create(meal *entity.Meal) (string, error) {

	var id string

	stmt, err := r.DB.Prepare("INSERT INTO meals (id, user_id, name, description, date, is_diet) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	err = stmt.QueryRow(meal.ID, meal.UserID, meal.Name, meal.Description, meal.Date, meal.IsDiet).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *MealRepositoryPG) GetMealByID(id string) (*entity.Meal, error) {
	stmt, err := r.DB.Prepare("SELECT id, user_id, name, description, date, is_diet FROM meals WHERE id = $1 ")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var meal entity.Meal

	err = stmt.QueryRow(id).Scan(
		&meal.ID,
		&meal.UserID,
		&meal.Name,
		&meal.Description,
		&meal.Date,
		&meal.IsDiet)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &meal, nil
}

func (r *MealRepositoryPG) GetAllMealsByUser(userID string) ([]entity.Meal, error) {
	query := "SELECT * FROM meals WHERE user_id = $1"

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}

	var meals []entity.Meal

	for rows.Next() {
		var meal entity.Meal

		err = rows.Scan(
			&meal.ID,
			&meal.UserID,
			&meal.Name,
			&meal.Description,
			&meal.Date,
			&meal.IsDiet)

		if err != nil {
			return nil, err
		}

		meals = append(meals, meal)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return meals, nil
}

func (r *MealRepositoryPG) UpdateMeal(id string, meal *entity.Meal) error {
	query := "UPDATE meals SET name = $1, description = $2, date = $3, is_diet = $4  WHERE id = $5"

	mealFind, err := r.GetMealByID(id)
	if err != nil {
		return err
	}

	if mealFind == nil {
		return errors.New("invalid meal id")
	}

	result, err := r.DB.Exec(query, meal.Name, meal.Description, meal.Date, meal.IsDiet, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no meal found with the given id")
	}

	return nil
}

func (r *MealRepositoryPG) GetMealsByDay(date string) ([]entity.Meal, error) {
	query := "SELECT * FROM meals WHERE DATE(date) = $1"

	rows, err := r.DB.Query(query, date)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var meals []entity.Meal
	for rows.Next() {
		var meal entity.Meal
		err := rows.Scan(
			&meal.ID,
			&meal.UserID,
			&meal.Name,
			&meal.Description,
			&meal.Date,
			&meal.IsDiet)
		if err != nil {
			return nil, err
		}

		meals = append(meals, meal)
	}

	return meals, nil
}

func (r *MealRepositoryPG) DeleteMeal(id string) error {
	query := "DELETE FROM meals WHERE id = $1"

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no meal found with the given id")
	}

	return nil
}
