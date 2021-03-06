package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func iAmRoot(w http.ResponseWriter, req *http.Request) {
	json := map[string]interface{}{
		"version": "0.1",
	}

	writeJSON(json, w, http.StatusOK)
}

func listProjects(w http.ResponseWriter, req *http.Request) {
	conn := pool.Get()
	defer conn.Close()

	ids, projects, moods, err := readProjects(conn)
	if err != nil {
		writeError(err, w)
		return
	}

	data := []Project{}
	for i := 0; i < len(projects); i++ {
		data = append(data, toProject(ids[i], projects[i], moods[i]))
	}

	writeJSON(data, w, http.StatusOK)
}

func removeProject(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	conn := pool.Get()
	defer conn.Close()

	err := deleteProject(id, conn)
	if err != nil {
		writeError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func newProject(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writeError(err, w)
		return
	}

	var project Project
	err = json.Unmarshal(body, &project)
	if err != nil {
		writeError(err, w)
		return
	}

	id, err := idOrNewUUID(project.ID)
	if err != nil {
		writeError(err, w)
		return
	}

	conn := pool.Get()
	defer conn.Close()

	projectDB, moodsByWeekDB := toProjectDB(project)
	err = writeProject(id, projectDB, moodsByWeekDB, conn)
	if err != nil {
		writeError(err, w)
		return
	}

	// read and output the written project and related moods.
	projectDB, moodsByWeekDB, err = readProject(id, conn)
	if err != nil {
		writeError(err, w)
		return
	}

	project = toProject(id, projectDB, moodsByWeekDB)
	writeJSON(project, w, http.StatusOK)
}

func setMood(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, week := vars["id"], vars["weekNo"]

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writeError(err, w)
		return
	}

	var mood Mood
	err = json.Unmarshal(body, &mood)
	if err != nil {
		writeError(err, w)
		return
	}

	conn := pool.Get()
	defer conn.Close()

	moodDB := toMoodDB(mood)
	err = writeMood(id, week, moodDB, conn)
	if err != nil {
		writeError(err, w)
	}

	w.WriteHeader(http.StatusOK)
}

func writeError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func writeJSON(data interface{}, w http.ResponseWriter, statusCode int) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to json.Marshal", data, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(payload)
}
