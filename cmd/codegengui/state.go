package main

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/jackmanlabs/errors"
)

type PackageState struct {
	ImportPath    string
	SelectedTypes []string
	BindingsPath  string
	SchemaPath    string
	InterfacePath string
	SqlDriver     string
	WriteTests    bool
}

func loadState() (map[string]PackageState, error) {

	state := make(map[string]PackageState)

	statePath := os.Getenv("HOME")
	if statePath == "" {
		log.Print("WARNING: $HOME is not set. Application state will not be saved.")
		return state, nil
	}

	statePath = path.Join(statePath, ".codegen", "state.json")

	f, err := os.Open(statePath)
	if os.IsNotExist(err) {
		// Do nothing. This file will be created when it's saved.
		return state, nil
	} else if err != nil {
		return nil, errors.Stack(err)
	}

	err = json.NewDecoder(f).Decode(&state)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return state, nil
}

func saveState(state map[string]PackageState) error {

	statePath := os.Getenv("HOME")
	if statePath == "" {
		log.Print("WARNING: $HOME is not set. Application state will not be saved.")
		return nil
	}

	statePath = path.Join(statePath, ".codegen", "state.json")

	dir := path.Dir(statePath)
	f, err := os.Open(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm|os.ModeDir)
		if err != nil {
			return errors.Stack(err)
		}
	} else if err != nil {
		return errors.Stack(err)
	} else {
		err = f.Close()
		if err != nil {
			return errors.Stack(err)
		}
	}

	f, err = os.Create(statePath)
	if err != nil {
		return errors.Stack(err)
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	err = enc.Encode(state)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
