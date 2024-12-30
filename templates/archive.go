// File: archive.go
// Authors: Charisman Aravinthan and Justin Dai
// Last Modified: June 2, 2023
// Version: 1.1.0
// This file contains the backend for the archive page

// Problem Archive Page
// Lists the list of problems to user
// Allows user to click on problems and filter them for specific problems
// Allows addition and deletion of problems for admin
package backend

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"strings"
	"io/ioutil"
)

// Function to list the problems. The function takes an HTTP response writer and an HTTP request as arguments.
// Returns None.
// Contributers: Justin Dai, Charisman Aravinthan.
func ListProblems (w http.ResponseWriter, r *http.Request) {

	// Retrieve data and set its state to "0"; no response needed from page.
	data := GetData(w, r)
	data.State = "0"

	// Retrieve filter parameters from the URL query string.
	filterSearch := r.URL.Query().Get("search")
	filterTag := r.URL.Query().Get("tag")
	filterDifficulty := r.URL.Query().Get("difficulty")
	filterStatus := r.URL.Query().Get("status")

	// Depending on the authorization of the user, choose a different set of templates.
	files := []string{}
	if AdminAuth(data.Authen, data.Remail) {
		files = []string{
			"templates/base.html",
			"templates/pages/adminarchive.html",
		}
	} else {
		files = []string{
			"templates/base.html",
			"templates/pages/archive.html",
		}
	}

	// Filter problems according to search parameter.
	// Creates a new slice specified by the filter and sets data.Problems to it
	if filterSearch != "" {
		filteredProblemsSearch := []ProblemInfo{}
		for _, p := range data.Problems {
			// Neglect case and leading/trailing spaces when filtering by search.
			if strings.Contains(strings.TrimSpace(strings.ToLower(p.Name)), strings.TrimSpace(strings.ToLower(filterSearch))) {
				filteredProblemsSearch = append(filteredProblemsSearch, p)
			}
		}
		data.Problems = filteredProblemsSearch
	}
	
	// Repeat filtering process for each filter.
	// Uses updated data.Problems with filter so that multiple filters work as intended
	if filterTag != "" {
		filteredProblemsTag := []ProblemInfo{}
		for _, p := range data.Problems {
			if contains(p.Tags, filterTag) {
				filteredProblemsTag = append(filteredProblemsTag, p)
			}
		}
		data.Problems = filteredProblemsTag
	}

	if filterDifficulty != "" {
		filteredProblemsDiff := []ProblemInfo{}
		for _, p := range data.Problems {
			if p.Difficulty == filterDifficulty {
				filteredProblemsDiff = append(filteredProblemsDiff, p)
			}
		}
		data.Problems = filteredProblemsDiff
	}

	if filterStatus != "" {
		filteredProblemsStat := []ProblemInfo{}
		for _, p := range data.Problems {
			if p.Status == filterStatus {
				filteredProblemsStat = append(filteredProblemsStat, p)
			}
		}
		data.Problems = filteredProblemsStat
	}

	// Parse the HTML templates.
	ts, err := template.ParseFiles(files...)

	// If there is an error, log it and send a 500 Internal Server Error response.
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// If the request method is not POST, execute the template with the data and return.
	if r.Method != http.MethodPost {
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	// Otherwise, there is user input
	// Retrieve form values signifying problem addition or deletion.
	addProb := r.FormValue("add")
	deleteProb := r.FormValue("delete")

	// Open a connection to the PostgreSQL database.
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Print(err.Error())
		return
	}

	// Form value is "add", meaning user requests to add problem to the database.
	if addProb == "add" {
		// Handles the addition of a new problem into the database table problems_lists.
		// It retrieves the problem's details from the form.
		// Name, difficulty, tag directly stored as strings; leading/trailing spaces neglected
		// Description, input, output, sample solution are read in as text files, then stored as strings for easier manipulation in database
		problemName := strings.TrimSpace(r.FormValue("name"))
		problemDifficulty := r.FormValue("difficulty")
		problemTags := strings.TrimSpace(r.FormValue("tag"))

		r.ParseMultipartForm(10 << 20)

		descriptionFile, _, err := r.FormFile("description")
		if err != nil {
			log.Fatal(err)
		}
		defer descriptionFile.Close()
		inFile, _, err := r.FormFile("input")
		if err != nil {
			log.Fatal(err)
		}
		defer inFile.Close()
		outFile, _, err := r.FormFile("output")
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
		solutionFile, _, err := r.FormFile("solution")
		if err != nil {
			log.Fatal(err)
		}
		defer solutionFile.Close()

		problemDescriptionBytes, err := ioutil.ReadAll(descriptionFile)
		if err != nil {
			log.Fatal(err)
		}
		problemInputBytes, err := ioutil.ReadAll(inFile)
		if err != nil {
			log.Fatal(err)
		}
		problemOutputBytes, err := ioutil.ReadAll(outFile)
		if err != nil {
			log.Fatal(err)
		}
		problemSolutionBytes, err := ioutil.ReadAll(solutionFile)
		if err != nil {
			log.Fatal(err)
		}
	
		// Trimming trailing spaces is only needed for I/O
		// Needed for output for answer comparison
		// Not really needed for input, but it makes the input seem less weird to the user
		// Description and solution rely on the problem writer to properly format it
		problemDescription := string(problemDescriptionBytes)
		problemInput := strings.TrimSpace(string(problemInputBytes))
		problemOutput := strings.TrimSpace(string(problemOutputBytes))
		problemSolution := string(problemSolutionBytes)
	
		// Creation of problem with the information inputted by user
		newProblem := ProblemsList{
			Name:        problemName,
			Difficulty:  problemDifficulty,
			Tags:        problemTags,
			Description: problemDescription,
			Answer:      problemOutput,
			Input:       problemInput,
			Solution:    problemSolution,
		}
		
		// Inserts problem into database
		// Does not return here after the error so the website doesn't crash
		// If user tries inserting duplicate problem
		result := db.Create(&newProblem)
		if result.Error != nil {
			log.Printf("Error inserting new problem: %s", result.Error.Error())
		}

		dbcon, err := db.DB()
		if err != nil {
			log.Print(err.Error())
			return
		}
		dbcon.Close()

		// Refresh page
		http.Redirect(w, r, "/archive", http.StatusSeeOther)
		return
	} else if deleteProb == "delete" {
		// If the "delete" form value is set, delete a problem from the database as requested by user.
		// Gets to-be-deleted problem name.
		// Searches for the problem in the database using name.
		// Deletes the problem.
		delProbN := r.FormValue("delPName")
		var delProb ProblemsList
		db.Where("name = ?", delProbN).First(&delProb)
		db.Where("name = ?", delProbN).Delete(&delProb)

		dbcon, err := db.DB()
		if err != nil {
			log.Print(err.Error())
			return
		}
		dbcon.Close()

		http.Redirect(w, r, "/archive", http.StatusSeeOther)
		return
	}

	// If neither "add" nor "delete" was set, just execute the template with the data.
	// Technically, if neither of the options do happen then the method is not POST, but just in case
	// if there are future additions that allow a post method with neither add nor delete, the page doesn't crash
	dbcon, err := db.DB()
	if err != nil {
		log.Print(err.Error())
		return
	}
	dbcon.Close()

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
