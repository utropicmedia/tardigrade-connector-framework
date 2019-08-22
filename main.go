// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"unsafe"

	"mongo"
	"storj"

	"github.com/urfave/cli"
)

const dbConfigFile = "./config/db_property.json"
const storjConfigFile = "./config/storj_config.json"

var gbDEBUG = false

// Create command-line tool to read from CLI.
var app = cli.NewApp()

// SetAppInfo sets information about the command-line application.
func setAppInfo() {
	app.Name = "Storj Connector Framework"
	app.Usage = "Framework to connect Storj network to various endpoints and data sources "
	app.Author = "Satyam Shivam - Utropicmedia\nKarl Mozurkewich - Utropicmedia"
	app.Version = "0.5"

}

// helper function to flag debug
func setDebug(debugVal bool) {
	gbDEBUG = debugVal
	mongo.DEBUG = debugVal
	storj.DEBUG = debugVal
}

// setCommands sets various command-line options for the app.
func setCommands() {

	app.Commands = []cli.Command{
		{
			Name:    "parse",
			Aliases: []string{"p"},
			Usage:   "Command to read and parse JSON information about MongoDB instance properties and then fetch ALL its collections. ",
			//\narguments-\n\t  fileName [optional] = provide full file name (with complete path), storing mongoDB properties if this fileName is not given, then data is read from ./config/db_connector.json\n\t  example = ./storj_mongodb d ./config/db_property.json\n",
			Action: func(cliContext *cli.Context) {
				var fullFileName = dbConfigFile

				// process arguments
				if len(cliContext.Args()) > 0 {
					for i := 0; i < len(cliContext.Args()); i++ {

						// Incase, debug is provided as argument.
						if cliContext.Args()[i] == "debug" {
							setDebug(true)
						} else {
							fullFileName = cliContext.Args()[i]
						}
					}
				}

				// Connect to Database and process data
				data, dbname, err := mongo.ConnectToDBFetchData(fullFileName)

				if err != nil {
					log.Fatalf("mongo.ConnectToDBFetchData: %s", err)
				} else {
					fmt.Println("Reading ALL collections from the MongoDB database...Complete!")
				}

				if gbDEBUG {
					fmt.Println("Size of fetched data from database :", dbname, unsafe.Sizeof(data))
				}
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Command to read and parse JSON information about Storj network and upload sample JSON data",
			//\n arguments- 1. fileName [optional] = provide full file name (with complete path), storing Storj configuration information if this fileName is not given, then data is read from ./config/storj_config.json example = ./storj_mongodb s ./config/storj_config.json\n\n\n",
			Action: func(cliContext *cli.Context) {

				// Default Storj configuration file name.
				var fullFileName = storjConfigFile

				// process arguments
				if len(cliContext.Args()) > 0 {
					for i := 0; i < len(cliContext.Args()); i++ {

						// Incase, debug is provided as argument.
						if cliContext.Args()[i] == "debug" {
							setDebug(true)
						} else {
							fullFileName = cliContext.Args()[i]
						}
					}
				}

				// Sample database name and data to be uploaded
				dbName := "testdb"
				jsonData := "{'testKey': 'testValue'}"

				// Converting JSON data to bson data.  TODO: convert to BSON using call to mongo library
				bsonData, _ := json.Marshal(jsonData)

				if gbDEBUG {
					t := time.Now()
					time := t.Format("2006-01-02_15:04:05")
					var fileName = "uploaddata_" + time + ".bson"

					err := ioutil.WriteFile(fileName, bsonData, 0644)
					if err != nil {
						fmt.Println("Error while writting to file ")
					}
				}

				err := storj.ConnectStorjUploadData(fullFileName, []byte(bsonData), dbName)
				if err != nil {
					fmt.Println("Error while uploading data to the Storj bucket")
				}
			},
		},
		{
			Name:    "store",
			Aliases: []string{"s"},
			Usage:   "Command to connect and transfer ALL collections from a desired MongoDB instance to given Storj Bucket in BSON format",
			//\n    arguments-\n      1. fileName [optional] = provide full file name (with complete path), storing mongoDB properties in JSON format\n   if this fileName is not given, then data is read from ./config/db_property.json\n      2. fileName [optional] = provide full file name (with complete path), storing Storj configuration in JSON format\n     if this fileName is not given, then data is read from ./config/storj_config.json\n   example = ./storj_mongodb c ./config/db_property.json ./config/storj_config.json\n",
			Action: func(cliContext *cli.Context) {

				// Default configuration file names.
				var fullFileNameStorj = storjConfigFile
				var fullFileNameMongoDB = dbConfigFile

				// process arguments - Reading fileName from the command line.
				var foundFirstFileName = false
				if len(cliContext.Args()) > 0 {
					for i := 0; i < len(cliContext.Args()); i++ {
						// Incase debug is provided as argument.
						if cliContext.Args()[i] == "debug" {
							setDebug(true)
						} else {
							if !foundFirstFileName {
								fullFileNameMongoDB = cliContext.Args()[i]
								foundFirstFileName = true
							} else {
								fullFileNameStorj = cliContext.Args()[i]
							}
						}
					}
				}

				// Fetching data from mongodb.
				data, dbname, err := mongo.ConnectToDBFetchData(fullFileNameMongoDB)

				if err != nil {
					log.Fatalf("mongo.ConnectToDBFetchData: %s", err)
				}

				if gbDEBUG {
					fmt.Println("Size of fetched data from database: ", dbname, unsafe.Sizeof(data))
				}
				// Connecting to storj network for uploading data.
				err = storj.ConnectStorjUploadData(fullFileNameStorj, []byte(data), dbname)
				if err != nil {
					fmt.Println("Error while uploading data to bucket ", err)
				}
			},
		},
	}
}

func main() {

	setAppInfo()
	setCommands()

	setDebug(false)

	err := app.Run(os.Args)

	if err != nil {
		log.Fatalf("app.Run: %s", err)
	}
}
