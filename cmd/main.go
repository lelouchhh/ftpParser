package main

import (
	"flag"
	"fmt"
	"ftpParser/pkg"
	"github.com/secsy/goftp"
	"log"
)

var list = []string{
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

func main() {
	log.Println("Starting...")
	day := flag.String("day", "2022061700", "Дата для выгрузки\nФормат: ГодМесяцДень00\nнапример: 2022061700")
	outputZip := flag.String("outputZip", "./", "Папка для выгрузки zip")
	outputXml := flag.String("outputXML", "./", "Папка для выгрузки XMLs")
	flag.Parse()
	nextDay := pkg.GetNextDay(*day)

	log.Println(*day, nextDay)

	config := goftp.Config{
		User:            "free",
		Password:        "free",
		ActiveTransfers: true,
	}

	ftpConn, err := goftp.DialConfig(config, "ftp.zakupki.gov.ru:21")
	if err != nil {
		log.Fatalln(err)
	}
	for i, item := range list {
		path := fmt.Sprintf("/fcs_regions/%s/notifications/currMonth/", item)
		folder := pkg.GetFolders(ftpConn, fmt.Sprintf("/fcs_regions/%s/notifications/currMonth/", item))

		if len(folder) == 0 {
			log.Fatalln("folder is empty")
		} else {
			matchedList := pkg.MatchingFolder(folder, fmt.Sprintf("notification_%s_%s_%s_\n", item, *day, nextDay))
			pkg.DownloadZips(ftpConn, matchedList, path, *outputZip)
			go func(i int) {
				pkg.ExtractZip(matchedList, *outputZip, *outputXml)
			}(i)
		}
	}
	log.Println("Done!")
}
