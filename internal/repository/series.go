package repository

import (
	"context"
	"sycinema/internal/models"
)

func GetAllSeries() ([]models.Series, error) {
	query := `SELECT id, title, description, poster_url, release_year FROM series`
	rows, err := DB.Query(context.Background(), query)
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

func GetSeriesByID(id int) (models.Series, error) {
	query := `SELECT id, title, description, poster_url, release_year FROM series WHERE id = $1`
	var s models.Series
	err := DB.QueryRow(context.Background(), query, id).Scan(&s.ID, &s.Title, &s.Description, &s.PosterURL, &s.ReleaseYear)
	return s, err
}

func GetEpisodesBySeriesID(seriesID int, userID int) ([]models.Episode, error) {
	// Делаем LEFT JOIN с таблицей лайков. Если совпадение есть, is_favorite = true.
	query := `
		SELECT e.id, e.series_id, e.season, e.episode_num, e.title, e.video_path,
		       (fe.user_id IS NOT NULL) AS is_favorite
		FROM episodes e
		LEFT JOIN favorite_episodes fe ON e.id = fe.episode_id AND fe.user_id = $1
		WHERE e.series_id = $2
		ORDER BY e.season, e.episode_num`
		
	rows, err := DB.Query(context.Background(), query, userID, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []models.Episode
	for rows.Next() {
		var ep models.Episode
		// Обязательно добавляем &ep.IsFavorite в Scan!
		if err := rows.Scan(&ep.ID, &ep.SeriesID, &ep.Season, &ep.EpisodeNum, &ep.Title, &ep.VideoPath, &ep.IsFavorite); err != nil {
			return nil, err
		}
		episodes = append(episodes, ep)
	}
	return episodes, nil
}

func GetEpisodeByID(id int) (models.Episode, error) {
	query := `SELECT id, series_id, season, episode_num, title, video_path FROM episodes WHERE id = $1`
	var ep models.Episode
	err := DB.QueryRow(context.Background(), query, id).Scan(&ep.ID, &ep.SeriesID, &ep.Season, &ep.EpisodeNum, &ep.Title, &ep.VideoPath)
	return ep, err
}