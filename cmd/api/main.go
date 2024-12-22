package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type Meeting struct {
    ID           int       `json:"id"`
    Title        string    `json:"title"`
    InviteeEmail string    `json:"invitee_email"`
    HostEmail    string    `json:"host_email"`
    StartAt      time.Time `json:"start_at"`
    EndAt        time.Time `json:"end_at"`
    Status       string    `json:"status"`
}

var (
    totalMeetings = promauto.NewCounter(prometheus.CounterOpts{
        Name: "meetings_total",
        Help: "Total number of meetings",
    })
    
    meetingsByStatus = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "meetings_by_status",
            Help: "Number of meetings by status",
        },
        []string{"status"},
    )
)

var db *sql.DB

func main() {
    connStr := "postgresql://nexu:nexu@postgres:5432/meetings?sslmode=disable"
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    r := mux.NewRouter()
    r.HandleFunc("/meetings", handleMeeting).Methods("POST")
    r.Handle("/metrics", promhttp.Handler())

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

func handleMeeting(w http.ResponseWriter, r *http.Request) {
    var meeting Meeting
    if err := json.NewDecoder(r.Body).Decode(&meeting); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var err error
    if meeting.ID != 0 {
        // Update existing meeting
        _, err = db.Exec(`
            UPDATE meetings 
            SET title=$1, invitee_email=$2, host_email=$3, start_at=$4, end_at=$5, status=$6 
            WHERE id=$7`,
            meeting.Title, meeting.InviteeEmail, meeting.HostEmail,
            meeting.StartAt, meeting.EndAt, meeting.Status, meeting.ID)
    } else {
        // Create new meeting
        err = db.QueryRow(`
            INSERT INTO meetings (title, invitee_email, host_email, start_at, end_at, status)
            VALUES ($1, $2, $3, $4, $5, $6)
            RETURNING id`,
            meeting.Title, meeting.InviteeEmail, meeting.HostEmail,
            meeting.StartAt, meeting.EndAt, meeting.Status).Scan(&meeting.ID)
        
        totalMeetings.Inc()
    }

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Update metrics
    meetingsByStatus.With(prometheus.Labels{"status": meeting.Status}).Inc()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(meeting)
}