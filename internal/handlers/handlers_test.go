package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuyoshio/bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	//{"post-search-avail", "/search-availability", "POST", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-02"},
	//}, http.StatusOK},
	//{"post-search-avail-json", "/search-availability-json", "POST", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-02"},
	//}, http.StatusOK},
	//{"make-reservation-post", "/make-reservation", "POST", []postData{
	//	{key: "first_name", value: "Liu"},
	//	{key: "last_name", value: "Yoshio"},
	//	{key: "email", value: "me@me.com"},
	//	{key: "phone", value: "0000000"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, exptected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler return wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// testcase where reservation is not in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// test with non-existing room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "first_name=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=abc@abc.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// put fake reservation in session
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	session.Put(ctx, "reservation", reservation)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong status code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// test for missing post body
	reqBody = "first_name=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=abc@abc.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")

	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong status code for missing post body: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// test for missing reservation in session
	// test for missing post body
	reqBody = "first_name=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=abc@abc.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong status code for missing post body: got %d, expected %d", rr.Code, http.StatusOK)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// first-case rooms are not available
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	log.Println(reqBody)
	// create the request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("fail to parse json")
	}
	log.Println(j)
}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
