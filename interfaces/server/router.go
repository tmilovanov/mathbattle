package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mathbattle/infrastructure"
	"mathbattle/interfaces/server/handlers"

	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("NOT FOUND HANDLER FOR URL: %v", r.URL)

	w.WriteHeader(http.StatusNotFound)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mathbattle")
}

func Start(container infrastructure.Container) {
	myRouter := mux.NewRouter()

	myRouter.NotFoundHandler = http.HandlerFunc(notFound)

	// Home Page
	myRouter.Handle("/", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(homePage)))

	// Stat
	sh := handlers.StatHandler{Ss: container.StatService()}
	myRouter.HandleFunc("/stat", sh.Stat)

	// Rounds
	rh := handlers.RoundHandler{Rs: container.RoundService()}
	myRouter.HandleFunc("/rounds/start", rh.StartNew).Methods("POST")
	myRouter.HandleFunc("/rounds/start_review", rh.StartReviewStage).Methods("POST")
	myRouter.HandleFunc("/rounds", rh.GetAll).Methods("GET")
	myRouter.HandleFunc("/rounds/running", rh.GetRunning).Methods("GET")
	myRouter.HandleFunc("/rounds/review_pending", rh.GetReviewPending).Methods("GET")
	myRouter.HandleFunc("/rounds/review_running", rh.GetReviewRunning).Methods("GET")
	myRouter.HandleFunc("/rounds/last", rh.GetLast).Methods("GET")
	myRouter.HandleFunc("/rounds/review_stage_distribution", rh.GetReivewStageDistribution).Methods("GET")
	myRouter.HandleFunc("/rounds/problem_descriptors/{participant_id}", rh.GetProblemDescriptors).Methods("GET")
	myRouter.HandleFunc("/rounds/{id}", rh.GetByID).Methods("GET")

	// Participants
	ph := handlers.ParticipantHandler{Ps: container.ParticipantService()}
	myRouter.Handle("/participants", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(ph.Store))).Methods("POST")
	myRouter.Handle("/participants", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(ph.GetAll))).Methods("GET")
	myRouter.Handle("/participants/{id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(ph.GetByID))).Methods("GET")
	myRouter.Handle("/participants/telegram/{id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(ph.GetByTelegramID))).Methods("GET")
	myRouter.Handle("/participants/{id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(ph.Update))).Methods("PUT")
	myRouter.Handle("/participants/{id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(ph.Delete))).Methods("DELETE")
	myRouter.Handle("/participants/unsubscribe/{id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(ph.Unsubscribe))).Methods("POST")

	// Solutions
	slh := handlers.SolutionHandler{Ss: container.SolutionService()}
	myRouter.Handle("/solutions", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(slh.Create))).Methods("POST")
	myRouter.Handle("/solutions/{id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(slh.GetByID))).Methods("GET")
	myRouter.Handle("/solutions/find/descriptor", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(slh.Find))).Methods("GET")
	myRouter.Handle("/solutions/append_part/{id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(slh.AppendPart))).Methods("POST")
	myRouter.Handle("/solutions/{id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(slh.Delete))).Methods("DELETE")
	myRouter.Handle("/solutions/descriptors/{participant_id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(slh.GetProblemDescriptors))).Methods("GET")

	// Reviews
	rs := handlers.ReviewHandler{Rs: container.ReviewService()}
	myRouter.Handle("/reviews", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(rs.Create))).Methods("POST")
	myRouter.Handle("/reviews/find/descriptor", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(rs.FindMany))).Methods("GET")
	myRouter.Handle("/reviews/{id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(rs.Delete))).Methods("DELETE")
	myRouter.Handle("/reviews/descriptors/{participant_id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(rs.GetSolutionDescriptors))).Methods("GET")

	// Problems
	prh := handlers.ProblemHandler{Ps: container.ProblemService()}
	myRouter.Handle("/problems/{id}", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(prh.GetByID))).Methods("GET")

	// Postman
	psth := handlers.PostmanHandler{Ps: container.Postman()}
	myRouter.Handle("/postman/send_to_users", ghandlers.LoggingHandler(os.Stdout, http.HandlerFunc(psth.SendToUsers)))

	log.Fatal(http.ListenAndServe(container.Config().APIUrl, myRouter))
}
