package models

import "fmt"

type Partida struct {
	Ano_campeonato                 *uint
	Data                           *string
	Rodada                         *int
	Estadio                        *string
	Arbitro                        *string
	Publico                        *uint
	Publico_max                    *uint
	Time_mandante                  *string
	Time_visitante                 *string
	Tecnico_mandante               *string
	Tecnico_visitante              *string
	Colocacao_mandante             *uint
	Colocacao_visitante            *uint
	Valor_equipe_titular_mandante  *uint
	Valor_equipe_titular_visitante *uint
	Idade_media_titular_mandante   *float64
	Idade_media_titular_visitante  *float64
	Gols_mandante                  *uint
	Gols_visitante                 *uint
	Gols_1_tempo_mandante          *uint
	Gols_1_tempo_visitante         *uint
	Escanteios_mandante            *uint
	Escanteios_visitante           *uint
	Faltas_mandante                *uint
	Faltas_visitante               *uint
	Chutes_bola_parada_mandante    *uint
	Chutes_bola_parada_visitante   *uint
	Defesas_mandante               *uint
	Defesas_visitante              *uint
	Impedimentos_mandante          *uint
	Impedimentos_visitante         *uint
	Chutes_mandante                *uint
	Chutes_visitante               *uint
	Chutes_fora_mandante           *uint
	Chutes_fora_visitante          *uint
}

func (p Partida) String() string {
	return fmt.Sprintf("%s %d x %d %s - %do Rodada %d", *p.Time_mandante, *p.Gols_mandante, *p.Gols_visitante, *p.Time_visitante, *p.Rodada, *p.Ano_campeonato)
}

func CalcPontosByTime(partidas []Partida, time string) int {
	var pts int
	for _, p := range partidas {
		if *p.Time_mandante == time {
			if *p.Gols_mandante > *p.Gols_visitante {
				pts += 3
			} else if *p.Gols_mandante == *p.Gols_visitante {
				pts += 1
			}
		}

		if *p.Time_visitante == time {
			if *p.Gols_mandante < *p.Gols_visitante {
				pts += 3
			} else if *p.Gols_mandante == *p.Gols_visitante {
				pts += 1
			}
		}
	}

	return pts
}