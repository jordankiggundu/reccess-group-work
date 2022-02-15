package main


import (
    "fmt"
    "log"
	"flag"
	"embed"
	"io/fs"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	uc "github.com/kindercare2022/kindercare-server/usecases"
)

//go:embed client/web/*
var static embed.FS

//go:embed client/web/admin/*
var adminstatic embed.FS

//go:embed client/web/teacher/*
var teacherstatic embed.FS

//go:embed client/web/pupil/*
var pupilstatic embed.FS

//go:embed client/web/status/*
var statusstatic embed.FS

//go:embed client/web/assignmentdetails/*
var assignmentdetailsstatic embed.FS

//go:embed client/web/comment/*
var commentstatic embed.FS

//go:embed client/web/pupilslist/*
var pupilslisting embed.FS

func index(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(struct{Success string}{Success: "API home"})
}

func resourceNotFound(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(struct{Success string}{Success: "The API doesn't have what you are looking for !"})
}

func getRouter() *mux.Router {
	website, _ := fs.Sub(static, "client/web")
	adminpage, _ := fs.Sub(adminstatic, "client/web/admin")
	pupilpage, _ := fs.Sub(pupilstatic, "client/web/pupil")
	statuspage, _ := fs.Sub(statusstatic, "client/web/status")
	teacherpage, _ := fs.Sub(teacherstatic, "client/web/teacher")
	commentpage, _ := fs.Sub(commentstatic, "client/web/comment")
	pupilslistt, _ := fs.Sub(pupilslisting, "client/web/pupilslist")
	assignmentdetailspage, _ := fs.Sub(assignmentdetailsstatic, "client/web/assignmentdetails")
	router := mux.NewRouter()
	//---
	router.HandleFunc("/user/login", uc.UserLogin).Methods("POST")
	router.HandleFunc("/pupils/list", uc.ReadAllPupils).Methods("GET")
	router.HandleFunc("/register/pupil", uc.RegisterPupil).Methods("POST")
	router.HandleFunc("/register/teacher", uc.RegisterTeacher).Methods("POST")
	router.HandleFunc("/activate/me", uc.SendActivationRequest).Methods("POST")
	router.HandleFunc("/activate/pupil/{pupilid}", uc.ActvatePupil).Methods("POST")
	router.HandleFunc("/assignment/list", uc.SubmitAssignmentsList).Methods("POST")
	router.HandleFunc("/comments/{pupilid}", uc.ReadTeachersComments).Methods("GET")
	router.HandleFunc("/create/assignment", uc.SubmitAssignmentsList).Methods("POST")
	router.HandleFunc("/solution/submit", uc.RecordAssignmentSolution).Methods("POST")
	router.HandleFunc("/deactivate/pupil/{pupilid}", uc.DeactvatePupil).Methods("POST")
	router.HandleFunc("/activation/requests", uc.ReadActivationRequests).Methods("GET")
	router.HandleFunc("/assignments/list/all", uc.ReadAllAssignmentsList).Methods("GET")
	router.HandleFunc("/comment/{solutionid}", uc.AttachCommentOnSolution).Methods("POST")
	router.HandleFunc("/assignments/list/{from}/{to}", uc.ReadAssignmentsList).Methods("GET")
	router.HandleFunc("/assignments/viewall/{pupilid}", uc.ViewAllAssignments).Methods("GET")
	router.HandleFunc("/assignment/attempt/{assignmentListid}", uc.AttemptAssignmentList).Methods("GET")
	router.HandleFunc("/solutions/{assignmentlistid}", uc.ReadSubmittedAssignmentSolutions).Methods("GET")
	router.HandleFunc("/assignment/viewone/{assignmentid}/{pupilid}", uc.ViewOneAssignmentDetails).Methods("GET")
	//pages
	router.PathPrefix("/admin").Handler(http.StripPrefix("/admin/", http.FileServer(http.FS(adminpage))) ).Methods("GET")    
	router.PathPrefix("/teacher").Handler(http.StripPrefix("/teacher/", http.FileServer(http.FS(teacherpage))) ).Methods("GET")    
	router.PathPrefix("/pupil").Handler(http.StripPrefix("/pupil/", http.FileServer(http.FS(pupilpage))) ).Methods("GET")
	router.PathPrefix("/status").Handler(http.StripPrefix("/status/", http.FileServer(http.FS(statuspage))) ).Methods("GET")    
	router.PathPrefix("/assignmentdetails").Handler(http.StripPrefix("/assignmentdetails/", http.FileServer(http.FS(assignmentdetailspage))) ).Methods("GET")    
	router.PathPrefix("/comment").Handler(http.StripPrefix("/comment/", http.FileServer(http.FS(commentpage))) ).Methods("GET")    
	router.PathPrefix("/list/pupils").Handler(http.StripPrefix("/list/pupils/", http.FileServer(http.FS(pupilslistt))) ).Methods("GET") 
	router.HandleFunc("/", index ).Methods("POST")
	router.PathPrefix("/").Handler( http.FileServer(http.FS(website)) ).Methods("GET")
	//Not found
	router.NotFoundHandler = http.HandlerFunc(resourceNotFound)
	
	return router
}

func main() {
    //++++| os.Args |+++++
    wsEndPoint := ":6200" 
    addr := flag.String("addr", wsEndPoint, "KinderCare API service address") 
    flag.Parse()
    //++++++++++++++++++++
    uc.Init()
    
    fmt.Println("Server listening on port: "+(strings.Split(wsEndPoint,":")[1])) 
    log.Fatal(http.ListenAndServe(*addr, getRouter()))
}








