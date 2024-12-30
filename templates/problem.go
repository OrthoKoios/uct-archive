// File: problem.go
// Authors: Charisman Aravinthan and Justin Dai
// Last Modified: June 2, 2023
// Version: 1.1.0
// This file contains the backend for the problem page

// Problem Page
// Allows user to view specific problem
// Allows user to attempt problem and receive answer feedback
// Allows user to view answer and solution
package backend

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// Function to display specific problem. The function takes an HTTP response writer and an HTTP request as arguments.
// Returns None.
// Contributers: Justin Dai, Charisman Aravinthan.
func Problem(w http.ResponseWriter, r *http.Request) {

	// Retrieve data state and set to 0; no response from page required.
	data := GetData(w, r)
	data.State = "0"

	// Retrieve problem name variable from URL.
	vars := mux.Vars(r)
	problemName := vars["problem_name"]

	files := []string{
		"templates/base.html",
		"templates/pages/problempage.html",
	}

	// Find problem in database using problem name.
    var currProb *ProblemInfo
    for _, p := range data.Problems {
        if p.Name == problemName {
            currProb = &p
            break
        }
    }

	// If problem cannot be found, redirects to 404 page.
    if currProb == nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
        return
    }

	// Updates data problem name, description, answer, input, solution values to the respective values of the problem user requests.
	data.Pname = currProb.Name
	data.Pdescription = currProb.Description
	data.Panswer = currProb.Answer
	data.Pinput = currProb.Input
	data.Psolution = currProb.Solution

	// Connection to database.
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Print(err.Error())
		return
	}

	// Retrieves the requested problem from database for updating.
	// Although the problem was located earlier, it does not allow editing to the problem in the database because it is searched from .Problems.
	// Meanwhile, it is also necessary to have the ProblemInfo struct to display the info on the webpage.
	// As a result, a search is required in both .Problems and in the problems database.
	problem := ProblemsList{}
	result := db.First(&problem, "Name = ?", data.Pname)
	if result.Error != nil {
		log.Printf("Error finding problem: %s", result.Error.Error())
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {

		dbcon, derr := db.DB()
		if derr != nil {
			log.Print(err.Error())
			return
		}
		dbcon.Close()

		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// If the request method is not POST, execute the template with the data and return.
	if r.Method != http.MethodPost {

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
		return
	}

	// Otherwise, there is user input. Retrieves user answer from form.
	// Trimspace to allow user to add random spaces without getting wrong answer.
	userAns := strings.TrimSpace(r.FormValue("answer"))

	// Checks user submission and updates data state based on the submission
	// data state tells webpage how to respond to user
	if data.Username == "" {
		// Does not allow user to submit without being logged in
		data.State = "4"
	} else if userAns == "" {
		// Does not allow user to submit nothing
		data.State = "1"
	} else if userAns == data.Panswer {
		// User answer is correct
		// Adds user to list of users who correctly solved problem if user is not there already
		// Removes user to list of users who attempted the problem if user is there
		data.State = "2"
		foundc := index(currProb.CorrectUsers, data.Remail)
		foundw := index(currProb.WrongUsers, data.Remail)
		if foundc == -1 {
			currProb.CorrectUsers = append(currProb.CorrectUsers, data.Remail)
			result = db.Model(&problem).Where("Name = ?", data.Pname).Update("correct_users", strings.Join(currProb.CorrectUsers, " , "))
			if result.Error != nil {
				log.Printf("Error updating problem: %s", result.Error.Error())
			}
		}
		if foundw != -1 {
			currProb.WrongUsers = remove(currProb.WrongUsers, foundw)
			result = db.Model(&problem).Where("Name = ?", data.Pname).Update("wrong_users", strings.Join(currProb.WrongUsers, " , "))
			if result.Error != nil {
				log.Printf("Error updating problem: %s", result.Error.Error())
			}
		}
	} else if userAns != data.Panswer {
		// User answer is incorrect
		// Adds user to list of users who attempted the problem only if user has not solved the problem correctly yet
		data.State = "3"
		foundc := index(currProb.CorrectUsers, data.Remail)
		foundw := index(currProb.WrongUsers, data.Remail)
		if foundw == -1 {
			if foundc == -1 {
				currProb.WrongUsers = append(currProb.WrongUsers, data.Remail)
				result = db.Model(&problem).Where("Name = ?", data.Pname).Update("wrong_users", strings.Join(currProb.WrongUsers, " , "))
				if result.Error != nil {
					log.Printf("Error updating problem: %s", result.Error.Error())
				}
			}
		}
	}

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
