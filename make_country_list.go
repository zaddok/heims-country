package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type UN struct {
	Id     string
	Alpha2 string
	Alpha3 string
	Name   string
}

type AU struct {
	Id   string
	Name string
	UNId string
}

type Country struct {
	UN UN
	AU AU
}

func main() {
	var countries []*Country
	unCode := map[string]*Country{}

	f, err := os.Open("un.csv")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, line := range lines {
		if line[0] == "name" {
			continue
		}
		un := UN{
			Name:   strings.TrimSpace(line[0]),
			Alpha2: strings.TrimSpace(line[1]),
			Alpha3: strings.TrimSpace(line[2]),
			Id:     strings.TrimSpace(line[3]),
		}
		country := &Country{UN: un}
		countries = append(countries, country)
		unCode[un.Id] = country
	}

	f, err = os.Open("australia-un.csv")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	lines, err = csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, line := range lines {
		au := AU{
			Id:   strings.TrimSpace(line[0]),
			Name: strings.TrimSpace(line[1]),
			UNId: strings.TrimSpace(line[2]),
		}
		var country *Country
		if au.UNId != "" {
			country = unCode[au.UNId]
			if country == nil {
				unCode[au.UNId] = country
				fmt.Println("WARNING: Australian data set refers to a UN country ID that the UN data set does not list", au.Name)
			}
			country.AU = au
			if au.Name != country.UN.Name {
				fmt.Printf("INFO: Australian UN Naming difference. AU:'%s' UN:'%s'\n", au.Name, country.UN.Name)
			}
		} else {
			country = &Country{AU: au}
			countries = append(countries, country)
		}
	}

	x := &Country{UN: UN{"804", "Ukraine", "UA", "UKR"}, AU: AU{"3312", "Ukraine", ""}}
	fmt.Println(x)
	for _, country := range countries {
		if country.UN.Id != "" && country.AU.Id == "" {
			fmt.Println("WARNING: UN data set refers to a country ID that the Australian data set does not list:", country.UN.Name)
		}
	}

	for _, country := range countries {
		fmt.Printf("\tCountry{UN: UN{\"%s\",\"%s\",\"%s\",\"%s\"}, AU: AU{\"%s\",\"%s\", \"%s\"}},\n", country.UN.Id, country.UN.Name, country.UN.Alpha2, country.UN.Alpha3, country.AU.Id, country.AU.Name, country.AU.UNId)
	}

}
