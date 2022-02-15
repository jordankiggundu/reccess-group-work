package usecases


import (
    "fmt"
    "time"
    "errors"
    "strconv"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/kindercare2022/kindercare-server/db"
)

func ReadAssignmentsList(w http.ResponseWriter, r *http.Request) {
	var assignmentLists []db.AssignmentList
	params := mux.Vars(r)
    from  := params["from"]
    to    := params["to"]
    
    fmt.Println("from = ", from)
    fmt.Println("to = ", to)
	response := database.Where("created_at >= ? AND created_at <= ?", from, to).Find(&assignmentLists)
	if response.RowsAffected < 1 {
        respondToClient(w, 404, nil, "No assignments found")
        return
    }
    
	expiredAssingmentsCount := 0
	availableAssingmentsCount := 0
	for _, assignmentList := range assignmentLists {
	    if !assignmentList.Expired {
	        expiredAssingmentsCount++
	    }else{
	        availableAssingmentsCount++
	    }
	}
	msg := fmt.Sprintf("Found %d available, and %d expired assignments.", availableAssingmentsCount, expiredAssingmentsCount)
	respondToClient(w, 200, assignmentLists, msg)
}


func ReadAllAssignmentsList(w http.ResponseWriter, r *http.Request) {
	var assignmentLists []db.AssignmentList
	response := database.Where("expired = ?", false).Find(&assignmentLists)
	if response.RowsAffected < 1 {
        respondToClient(w, 404, assignmentLists, "No assignments found")
        return
    }
	msg := fmt.Sprintf("Found %d assignment lists.", response.RowsAffected)
	respondToClient(w, 200, assignmentLists, msg)
}

func pupilExists (identifier string) (bool, db.Pupil, error) {
    var pupil db.Pupil
    response := database.Where("id = ? OR profile_id = ?", identifier, identifier).First(&pupil)                   
    numberOfRowsFound := response.RowsAffected
    pupilExists := numberOfRowsFound > 0
    
    if !pupilExists {
        if id, err := strconv.Atoi(identifier); err == nil {
            resp := database.Where("id = ?", uint(id)).First(&pupil)
            rowsFound := resp.RowsAffected
            exists := rowsFound > 0
            return exists, pupil, response.Error
        }else{
            return false, pupil, errors.New("pupil id must be a number")
        } 
    }else{
        return pupilExists, pupil, response.Error
    }
}

func assignmentExists (identifier string) (bool, db.Assignment, error) {
    var assignment db.Assignment
    response := database.Where("id = ?", identifier).First(&assignment)                   
    numberOfRowsFound := response.RowsAffected
    assignmentExists := numberOfRowsFound > 0
    
    if !assignmentExists {
        if id, err := strconv.Atoi(identifier); err == nil {
            resp := database.Where("id = ?", uint(id)).First(&assignment)
            rowsFound := resp.RowsAffected
            exists := rowsFound > 0
            return exists, assignment, resp.Error
        }else{
            return false, assignment, errors.New("assignment id must be a number")
        } 
    }else{
        return assignmentExists, assignment, response.Error
    }
}

type assignmentsAttemptedOrNot struct {
	  ID                         uint      `json:"id"`
      AssignmentListId           string    `json:"assignment_list_id"`
      Assignment                 string    `json:"assignment"`
      AttemptTimeTaken           uint      `json:"attempt_time_taken"`
      Characterr                 string    `json:"characterr"`
      Score                      uint      `json:"score"`
	  PupilStatus                string    `json:"pupil_status"`
	  CreatedAt                  time.Time `json:"created_at"`
      Attempted                  string    `josn:"attempted"`
}

func ViewAllAssignments(w http.ResponseWriter, r *http.Request) {
	allAssignments := []assignmentsAttemptedOrNot{}
	attemptedAssignments := []assignmentsAttemptedOrNot{}
	unAttemptedAssignments := []assignmentsAttemptedOrNot{}
	params := mux.Vars(r)
    pupilid  := params["pupilid"]
    
    var solutions []db.Solution
    var listOfAsnmtIDs []uint
    
    ok, pupil, err := pupilExists(pupilid)
    if err != nil {
        respondToClient(w, 400, nil, err.Error())
        return
    }
    if !ok {
        respondToClient(w, 400, nil, "Specified pupil nt found")
        return
    }
    fmt.Printf("pupil.ID = %v\n", pupil.ID)
	_ = database.Table("solutions").Where("pupil_id = ?", pupil.ID).Select("assignment_id").Find(&solutions)
	
	fmt.Printf("solutions = %v\n", solutions)
	
	for _, solution := range solutions {
	    listOfAsnmtIDs = append(listOfAsnmtIDs, solution.AssignmentID)
	}
	
	printPupilStatus := func(condition bool) string {
	    if condition {
            return "blocked"
        }
        return "active"
	}
	
	response1 := database.Table("solutions as s").Joins("JOIN assignments as a on s.assignment_id=a.id").Select(`a.id, a.created_at, CONCAT(?) as pupil_status, CONCAT('attempted') as attempted, s.assignment_list_name as assignment, s.character_to_draw as characterr, s.auto_attached_score as score, s.attempt_time_taken_in_seconds as attempt_time_taken, s.assignment_list_id`, printPupilStatus(pupil.Deactivated) ).Find(&attemptedAssignments)    
	response2 := database.Table("assignments as a").Where("a.id NOT IN ?", listOfAsnmtIDs).Joins("JOIN assignment_lists as al on al.id=a.assignment_list_id").Select(`a.id, CONCAT(?) as pupil_status, a.created_at, a.character_to_draw as characterr, CONCAT("Not attempted") as attempted, al.name as assignment, al.id as assignment_list_id`, printPupilStatus(pupil.Deactivated) ).Find(&unAttemptedAssignments)    
	
	if (response1.RowsAffected < 1) && (response2.RowsAffected < 1) {
        respondToClient(w, 400, nil, "No assignments found")
        return
    }
    
    //eliminate duplicates from the list:
    
    
	msg := fmt.Sprintf("Found %d attempted, and %d not attempted.", response1.RowsAffected, response2.RowsAffected)
	//include all cassignments that are already attempted
	uniqueChars := []string{}
	for _, elem := range attemptedAssignments {
	    ch := elem.Characterr
	    alreadyIncluded := false
	    for _, c := range uniqueChars {
	        if  ch == c {
	            alreadyIncluded = true
	            break
	        }
	    }
	    
	    if !alreadyIncluded {
	        allAssignments = append(allAssignments, elem)
	        uniqueChars = append(uniqueChars, ch)
	    }
	}
	
	//include all cassignments that are not yet attempted
	uniqueChars2 := []string{}
	for _, elem := range unAttemptedAssignments {
	    ch := elem.Characterr
	    alreadyIncluded := false
	    for _, c := range uniqueChars2 {
	        if  ch == c {
	            alreadyIncluded = true
	            break
	        }
	    }
	    
	    if !alreadyIncluded {
	        allAssignments = append(allAssignments, elem)
	        uniqueChars2 = append(uniqueChars2, ch)
	    }
	}
	
	respondToClient(w, 200, allAssignments, msg)
}

func ViewOneAssignmentDetails (w http.ResponseWriter, r *http.Request) {
    attemptedAssignments := []assignmentsAttemptedOrNot{}
	unAttemptedAssignments := []assignmentsAttemptedOrNot{}
	params := mux.Vars(r)
    pupilid  := params["pupilid"]
    assignmentid  := params["assignmentid"]
    
    var solutions []db.Solution
    var attemptedAsnmtIDs []uint
    
    ok, pupil, err := pupilExists(pupilid)
    okk, assignmt, er := assignmentExists(assignmentid)
    if err != nil {
        respondToClient(w, 400, nil, err.Error())
        return
    }
    if er != nil {
        respondToClient(w, 400, nil, err.Error())
        return
    }
    if !ok {
        respondToClient(w, 400, nil, "Specified pupil not found")
        return
    }
    if !okk {
        respondToClient(w, 400, nil, "Specified assignment not found")
        return
    }
    
	_ = database.Table("solutions").Where("pupil_id = ? AND assignment_id = ?", pupil.ID, assignmt.ID).Select("assignment_id").Find(&solutions)
	for _, solution := range solutions {
	    attemptedAsnmtIDs = append(attemptedAsnmtIDs, solution.AssignmentID)
	}
	
	printPupilStatus := func(condition bool) string {
	    if condition {
            return "blocked"
        }
        return "active"
	}
	
	response1 := database.Table("solutions as s").Joins("JOIN assignments as a on s.assignment_id=a.id").Where("(a.id IN ?) AND (a.id = ?)", attemptedAsnmtIDs, assignmt.ID).Select(`a.id, a.created_at, CONCAT(?) as pupil_status, CONCAT('attempted') as attempted, s.assignment_list_name as assignment, s.character_to_draw as characterr, s.auto_attached_score as score, s.attempt_time_taken_in_seconds as attempt_time_taken, s.assignment_list_id`, printPupilStatus(pupil.Deactivated) ).Find(&attemptedAssignments)    
	if (response1.RowsAffected >= 1) {
        msg := fmt.Sprintf("Found %d record.", response1.RowsAffected)
        respondToClient(w, 200, attemptedAssignments[0], msg)
        return
    }
    
	response2 := database.Table("assignments as a").Joins("JOIN assignment_lists as al on al.id=a.assignment_list_id").Where("a.id = ?", assignmt.ID).Select(`a.id, CONCAT(?) as pupil_status, a.created_at, a.character_to_draw as characterr, CONCAT("Not attempted") as attempted, al.name as assignment, al.id as assignment_list_id`, printPupilStatus(pupil.Deactivated) ).Find(&unAttemptedAssignments)    
	if (response1.RowsAffected < 1) && (response2.RowsAffected < 1) {
        respondToClient(w, 400, nil, "No assignments found")
        return
    }
    
	msg := fmt.Sprintf("Found %d record.", response1.RowsAffected + response2.RowsAffected)
	respondToClient(w, 200, unAttemptedAssignments[0], msg)
}

func AttemptAssignmentList (w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var assignments []db.Assignment
    assignmentListid  := params["assignmentListid"]
    fmt.Println("Attempt: assignmentListid = ", assignmentListid)
    if id, err := strconv.Atoi(assignmentListid); err == nil {
        resp := database.Where("assignment_list_id = ?", uint(id)).Select("id, created_at, assignment_list_id, character_to_draw").Find(&assignments)
        rowsFound := resp.RowsAffected
        msg := fmt.Sprintf("Assignents available = %d", rowsFound)
        respondToClient(w, 200, assignments, msg)
    }else{
        respondToClient(w, 400, nil, "assignment list id must be a number")
    }
}

func RecordAssignmentSolution(w http.ResponseWriter, r *http.Request) {
    var solutions []db.Solution
    var asnmtlist db.AssignmentList
    _ = json.NewDecoder(r.Body).Decode(&solutions)
    //retreive the target assignmentList, get its name and assign to solution
    if len(solutions) > 0 {
      soln := solutions[0]
      database.Where("id = ?", soln.AssignmentListID).First(&asnmtlist);
    }
    for i, solution := range solutions {
      solution.AssignmentListName = asnmtlist.Name
      solutions[i] = solution
    }
	database.Create(&solutions)
	respondToClient(w, 201, solutions, "Answers submitted succeffully.") 
}

func ReadTeachersComments(w http.ResponseWriter, r *http.Request) {
	var comments []db.Solution
	params := mux.Vars(r)
    pupilid  := params["pupilid"]
	response := database.Where("pupil_id = ? AND comment != ?", pupilid, "").Find(&comments)
	if response.RowsAffected < 1 {
        respondToClient(w, 400, nil, "You have not yet attempted any assignments") 
    }
    msg := fmt.Sprintf("Found %d comments.", response.RowsAffected)
	respondToClient(w, 200, comments, msg)
}

func SendActivationRequest(w http.ResponseWriter, r *http.Request) {
	var arequest db.ActivationRequest
	_ = json.NewDecoder(r.Body).Decode(&arequest)
	database.Save(&arequest)
	respondToClient(w, 201, arequest, "Request submitted.") 
}





