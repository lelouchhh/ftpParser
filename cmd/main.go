package main

import (
	"flag"
	"fmt"
	"ftpParser/pkg"
	"github.com/secsy/goftp"
	"os"
	"sync"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Используйте 'cmd.exe -h' для получения справки")
		os.Exit(1)
	}
	//name := flag.String("region", "Dagestan_Resp", "Строка с названием региона")
	day := flag.String("day", "2022061700", "Дата для выгрузки\nФормат: ГодМесяцДень00\nнапример: 2022061700")
	outputZip := flag.String("outputZip", "./", "Папка для выгрузки zip")
	outputXml := flag.String("outputXML", "./", "Папка для выгрузки XMLs")
	flag.Parse()
	config := goftp.Config{
		User:            "free",
		Password:        "free",
		ActiveTransfers: true,
	}
	list := []string{
		"Adygeja_Resp",
		"Altaj_Resp",
		"Altajskij_kraj",
		"Amurskaja_obl",
		"Arkhangelskaja_obl",
		"Astrakhanskaja_obl",
		"Bajkonur_g",
		"Bashkortostan_Resp",
		"Belgorodskaja_obl",
		"Dagestan_Resp",
		"Brjanskaja_obl",
		"Burjatija_Resp",
		"Chechenskaja_Resp",
		"Cheljabinskaja_obl",
		"Chukotskij_AO",
		"Chuvashskaja_Resp",
		"Evrejskaja_Aobl",
		"Ingushetija_Resp",
		"Irkutskaja_obl",
		"Ivanovskaja_obl",
		"Jamalo-Neneckij_AO",
		"Jaroslavskaja_obl",
		"Kabardino-Balkarskaja_Resp",
		"Kaliningradskaja_obl",
		"Kalmykija_Resp",
		"Kaluzhskaja_obl",
		"Kamchatskij_kraj",
		"Karachaevo-Cherkesskaja_Resp",
		"Karelija_Resp",
		"Kemerovskaja_obl",
		"Khabarovskij_kraj",
		"Khakasija_Resp",
		"Khanty-Mansijskij_AO-Jugra_AO",
		"Kirovskaja_obl",
		"Komi_Resp",
		"Kostromskaja_obl",
		"Krasnodarskij_kraj",
		"Krasnojarskij_kraj",
		"Krim_Resp",
		"Kurganskaja_obl",
		"Kurskaja_obl",
		"Leningradskaja_obl",
		"Lipeckaja_obl",
		"Magadanskaja_obl",
		"Marij_El_Resp",
		"Mordovija_Resp",
		"Moskovskaja_obl",
		"Moskva",
		"Murmanskaja_obl",
		"Neneckij_AO",
		"Nizhegorodskaja_obl",
		"Novgorodskaja_obl",
		"Novosibirskaja_obl",
		"Omskaja_obl",
		"Orenburgskaja_obl",
		"Orlovskaja_obl",
		"Penzenskaja_obl",
		"Permskij_kraj",
		"Primorskij_kraj",
		"Pskovskaja_obl",
		"Rjazanskaja_obl",
		"Rostovskaja_obl",
		"Sakha_Jakutija_Resp",
		"Sakhalinskaja_obl",
		"Samarskaja_obl",
		"Sankt-Peterburg",
		"Saratovskaja_obl",
		"Sevastopol_g",
		"Severnaja_Osetija-Alanija_Resp",
		"Smolenskaja_obl",
		"Stavropolskij_kraj",
		"Sverdlovskaja_obl",
		"Tambovskaja_obl",
		"Tatarstan_Resp",
		"Tjumenskaja_obl",
		"Tomskaja_obl",
		"Tulskaja_obl",
		"Tverskaja_obl",
		"Tyva_Resp",
		"Udmurtskaja_Resp",
		"Uljanovskaja_obl",
		"Vladimirskaja_obl",
		"Volgogradskaja_obl",
		"Vologodskaja_obl",
		"Voronezhskaja_obl",
		"Zabajkalskij_kraj",
	}
	var wg sync.WaitGroup
	wg.Add(len(list))
	fmt.Println(list)
	for i, item := range list {
		fmt.Printf("i: %v, val: %v\n", i, *day)
		path := fmt.Sprintf("/fcs_regions/%s/notifications/currMonth/", item)
		ftpConn, _ := goftp.DialConfig(config, "ftp.zakupki.gov.ru:21")
		c := pkg.GetFolders(ftpConn, fmt.Sprintf("/fcs_regions/%s/notifications/currMonth/", item))
		if len(c) == 0 {
			os.Exit(1)
		} else {
			nextDay := pkg.GetNextDay(*day)
			list := pkg.MatchingFolder(c, fmt.Sprintf("notification_%s_%s_%s_\n", item, *day, nextDay))
			pkg.Write(ftpConn, list, path, *outputZip)
			go func(i int) {
				defer wg.Done()
				pkg.Extract(list, *outputZip, *outputXml)

			}(i)

		}
	}
}
