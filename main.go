package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	//SQL Login
	server := "localhost"
	port := "1433"
	user := "invadmin"
	password := "invadmin"
	database := "CHEMINVDB2"

	input_value := ""
	fmt.Print("INVENTORY SQL ADDRESS:" + server + ">")
	fmt.Scanln(&input_value)
	if input_value != "" {
		server = input_value
	}

	input_value = ""
	fmt.Print("INVENTORY SQL Port:" + port + ">")
	fmt.Scanln(&input_value)
	if input_value != "" {
		port = input_value
	}

	input_value = ""
	fmt.Print("INVENTORY Username:" + user + ">")
	fmt.Scanln(&input_value)
	if input_value != "" {
		user = input_value
	}

	input_value = ""
	fmt.Print("INVENTORY password:" + password + ">")
	fmt.Scanln(&input_value)
	if input_value != "" {
		password = input_value
	}

	input_value = ""
	fmt.Print("INVENTORY database name:" + database + ">")
	fmt.Scanln(&input_value)
	if input_value != "" {
		database = input_value
	}

	int_port, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("illegal port:", err.Error())
		os.Exit(0)
	}

	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s", server, int_port, database, user, password)
	fmt.Println(connString)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
		os.Exit(0)
	}

	rows, err := conn.Query("select * from INV_COMPOUNDS WHERE BASE64_CDX not like 'Vmp%';")
	if err != nil {
		log.Fatal(err)
	}

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var result []map[string]string

	for rows.Next() {
		record := make(map[string]string)
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			record[string(columns[i])] = string(value)
		}
		result = append(result, record)
	}
	defer rows.Close()
	defer conn.Close()

	for line := range result {
		this_line := result[line]
		fmt.Println("----------")
		fmt.Println("COMPOUND ID:" + this_line["COMPOUND_ID"])
		fmt.Println("CAS:" + this_line["CAS"])
		fmt.Println("SUBSTANCE NAME:" + this_line["SUBSTANCE_NAME"])

		SMILES := ""
		urladdr := "https://www.ncbi.nlm.nih.gov/pccompound?term=" + this_line["CAS"]
		resp, err := http.Get(urladdr)
		if err != nil {
			log.Fatal(err)
		}
		bytebody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		content := string(bytebody)

		if strings.Contains(content, " <meta property=\"og:url\" content=\"https://pubchem.ncbi.nlm.nih.gov/compound/") {
			CID := strings.Split(content, " <meta property=\"og:url\" content=\"https://pubchem.ncbi.nlm.nih.gov/compound/")[1]
			CID = strings.Split(CID, "\"")[0]

			urladdr = "https://pubchem.ncbi.nlm.nih.gov/rest/pug_view/data/compound/" + CID + "/JSON"
			resp, err := http.Get(urladdr)
			if err != nil {
				log.Fatal(err)
			}
			bytebody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			content := string(bytebody)

			content = strings.Split(content, "\"Name\": \"Canonical SMILES\",")[1]
			content = strings.Split(content, "}")[0]

			SMILES = strings.Split(content, "\"StringValue\": \"")[1]
			SMILES = strings.Split(SMILES, "\"")[0]

		}

		input_value = ""
		fmt.Print("SMILES:" + SMILES + ">")
		fmt.Scanln(&input_value)
		if input_value != "" {
			SMILES = input_value
		}
		ioutil.WriteFile("smiles.input", []byte(SMILES), 0644)
		cmd := exec.Command("convert.exe")
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		bytebody, err = ioutil.ReadFile("cdx.output")
		if err != nil {
			log.Fatal(err)
		}
		content = string(bytebody)
		CDX := content

		SQL := "UPDATE INV_COMPOUNDS SET BASE64_CDX='" + CDX + "' WHERE  COMPOUND_ID=" + this_line["COMPOUND_ID"] + ";"
		conn.Exec(SQL)

	}
	defer conn.Close()
}
