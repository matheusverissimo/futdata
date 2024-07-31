package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const databaseFile string = "./data/futdata.db"

var dbconn *sql.DB

func convertCSVRecordToSqlNullString(record []string) []sql.NullString {
	var returned []sql.NullString
	for _, s := range record {
		var ns sql.NullString
		if len(s) == 0 {
			ns = sql.NullString{}
		} else {
			ns = sql.NullString{
				String: s,
				Valid:  true,
			}
		}
		returned = append(returned, ns)
	}
	return returned
}

func DropDatabase() {
	err := os.Remove(databaseFile)
	if err != nil {
		fmt.Println(err)
	}
}

func connectToDatabase() error {
	var err error
	dbconn, err = sql.Open("sqlite3", databaseFile)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return err
	}
	err = dbconn.Ping()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func InitDatabase() error {
	connectToDatabase()

	if !checkSchemaState() {
		LoadDatabaseFromCSVFile("./data/brasileirao_serie_a.csv")
	}

	return nil
}

func checkSchemaState() bool {
	row := dbconn.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND tbl_name IN ('times', 'partidas', 'arbitros', 'estadios', 'tecnicos')")

	var count int
	row.Scan(&count)

	if count != 5 {
		return false
	}

	return true
}

func createTables() error {
	errs := []error{
		createArbitrosTable(dbconn),
		createEstadiosTable(dbconn),
		createTimesTable(dbconn),
		createTecnicosTable(dbconn),
		createPartidasTable(dbconn),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func createTimesTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS times(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT UNIQUE
	);`)

	if err != nil {
		return err
	}

	return nil
}

func createArbitrosTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS arbitros(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT UNIQUE
	);`)

	if err != nil {
		return err
	}

	return nil
}

func createEstadiosTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS estadios(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT UNIQUE
	);`)

	if err != nil {
		return err
	}

	return nil
}

func createTecnicosTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tecnicos(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT UNIQUE
	);`)

	if err != nil {
		return err
	}

	return nil
}

func createPartidasTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS partidas(
			ano_campeonato INTEGER,
			data TEXT, 
			rodada INTEGER,
			estadio_id INTEGER,
			arbitro_id INTEGER,
			publico INTEGER,
			publico_max INTEGER,
			time_mandante_id INTEGER,
			time_visitante_id INTEGER,
			tecnico_mandante_id INTEGER,
			tecnico_visitante_id INTEGER,
			colocacao_mandante INTEGER,
			colocacao_visitante INTEGER,
			valor_equipe_titular_mandante INTEGER,
			valor_equipe_titular_visitante INTEGER,
			idade_media_titular_mandante FLOAT,
			idade_media_titular_visitante FLOAT,
			gols_mandante INTEGER,
			gols_visitante INTEGER,
			gols_1_tempo_mandante INTEGER,
			gols_1_tempo_visitante INTEGER,
			escanteios_mandante INTEGER,
			escanteios_visitante INTEGER,
			faltas_mandante INTEGER,
			faltas_visitante INTEGER,
			chutes_bola_parada_mandante INTEGER,
			chutes_bola_parada_visitante INTEGER,
			defesas_mandante INTEGER,
			defesas_visitante INTEGER,
			impedimentos_mandante INTEGER,
			impedimentos_visitante INTEGER,
			chutes_mandante INTEGER,
			chutes_visitante INTEGER,
			chutes_fora_mandante INTEGER,
			chutes_fora_visitante INTEGER,
			PRIMARY KEY(ano_campeonato, rodada, time_mandante_id, time_visitante_id),
			FOREIGN KEY(estadio_id) REFERENCES estadios(id),
			FOREIGN KEY(arbitro_id) REFERENCES arbitros(id),
			FOREIGN KEY(time_mandante_id) REFERENCES times(id),
			FOREIGN KEY(time_visitante_id) REFERENCES times(id),
			FOREIGN KEY(tecnico_mandante_id) REFERENCES tecnicos(id),
			FOREIGN KEY(tecnico_visitante_id) REFERENCES tecnicos(id)
	);`)

	if err != nil {
		return err
	}

	return nil
}

func LoadDatabaseFromCSVFile(csvFile string) error {
	var err error
	if dbconn == nil {
		err = connectToDatabase()
	}

	if err != nil {
		return err
	}

	err = createTables()

	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	reader.Read()
	for {
		row, err := reader.Read()
		record := convertCSVRecordToSqlNullString(row)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}

		tx, err := dbconn.Begin()

		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		var errs [5]error

		_, errs[0] = tx.Exec("INSERT OR IGNORE INTO times(nome) VALUES (?), (?);", record[7], record[8])
		_, errs[1] = tx.Exec("INSERT OR IGNORE INTO estadios(nome) VALUES (?);", record[3])
		_, errs[2] = tx.Exec("INSERT OR IGNORE INTO arbitros(nome) VALUES (?);", record[4])
		_, errs[3] = tx.Exec("INSERT OR IGNORE INTO tecnicos(nome) VALUES (?), (?);", record[9], record[10])
		_, errs[4] = tx.Exec(`INSERT OR IGNORE INTO partidas VALUES(?,?,?,
			(SELECT id from estadios WHERE nome = ?),
			(SELECT id from arbitros WHERE nome = ?),?,?,
			(SELECT id from times WHERE nome = ?),
			(SELECT id from times WHERE nome = ?),
			(SELECT id from tecnicos WHERE nome = ?),
			(SELECT id from tecnicos WHERE nome = ?),
			?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, 
			record[0],
			record[1],
			record[2],
			record[3],
			record[4],
			record[5],
			record[6],
			record[7],
			record[8],
			record[9],
			record[10],
			record[11],
			record[12],
			record[13],
			record[14],
			record[15],
			record[16],
			record[17],
			record[18],
			record[19],
			record[20],
			record[21],
			record[22],
			record[23],
			record[24],
			record[25],
			record[26],
			record[27],
			record[28],
			record[29],
			record[30],
			record[31],
			record[32],
			record[33],
			record[34])

		for _, e := range errs {
			if e != nil {
				fmt.Println(e.Error())
			}
		}

		err = tx.Commit()
	}

	if err != nil {
		return err
	}

	return nil
}

func GetAllMatchesByYearAndTeam(year int, team string) ([]Partida, error) {
	partidas := []Partida{}
	rows, err := dbconn.Query(fmt.Sprintf(`SELECT p.ano_campeonato, p.rodada, 
											(SELECT nome FROM times WHERE id = p.time_mandante_id),
											(SELECT nome FROM times WHERE id = p.time_visitante_id),
											p.gols_mandante,
											p.gols_visitante
											FROM partidas p
											INNER JOIN times t ON t.id = p.time_mandante_id OR t.id = p.time_visitante_id 
											WHERE p.ano_campeonato = %v AND t.nome = '%v';`, year, team))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		partida := Partida{}
		err := rows.Scan(&partida.ano_campeonato,
			&partida.rodada,
			&partida.time_mandante,
			&partida.time_visitante,
			&partida.gols_mandante,
			&partida.gols_visitante)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		partidas = append(partidas, partida)
	}

	return partidas, nil
}
