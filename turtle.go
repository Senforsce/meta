package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/nvkp/turtle"
)

type RDFTriple struct {
	Subject   string `turtle:"subject"`
	Predicate string `turtle:"predicate"`
	Object    string `turtle:"object"`
	Label     string `turtle:"label"`
	DataType  string `turtle:"datatype"`
}

type RDFTripleList []RDFTriple

func readTurtleString(rdf []byte) (RDFTripleList, error) {
	var triples = RDFTripleList{}

	err := turtle.Unmarshal(
		[]byte(rdf),
		&triples,
	)
	if err != nil {
		return nil, fmt.Errorf("(ttl Marshalling failed)")
	}

	return triples, nil

}

var userTurtlePath = filepath.Join(os.Getenv("O8ROOT"), "o8", "config.ttl")

func loadUserTurtle() (RDFTripleList, error) {
	f, readFileErr := loadUserTurtleFile()

	if readFileErr != nil {
		return nil, readFileErr
	}
	defer f.Close()

	byteValue, _ := io.ReadAll(f)

	triples, marshallErr := readTurtleString(byteValue)

	if marshallErr != nil {
		return nil, marshallErr
	}

	return triples, nil
}

func loadUserTurtleFile() (fs.File, error) {
	return os.Open(userTurtlePath)
}

func saveUserTurtle(config Config) error {
	config.Input = ""
	config.Output = ""
	config.Interactive = false

	err := os.MkdirAll(filepath.Dir(userConfigPath), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(userConfigPath)
	if err != nil {
		return err
	}
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = f.Write(b)

	print(userConfigPath, "SAVED")

	return err
}
