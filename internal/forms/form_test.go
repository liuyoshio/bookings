package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b")

	if form.Valid() {
		t.Error("got valid when should have been invalid")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b")

	if !form.Valid() {
		t.Error("get invalid when should have been invalid")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "")
	r.PostForm = postedData

	err := r.ParseForm()
	if err != nil {
		t.Error("error when parsing form")
	}

	form := New(r.PostForm)

	if !form.Has("a", r) {
		t.Error("got false when should have been true")
	}

	if form.Has("b", r) {
		t.Error("got true when should have been false")
	}

	if form.Has("c", r) {
		t.Error("got true when should have been false")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	postedData := url.Values{}
	postedData.Add("name1", "aaaa")
	postedData.Add("name2", "s")
	r.PostForm = postedData

	err := r.ParseForm()
	if err != nil {
		t.Error("error when parsing form")
	}

	form := New(r.PostForm)

	if !form.MinLength("name1", 3, r) {
		t.Error("got false when should have true")
	}

	if form.MinLength("name2", 3, r) {
		t.Error("got true when should have false")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	postedData := url.Values{}
	postedData.Add("email1", "aaa@a.com")
	postedData.Add("email2", "abc")
	r.PostForm = postedData

	err := r.ParseForm()
	if err != nil {
		t.Error("error when parsing form")
	}
	form := New(r.PostForm)

	if !form.IsEmail("email1") {
		t.Error("got false when should have true")
	}

	if form.IsEmail("email2") {
		t.Error("got true when should have false")
	}

}
