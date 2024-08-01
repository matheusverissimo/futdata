package db

import (
	"database/sql"
	"fmt"
	"futdata/pkg/models"
)

type Repository struct {
	conn *sql.DB
}

func NewRepository() *Repository {
	return &Repository{
		conn: dbconn,
	}
}

func (r *Repository) FindAllTimes() ([]string, error) {
	rows, err := dbconn.Query("SELECT nome FROM times;")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var times []string

	for rows.Next() {
		var time string
		if err := rows.Scan(&time); err != nil {
			return nil, err
		}
		times = append(times, time)
	}

	return times, nil
}

func (r *Repository) FindAllPartidasByAnoTime(year int, team string) ([]models.Partida, error) {
	partidas := []models.Partida{}
	rows, err := dbconn.Query(fmt.Sprintf(`SELECT p.ano_campeonato, p.rodada, 
											(SELECT nome FROM times WHERE id = p.time_mandante_id),
											(SELECT nome FROM times WHERE id = p.time_visitante_id),
											p.gols_mandante,
											p.gols_visitante
											FROM partidas p
											INNER JOIN times t ON t.id = p.time_mandante_id OR t.id = p.time_visitante_id 
											WHERE p.ano_campeonato = %v AND t.nome = '%v';`, year, team))

	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		partida := models.Partida{}
		err := rows.Scan(&partida.Ano_campeonato,
			&partida.Rodada,
			&partida.Time_mandante,
			&partida.Time_visitante,
			&partida.Gols_mandante,
			&partida.Gols_visitante)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		partidas = append(partidas, partida)
	}

	return partidas, nil
}

func (r *Repository) FindAnosByTime(time string) ([]string, error) {
	rows, err := r.conn.Query("SELECT DISTINCT ano_campeonato FROM partidas WHERE time_mandante_id = (SELECT id FROM times where nome = ?);", time)
	if err != nil {
		return nil, err
	}

	var anos []string
	
	for rows.Next(){
		var ano string
		rows.Scan(&ano)
		anos = append(anos, ano)
	}

	return anos, nil
}