package usecases


import (
    "fmt"
    "errors"
    "strconv"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/kindercare2022/kindercare-server/db"
)

func RegisterTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher db.Teacher
	_ = json.NewDecoder(r.Body).Decode(&teacher)
	fmt.Printf("teacher.Profile.Password= %s, \nteacher.Profile=%v \n\n ", teacher.Profile.Password,  teacher.Profile)
	teacher.Profile.UserType = "teacher"
	database.Save(&teacher)
	respondToClient(w, 201, teacher, "Teacher registered succeffully.") 
}

type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login requeat received")
	var credentials Credentials
	_ = json.NewDecoder(r.Body).Decode(&credentials)
	
    var user db.Profile
    result := database.Where("username = ? AND password = ?", credentials.Username, credentials.Password).First(&user)
    rows := result.RowsAffected
    if rows > 0 {
        fmt.Printf("\nSigned in succeffully.\n")
        fmt.Printf("user.Username= %s user.Password= %s\n", user.Username, user.Password)
        //if user.userType is pupil, get the pupil id
        if user.UserType == "pupil" {
            pupil := db.Pupil{}
            database.Where("profile_id = ?", user.ID).First(&pupil)
            respondToClient(w, 200, pupil, "Sign in successful.")
        }
        //else if user.userType = 'teacher' , get the teacher id
        fmt.Println("user = ",user)
        respondToClient(w, 200, user, "Sign in successful.")  
    }else{
        fmt.Println("Signed in failed.")
        respondToClient(w, 403, nil, "Access denied.")  
    }
}

func userExists (identifier string) (bool, db.Profile, error) {
    //the identifier can be ID, phone, email, username
    var user db.Profile
    response := database.Where("id = ? OR username = ?", identifier, identifier).First(&user)                   
    numberOfRowsFound := response.RowsAffected
    userExists := numberOfRowsFound > 0
    
    if !userExists {
        if id, err := strconv.Atoi(identifier); err == nil {
            resp := database.Where("id = ?", uint(id)).First(&user)
            rowsFound := resp.RowsAffected
            exists := rowsFound > 0
            return exists, user, response.Error
        }else{
            return false, user, errors.New("user id must be a number")
        } 
    }else{
        return userExists, user, response.Error
    }
}

func ReadUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    identifier := params["id"]
    
    ok, user, err := userExists(identifier)
    if err != nil {
        respondToClient(w, 400, nil, err.Error())
    }
    
    if !ok {
        respondToClient(w, 404, nil, "Specified User not found")
    }
    
    respondToClient(w, 200, user, "")
}

func ReadAllUsers(w http.ResponseWriter, r *http.Request) {
    var users []db.Profile
    response := database.Find(&users)
    numberOfRowsFound := response.RowsAffected
    msg := fmt.Sprintf("Found %d users", numberOfRowsFound)
    respondToClient(w, 200, users, msg)
}




