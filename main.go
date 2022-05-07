package main

import (
	bst "GoSchool_Assignment2/binarysearchtree"
	dll "GoSchool_Assignment2/doublylinkedlist"
	util "GoSchool_Assignment2/utility"
	"fmt"
	"os"
)

func adminMenu(sessionList, patientList, dentistList **dll.DoublyLinkedlist, appointmentTree **bst.Binarysearchtree) {
	var optionSelected int

	fmt.Println("\nAdmin Function")
	fmt.Println("---------------")
	fmt.Println("[1] Browse dentist appointments")
	fmt.Println("[2] Search patients")
	fmt.Println("[3] Search appointment")
	fmt.Println("[4] Return to main menu")

selection:
	for {
		fmt.Println("\nSelect your choice:")
		if r, err := util.ReadInputAsInt(); err == nil {
			optionSelected = r
			break
		} else {
			fmt.Println(err)
		}
	}

	switch optionSelected {
	case 1:
		err := browseDoctorAppointment(sessionList, dentistList, appointmentTree)
		if err != nil {
			fmt.Println(err)
		}
	case 2:
		searchPatient(sessionList, patientList, appointmentTree)
	case 3:
		err := searchAppointment(sessionList, dentistList, patientList, appointmentTree)
		if err != nil {
			fmt.Println(err)
		}
	case 4:
		mainMenu(sessionList, patientList, dentistList, appointmentTree)
	case 10:
		os.Exit(0)
	default:
		fmt.Println("\nInvalid selection, please select a valid selection.")
		goto selection
	}

	fmt.Println("\nPress [Enter] to continue...")
	_, _ = fmt.Scanln()

	adminMenu(sessionList, patientList, dentistList, appointmentTree)
}

func mainMenu(sessionList, patientList, dentistList **dll.DoublyLinkedlist, appointmentTree **bst.Binarysearchtree) {
	var optionSelected int

	fmt.Println("\nCentral City Dentist Clinic")
	fmt.Println("=============================")
	fmt.Println("[1] Make an appointment")
	fmt.Println("[2] Search Dentist Schedule")
	fmt.Println("[3] Edit appointment")
	fmt.Println("[4] Admin")
	fmt.Println("[5] Exit")

selection:
	for {
		fmt.Println("\nSelect your choice:")
		if r, err := util.ReadInputAsInt(); err == nil {
			optionSelected = r
			break
		} else {
			fmt.Println(err)
		}
	}

	switch optionSelected {
	case 1:
		makeAppointment(sessionList, patientList, dentistList, appointmentTree)
	case 2:
		err := viewDentistSchedule(sessionList, dentistList, appointmentTree)
		if err != nil {
			fmt.Println(err)
		}
	case 3:
		editAppointment(sessionList, patientList, dentistList, appointmentTree)
	case 4:
		adminMenu(sessionList, patientList, dentistList, appointmentTree)
	case 5:
		os.Exit(0)
	default:
		fmt.Println("\nInvalid selection, please select a valid selection.")
		goto selection
	}

	fmt.Println("\nPress [Enter] to continue...")
	_, _ = fmt.Scanln()

	mainMenu(sessionList, patientList, dentistList, appointmentTree)
}

func main() {

	// Initialize new doubly linkedlist and binary search tree
	var (
		appointmentTree                       = bst.New()
		dentistList                           = dll.New()
		patientList                           = dll.New()
		sessionList                           = dll.New()
		dentist, dentist2                     *Dentist
		patient, patient2, patient3, patient4 *Patient
	)

	// Initialize Sample Data
	sessionList.Add(Session{1, "09:00", "10:00"})
	sessionList.Add(Session{2, "10:00", "11:00"})
	sessionList.Add(Session{3, "11:00", "12:00"})
	sessionList.Add(Session{4, "13:00", "14:00"})
	sessionList.Add(Session{5, "14:00", "15:00"})
	sessionList.Add(Session{6, "15:00", "16:00"})
	sessionList.Add(Session{7, "16:00", "17:00"})

	dentistList.Add(&Dentist{"Dr. James Holden"})
	dentistList.Add(&Dentist{"Dr. Amos Burton"})
	dentistList.Add(&Dentist{"Dr. Camina Drummer"})

	patientList.Add(&Patient{"Rosy Tyler", 85551064})
	patientList.Add(&Patient{"Rose Tyler", 92976084})
	patientList.Add(&Patient{"Donna Noble", 97251271})
	patientList.Add(&Patient{"Martha Jones", 87166096})
	patientList.Add(&Patient{"Amy Pond", 96402533})
	patientList.Add(&Patient{"River Song", 96418406})
	patientList.Add(&Patient{"Clara Oswald", 92309613})
	patientList.Add(&Patient{"Bill Potts", 82187412})
	patientList.Add(&Patient{"Yasmin Khan", 87866727})

	// Sort Patient first so binary searching will be faster
	chSort := make(chan error)
	go sortPatientsByMobileNumber(&patientList, chSort)
	err := <-chSort
	if err != nil {
		fmt.Println(err)
	}

	// Dr. James Holden
	dd, err := dentistList.Get(1)
	if err == nil {
		dentist = dd.(*Dentist)
	}

	// Dr. Amos Burton
	dd2, err2 := dentistList.Get(2)
	if err2 == nil {
		dentist2 = dd2.(*Dentist)
	}

	// Clara Oswald [92309613]
	pp, err := patientList.Get(5)
	if err == nil {
		patient = pp.(*Patient)
	}

	// Bill Potts [82187412]
	pp2, err := patientList.Get(1)
	if err == nil {
		patient2 = pp2.(*Patient)
	}

	// Amy Pond [96402533]
	pp3, err := patientList.Get(7)
	if err == nil {
		patient3 = pp3.(*Patient)
	}

	// Rosy Tyler [85551064]
	pp4, err := patientList.Get(2)
	if err == nil {
		patient4 = pp4.(*Patient)
	}

	// Insert sample data into binary search tree
	appointmentTree.Add("2022-05-19", 5, dentist, patient)
	appointmentTree.Add("2022-05-19", 2, dentist2, patient2)
	appointmentTree.Add("2022-05-19", 6, dentist, patient3)
	appointmentTree.Add("2022-05-01", 3, dentist, patient3)
	appointmentTree.Add("2022-05-21", 6, dentist, patient3)
	appointmentTree.Add("2022-05-23", 4, dentist, patient2)
	appointmentTree.Add("2022-05-19", 2, dentist, patient4)

	mainMenu(&sessionList, &patientList, &dentistList, &appointmentTree)
}
