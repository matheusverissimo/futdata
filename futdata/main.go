package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	err := InitDatabase()
	year, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		return 
	}
	partidas, err := GetAllMatchesByYearAndTeam(int(year), os.Args[2])

	pts := CalcPontosByTime(partidas, os.Args[2])

	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}

	for _, p := range partidas {
		fmt.Printf("%v\n", p)
	}
	fmt.Printf("O %v fez %d pontos no ano de %d!", os.Args[2], pts, year)
}
