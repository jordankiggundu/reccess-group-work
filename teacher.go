package usecases


import (
    "fmt"
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


