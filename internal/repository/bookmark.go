package repository

import (
	"context"
	"sycinema/internal/models"
)

func IsBookmarked(userID, seriesID int) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id=$1 AND series_id=$2)`
	DB.QueryRow(context.Background(), query, userID, seriesID).Scan(&exists)
	return exists
}

func ToggleBookmark(userID, seriesID int) error {
	if IsBookmarked(userID, seriesID) {
		_, err := DB.Exec(context.Background(), `DELETE FROM bookmarks WHERE user_id=$1 AND series_id=$2`, userID, seriesID)
		return err
	}
	_, err := DB.Exec(context.Background(), `INSERT INTO bookmarks (user_id, series_id) VALUES ($1, $2)`, userID, seriesID)
	return err
}

func GetBookmarks(userID int) ([]models.Series, error) {
	query := `
		SELECT s.id, s.title, s.description, s.poster_url, s.release_year 
		FROM series s
		JOIN bookmarks b ON s.id = b.series_id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC`
	
	rows, err := DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var series []models.Series
	for rows.Next() {
		var s models.Series
		if err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.PosterURL, &s.ReleaseYear); err != nil {
			return nil, err
		}
		series = append(series, s)
	}
	return series, nil
}

func ToggleFavoriteEpisode(userID, episodeID int) error {
	var exists bool
	DB.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM favorite_episodes WHERE user_id=$1 AND episode_id=$2)`, userID, episodeID).Scan(&exists)
	
	if exists {
		_, err := DB.Exec(context.Background(), `DELETE FROM favorite_episodes WHERE user_id=$1 AND episode_id=$2`, userID, episodeID)
		return err
	}
	_, err := DB.Exec(context.Background(), `INSERT INTO favorite_episodes (user_id, episode_id) VALUES ($1, $2)`, userID, episodeID)
	return err
}