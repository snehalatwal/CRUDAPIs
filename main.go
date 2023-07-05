package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Models->different file

type Course struct {
	CourseID    string  `json:"course_id"`
	CourseName  string  `json:"course_name"`
	CoursePrice int     `json:"course_price"`
	Author      *Author `json:"course_author"`
}

type Author struct {
	FullName string `json:"author_name"`
	Website  string `json:"author_website"`
}

// fake DB-file
var courses []Course

func (c *Course) isEmpty() bool {
	// return true when course_name and course_id is missing
	// return c.CourseName == "" && c.CourseID == ""
	return c.CourseName == ""
}

// serveHome route
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to API by Sneha</h1>"))
}

// controllers

// get all courses route
func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all the courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

//to get a course from course id

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Return the course")
	w.Header().Set("Content-Type", "application/json")

	// query parameters passed
	params := mux.Vars(r)

	// loop through courses, find matching id, return the response
	for _, course := range courses {
		if course.CourseID == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	json.NewEncoder(w).Encode("No course found with given id")
	return
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create one course")
	w.Header().Set("Content-Type", "application/json")

	// if request body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please pass the body")
	}

	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)

	if course.isEmpty() {
		json.NewEncoder(w).Encode("No data inside the JSON passef")
		return
	}

	// generate unique id and convert them to string
	// append course to courses

	rand.Seed(time.Now().UnixNano())
	course.CourseID = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return

}

func updateTheCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update the course")
	w.Header().Set("Content-type", "application/json")

	// grab id from request
	params := mux.Vars(r)

	// loop, id ,remove, add with my Id(update)
	for index, course := range courses {
		if course.CourseID == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseID = params["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	// send the response when id is not found
	json.NewEncoder(w).Encode("No id passed")
	return

}

func deleteCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete the course")
	w.Header().Set("Content-type", "application/json")

	params := mux.Vars(r)

	// loop, id, remove, index, index+1

	for index, course := range courses {
		if course.CourseID == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode("The course is deleted")
			break
		}
	}

}

func main() {
	fmt.Println("BuildAPI")
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateTheCourse).Methods("PUT")
	r.HandleFunc("/course/{id}", deleteCourse).Methods("DELETE")

	// seeding
	courses = append(courses, Course{CourseID: "1", CourseName: "NodeJS", CoursePrice: 599, Author: &Author{

		FullName: "Sneha Latwal", Website: "sneha.dev",
	}})
	courses = append(courses, Course{CourseID: "2", CourseName: "Golang", CoursePrice: 199, Author: &Author{

		FullName: "Sneha", Website: "sneha/dev.in",
	}})

	log.Fatal(http.ListenAndServe(":4000", r))

}
