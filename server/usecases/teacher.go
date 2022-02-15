package usecases


import (
    "fmt"
    "time"
    "strconv"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/kindercare2022/kindercare-server/db"
)


func RegisterPupil(w http.ResponseWriter, r *http.Request) {
	var pupil db.Pupil
	_ = json.NewDecoder(r.Body).Decode(&pupil)
	pupil.Profile.UserType = "pupil"
	database.Save(&pupil)
	respondToClient(w, 201, pupil, "Pupil registered succeffully.") 
}

func DeactvatePupil(w http.ResponseWriter, r *http.Request) {
	var pupil db.Pupil
	params := mux.Vars(r)
    pupilid := params["pupilid"]
	id, _ := strconv.Atoi(pupilid);
	database.First(&pupil, "id = ?", id)
	pupil.Deactivated = true
	database.Save(&pupil)
	respondToClient(w, 201, pupil, "Deactivated.") 
}

type AcvnRequest struct {
    db.ActivationRequest
    db.Pupil
}

func ReadActivationRequests(w http.ResponseWriter, r *http.Request) {
	var requests []AcvnRequest
	response := database.Table("activation_requests as ar").Joins("JOIN pupils as p on p.id=ar.pupil_id").Select("ar.id, ar.pupil_id, ar.request_message, p.firstname, p.lastname, p.deactivated").Find(&requests)
	if response.RowsAffected < 1 {
        respondToClient(w, 200, requests, "No activation requests found.") 
    }
    msg := fmt.Sprintf("Found %d requests.", response.RowsAffected)
	respondToClient(w, 200, requests, msg)
}

func ActvatePupil(w http.ResponseWriter, r *http.Request) {
	var pupil db.Pupil
	params := mux.Vars(r)
    pupilid := params["pupilid"]
    id, _ := strconv.Atoi(pupilid);
	database.First(&pupil, "id = ?", id)
	pupil.Deactivated = false
	database.Save(&pupil)
	respondToClient(w, 201, pupil, "Activated.")
}

func SubmitAssignmentsList(w http.ResponseWriter, r *http.Request) {
    var assignmentList db.AssignmentList
	_ = json.NewDecoder(r.Body).Decode(&assignmentList)
	database.Save(&assignmentList)
	respondToClient(w, 201, assignmentList, "Assignment submited succeffully.")
}

func ReadSubmittedAssignmentSolutions(w http.ResponseWriter, r *http.Request) {
	var solutions []db.Solution
	params := mux.Vars(r)
    assignmentlistid  := params["assignmentlistid"]
	response := database.Table("solutions as s").Joins("JOIN pupils as p on p.id=s.pupil_id").Where("s.assignment_list_id = ? OR s.assignment_list_name = ?", assignmentlistid, assignmentlistid).Select("s.*, p.firstname, p.lastname").Find(&solutions)
	if response.RowsAffected < 1 {
        respondToClient(w, 400, solutions, "The assignment you've specified does not exist.") 
    }
    msg := fmt.Sprintf("Found %d submissions.", response.RowsAffected)
	respondToClient(w, 200, solutions, msg)
}

func AttachCommentOnSolution(w http.ResponseWriter, r *http.Request) {
	var solution db.Solution
	var commented db.Solution
	_ = json.NewDecoder(r.Body).Decode(&commented)
	fmt.Println("received = ", solution)
	params := mux.Vars(r)
    solutionid := params["solutionid"]
	database.First(&solution, "id = ? OR id = ?", solutionid, commented.ID)
	solution.Comment = commented.Comment
	database.Save(&solution)
	respondToClient(w, 201, solution, "Comment attached.") 
}

func ReadAllPupils (w http.ResponseWriter, r *http.Request) {
    var pupils []db.Pupil
    response := database.Find(&pupils)
    msg := fmt.Sprintf("Found %d records", response.RowsAffected)
    respondToClient(w, 201, pupils, msg) 
}

type PupilSolutions struct {
      ID                         uint      `json:"id"`
      AssignmentListId           string    `json:"assignment_list_id"`
      Firstname           string    `json:"firstname"`
      Lastname           string    `json:"lastname"`
      Phone           string    `json:"phone"`
      Assignment                 string    `json:"assignment"`
      AttemptTimeTaken           uint      `json:"attempt_time_taken"`
      Characterr                 string    `json:"characterr"`
      Score                      uint      `json:"score"`
      PupilStatus                string    `json:"pupilstatus"`
      Deactivated                bool      `json:"deactivated"`
      CreatedAt                  time.Time `json:"created_at"`
      Attempted                  string    `josn:"attempted"`
}

////////////////////////////////////////////////////////////////////////////////////
func ViewAllSolutionWithDetails(w http.ResponseWriter, r *http.Request) {
	allAssignments := []PupilSolutions{}
	attemptedAssignments := []PupilSolutions{}
	unAttemptedAssignments := []PupilSolutions{}
	params := mux.Vars(r)
        assignmentlistid  := params["assignmentlistid"]
    
        var solutions []db.Solution
        var listOfPupilIDs []uint
        alid, _ := strconv.Atoi(assignmentlistid)
        assignmentListName := "" 
        
	_ = database.Table("solutions").Where("assignment_list_id = ?", alid).Select("pupil_id, assignment_list_name").Find(&solutions)
	
	fmt.Printf("solutions under assignment-list id %d = %v\n", alid, solutions)
	
	for _, solution := range solutions {
	    pid := solution.PupilID
	    
	    //eliminate duplicates from the list:
	    alreadyIncluded := false
	    for _, idd := range listOfPupilIDs {
	        if  idd == pid {
	            alreadyIncluded = true
	            break
	        }
	    }
	    if !alreadyIncluded {
	        listOfPupilIDs = append(listOfPupilIDs, solution.PupilID)
	    	assignmentListName = solution.AssignmentListName
	    }
	}
	
	/*
	printPupilStatus := func(condition bool) string {
	    if condition {
            return "blocked"
        }
        return "active"
	}
	*/
	
	response1 := database.Table("solutions as s").Joins("JOIN assignments as a on s.assignment_id=a.id").Joins("JOIN pupils as p on s.pupil_id=p.id").Select(`a.id, a.created_at, p.lastname, p.firstname, p.phone, CONCAT('Yes') as attempted, s.assignment_list_name as assignment, s.character_to_draw as characterr, s.auto_attached_score as score, s.attempt_time_taken_in_seconds as attempt_time_taken, s.assignment_list_id`).Find(&attemptedAssignments)
	response2 := database.Table("pupils as p").Where("p.id NOT IN ?", listOfPupilIDs).Select(`p.firstname, p.lastname, p.phone, p.deactivated as deactivated, CONCAT("No") as attempted, CONCAT(?) as assignment,  CONCAT(?) as assignment_list_id`, assignmentListName, assignmentlistid ).Find(&unAttemptedAssignments)   
	
	if (response1.RowsAffected < 1) && (response2.RowsAffected < 1) {
        respondToClient(w, 400, allAssignments, "No solutions found")
        return
    	}
    
	msg := fmt.Sprintf("Found %d attempted, and %d not attempted.", response1.RowsAffected, response2.RowsAffected)
	
	allAssignments = append(allAssignments, attemptedAssignments...)
	allAssignments = append(allAssignments, unAttemptedAssignments...)
	
	respondToClient(w, 200, allAssignments, msg)
}

/////////////////////////////////////////////////////////////////////////////////
func ReadAllAssingmentLists(w http.ResponseWriter, r *http.Request) {
	assnmntLists := []db.AssignmentList{}
	response := database.Preload("Assignments").Find(&assnmntLists)
	msg := fmt.Sprintf("Found %d assignments", response.RowsAffected)
	respondToClient(w, 200, assnmntLists, msg)
}

/////////////////////////////////////////////////////////////////////////////////




