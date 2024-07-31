package main

import "fmt"

type Partida struct {
	ano_campeonato                 *uint
	data                           *string
	rodada                         *int
	estadio                        *string
	arbitro                        *string
	publico                        *uint
	publico_max                    *uint
	time_mandante                  *string
	time_visitante                 *string
	tecnico_mandante               *string
	tecnico_visitante              *string
	colocacao_mandante             *uint
	colocacao_visitante            *uint
	valor_equipe_titular_mandante  *uint
	valor_equipe_titular_visitante *uint
	idade_media_titular_mandante   *float64
	idade_media_titular_visitante  *float64
	gols_mandante                  *uint
	gols_visitante                 *uint
	gols_1_tempo_mandante          *uint
	gols_1_tempo_visitante         *uint
	escanteios_mandante            *uint
	escanteios_visitante           *uint
	faltas_mandante                *uint
	faltas_visitante               *uint
	chutes_bola_parada_mandante    *uint
	chutes_bola_parada_visitante   *uint
	defesas_mandante               *uint
	defesas_visitante              *uint
	impedimentos_mandante          *uint
	impedimentos_visitante         *uint
	chutes_mandante                *uint
	chutes_visitante               *uint
	chutes_fora_mandante           *uint
	chutes_fora_visitante          *uint
}

func (p Partida) String() string {
	return fmt.Sprintf("%s %d x %d %s - %do Rodada %d", *p.time_mandante, *p.gols_mandante, *p.gols_visitante, *p.time_visitante, *p.rodada, *p.ano_campeonato)
}

func CalcPontosByTime(partidas []Partida, time string) int {
	var pts int
	for _, p := range partidas {
		if *p.time_mandante == time {
			if *p.gols_mandante > *p.gols_visitante {
				pts += 3
			} else if *p.gols_mandante == *p.gols_visitante {
				pts += 1
			}
		}

		if *p.time_visitante == time {
			if *p.gols_mandante < *p.gols_visitante {
				pts += 3
			} else if *p.gols_mandante == *p.gols_visitante {
				pts += 1
			}
		}
	}

	return pts
}