package main

import (
	"fmt"
	"sync"
)

func main() {
	students := []Student{
		Student{Name: "Sophie"},
		Student{Name: "Ben"},
	}
	result := RegisterStudents(students, Course{Name: "Intro to Golang"})
	fmt.Print(result)
}

type RegisterStudentsResults struct {
	Results []StudentRegistrationResult
}

type StudentRegistrationResult struct {
	Registration StudentRegistration
	Error        error
}

type StudentRegistration struct {
	Student Student
	Course  Course
}

type Course struct {
	Name string
}

type Student struct {
	Name string
}

func RegisterStudents(students []Student, course Course) RegisterStudentsResults {
	output := make(chan RegisterStudentsResults)
	input := make(chan StudentRegistrationResult)
	var wg sync.WaitGroup
	go handleResults(input, output, &wg)
	defer close(output)
	for _, student := range students {
		wg.Add(1)
		go ConcurrentRegisterStudent(student, course, input)
	}

	wg.Wait()
	close(input)
	return <-output
}

func handleResults(input chan StudentRegistrationResult, output chan RegisterStudentsResults, wg *sync.WaitGroup) {
	var results RegisterStudentsResults
	for result := range input {
		results.Results = append(results.Results, result)
		wg.Done()
	}
	output <- results
}

func ConcurrentRegisterStudent(student Student, course Course, output chan StudentRegistrationResult) {
	result := RegisterStudent(student, course)
	output <- result
}

func RegisterStudent(student Student, course Course) StudentRegistrationResult {
	return StudentRegistrationResult{
		Registration: StudentRegistration{
			Student: student,
			Course:  course,
		},
	}
}
