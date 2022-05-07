package main

import (
	bst "GoSchool_Assignment2/binarysearchtree"
	dll "GoSchool_Assignment2/doublylinkedlist"
	util "GoSchool_Assignment2/utility"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

func makeAppointment(sessionList, patientList, dentistList **dll.DoublyLinkedlist, appointmentTree **bst.Binarysearchtree) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There was an error making a new appointment, returning to main menu....")
		}
	}()

	var (
		ok      bool
		patient *Patient
		dentist *Dentist
	)

	fmt.Println("\nMake an appointment")
	fmt.Println("---------------------")

patientInput:
	fmt.Println("\nPlease enter patient's mobile number:")
	patientMobile, err := util.ReadInputAsInt()
	if err == nil {
		if util.ValidateMobileNumber(patientMobile) {
			result := (**patientList).SearchByMobileNumber(patientMobile)
			if result != nil {
				patient = result.(*Patient)
				fmt.Println("Found patient:", patient.GetName())
			} else {
			yesNoInput:
				fmt.Println("\nPatient not found, create a new patient? (Y/N)")
				for {
					input := util.ReadInput()
					if input == "Y" || input == "y" {
						fmt.Println("\nPlease enter patient name")
						patientName := util.ReadInput()
						patient = &Patient{patientName, patientMobile}
						(**patientList).Add(patient)
						fmt.Printf("Patient [%v] created.\n", patient.GetName())
						//newPatient = true
						break
					} else if input == "N" || input == "n" {
						fmt.Println("\nPatient not created. Please enter another patient name.")
						goto patientInput
					} else {
						fmt.Println("\nInvalid Selection.")
						goto yesNoInput
					}
				}
			}
		} else {
			fmt.Println("Please enter a valid mobile number.")
		}
	} else {
		fmt.Println(err)
		goto patientInput
	}

	dentists := (**dentistList).GetList()
	printDentistList(dentists)

dentistInput:
	dentist, ok = getDentistFromSelection(dentists)
	if !ok {
		goto dentistInput
	}

apptInput:
	fmt.Println("\nPlease enter appointment date (yyyy-mm-dd):")
	inputDate := util.ReadInput()
	appointmentDate, err := time.Parse("2006-01-02", inputDate)
	if err != nil {
		fmt.Println(invalidDate)
		goto apptInput
	} else {
		currentTime := time.Now()
		if !currentTime.Before(appointmentDate) {
			fmt.Println(olderDate)
			goto apptInput
		}
	}

	schedule := (**appointmentTree).GetSchedule(appointmentDate.Format("2006-01-02"), dentist)
	printDentistSchedule(sessionList, appointmentDate.Format("2006-01-02"), dentist, &schedule)

sessionInput:
	fmt.Println("\nPlease enter session number [1-7]:")
	appointmentSession, err := util.ReadInputAsInt()
	if err == nil {
		for _, data := range schedule {
			if data.GetSession() == appointmentSession {
				fmt.Printf("Session [%v] is already booked, please select another session.\n", appointmentSession)
				goto sessionInput
			}
			if appointmentSession <= 0 || appointmentSession > 7 {
				fmt.Println("Invalid session, please select another session.")
				goto sessionInput
			}
		}
	} else {
		fmt.Println(err)
		goto sessionInput
	}

	chAppointment := make(chan error)
	chSort := make(chan error)

	go sortPatientsByMobileNumber(patientList, chSort)
	go addAppointment(appointmentTree, appointmentDate.Format("2006-01-02"), appointmentSession, dentist, patient, chAppointment)

	for i := 0; i < 2; i++ {
		select {
		case errSort := <-chSort:
			if errSort != nil {
				fmt.Println(errSort)
			}
		case errAdd := <-chAppointment:
			if errAdd != nil {
				fmt.Println(errAdd)
			} else {
				fmt.Println("\n--------------------")
				fmt.Println("Appointment created.")
				fmt.Println("--------------------")
			}
		}
	}
}

func editAppointment(sessionList, patientList, dentistList **dll.DoublyLinkedlist, appointmentTree **bst.Binarysearchtree) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("There was an error making an edit to appointment, returning to main menu....")
		}
	}()

	if i := (**patientList).GetSize(); i == 0 {
		panic(errors.New("\nerror: no patient"))
	}

	var (
		dateChanged    bool = false
		dentistChanged bool = false
		sessionChanged bool = false
		dentist        *Dentist
		patient        *Patient
	)

	fmt.Println("\nEdit appointment")
	fmt.Println("----------------")

patientInput:
	fmt.Println("\nPlease enter patient's mobile number:")
	patientMobile, err := util.ReadInputAsInt()
	if err == nil {
		if util.ValidateMobileNumber(patientMobile) {
			result := (**patientList).SearchByMobileNumber(patientMobile)
			if result != nil {
				patient = result.(*Patient)
				fmt.Println("Found existing patient:", patient.GetName())
			} else {
				fmt.Println("Patient does not exist, please enter correct patient's mobile number.")
				goto patientInput
			}
		} else {
			fmt.Println("Please enter a valid mobile number.")
		}
	} else {
		fmt.Println(err)
		goto patientInput
	}

	pSchedule := (**appointmentTree).GetUpComingSchedule(patient)

	if len(pSchedule) == 0 {
		fmt.Printf("\nThere are no appointment made for [%v]. Please make a new appointment first.\n", patient.GetName())
	} else {

		printPatientSchedule(sessionList, patient, pSchedule)

		var optionSelected int
		for {
			fmt.Println("\nPlease enter the S/N of the appointment which you would like to edit.")
			if r, err := util.ReadInputAsInt(); err == nil {
				optionSelected = r
				break
			} else {
				fmt.Println(err)
			}
		}

		node := pSchedule[optionSelected-1]
		dentists := (**dentistList).GetList()
		printDentistList(dentists)

	dentistInput:
		fmt.Println("\nPlease select a dentist. Enter 0 for no change:")
		dentistSelection, err := util.ReadInputAsInt()
		if err == nil {
			if dentistSelection > len(dentists) || dentistSelection < 0 {
				fmt.Println(invalidSelection)
				goto dentistInput
			} else if dentistSelection == 0 {
				dentist = node.GetDentist().(*Dentist)
			} else {
				dentistChanged = true
				dentist = dentists[dentistSelection-1].(*Dentist)
				fmt.Printf("Selected dentist [%s]\n", dentist.GetName())
			}
		} else {
			fmt.Println(err)
			goto dentistInput
		}

	apptInput:
		fmt.Println("\nPlease enter new appointment date (yyyy-mm-dd). Enter for no change.:")
		inputDate := util.ReadInput()

		var appointmentDate time.Time
		if len(inputDate) > 0 {
			appointmentDate, err = time.Parse("2006-01-02", inputDate)
			if err != nil {
				fmt.Println(invalidDate)
				goto apptInput
			} else {
				currentTime := time.Now()
				if !currentTime.Before(appointmentDate) {
					fmt.Println(olderDate)
					goto apptInput
				}
			}
			dateChanged = true
		} else {
			appointmentDate, _ = time.Parse("2006-01-02", node.GetDate())
		}

		dSchedule := (**appointmentTree).GetSchedule(appointmentDate.Format("2006-01-02"), dentist)
		printDentistSchedule(sessionList, appointmentDate.Format("2006-01-02"), dentist, &dSchedule)

	sessionInput:
		fmt.Println("\nPlease enter session number [1-7]:")
		appointmentSession, err := util.ReadInputAsInt()
		if err == nil {
			for _, data := range dSchedule {
				if data.GetSession() == appointmentSession {
					fmt.Printf("Session [%v] is already booked, please select another session.\n", appointmentSession)
					goto sessionInput
				}
				if appointmentSession <= 0 || appointmentSession > 7 {
					fmt.Println("Invalid session, please select another session.")
					goto sessionInput
				}
			}
			sessionChanged = true
		} else {
			fmt.Println(err)
			goto sessionInput
		}

		if dateChanged {
			err := (**appointmentTree).Remove(node)
			if err == nil {
				chAppointment := make(chan error)
				go addAppointment(appointmentTree, appointmentDate.Format("2006-01-02"), appointmentSession, dentist, patient, chAppointment)
				errAdd := <-chAppointment
				if errAdd != nil {
					fmt.Println(errAdd)
				}
			}
		} else {
			if dentistChanged {
				node.SetDentist(dentist)
			}
			if sessionChanged {
				node.SetSession(appointmentSession)
			}
		}

		fmt.Println("\nAppointment editied.")
		fmt.Println("---------------------")

		pSchedule = (**appointmentTree).GetUpComingSchedule(patient)
		printPatientSchedule(sessionList, patient, pSchedule)
	}
}

func searchAppointment(sessionList, dentistList, patientList **dll.DoublyLinkedlist, appointmentTree **bst.Binarysearchtree) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			//fmt.Println("Error searching appointment")
		}
	}()

	var (
		appointmentDate time.Time
		dentist         *Dentist
		patients        []interface{}
		patient         *Patient
		err             error
	)

	tSize := (**appointmentTree).GetSize()
	if tSize == 0 {
		return errors.New(noAppointment)
	}

	fmt.Println("\nSearch appointment")
	fmt.Println("---------------------")

apptInput:
	fmt.Println("\nPlease enter appointment date (yyyy-mm-dd) to filter. Enter to show all results:")
	inputDate := util.ReadInput()
	if len(inputDate) > 0 {
		appointmentDate, err = time.Parse("2006-01-02", inputDate)
		if err != nil {
			fmt.Println(invalidDate)
			goto apptInput
		}
	}

	fmt.Println("\nPlease enter patient name. Press enter to view all result:")
	patientName := util.ReadInput()
	if len(patientName) > 0 {
		patients = (**patientList).SearchByName(patientName)
		printPatientsSearchResult(patients, patientName)
	}
patientInput:
	if len(patients) > 0 {
		fmt.Println("\nPlease select a patient:")
		patientSelection, err := util.ReadInputAsInt()
		if err == nil {
			if patientSelection > len(patients) || patientSelection <= 0 {
				fmt.Println(invalidSelection)
				goto patientInput
			}
			patient = patients[patientSelection-1].(*Patient)
			fmt.Printf("Selected [%s]\n", patient.GetName())
		} else {
			fmt.Println(err)
			fmt.Println(invalidSelection)
			goto patientInput
		}
	}

	dentists := (**dentistList).GetList()
	printDentistList(dentists)
dentistInput:
	fmt.Println("\nPlease select a dentist to filter. Enter 0 to show all results:")
	dentistSelection, err := util.ReadInputAsInt()
	if err == nil {
		if dentistSelection > len(dentists) || dentistSelection < 0 {
			fmt.Println(invalidSelection)
			goto dentistInput
		}
		if dentistSelection > 0 {
			dentist = dentists[dentistSelection-1].(*Dentist)
			fmt.Printf("Selected dentist [%s]\n", dentist.GetName())
		}
	} else {
		fmt.Println(err)
		goto dentistInput
	}

	printSimpleSession(sessionList)
sessionInput:
	fmt.Println("\nPlease enter session number to filter [1-7]. Enter 0 to show all results:")
	appointmentSession, err := util.ReadInputAsInt()
	if err == nil {
		if appointmentSession < 0 || appointmentSession > 7 {
			fmt.Println("Invalid session, please select another session.")
			goto sessionInput
		}
	} else {
		fmt.Println(err)
		goto sessionInput
	}

	chSearchDate := make(chan []*bst.BinaryNode)
	chSearchPatient := make(chan []*bst.BinaryNode)
	chSearchDentist := make(chan []*bst.BinaryNode)
	chSearchSession := make(chan []*bst.BinaryNode)

	count := 0

	if len(inputDate) > 0 {
		count++
		go (**appointmentTree).SearchAllNodeByField("date", appointmentDate.Format("2006-01-02"), chSearchDate)
	}
	if len(patientName) > 0 {
		count++
		go (**appointmentTree).SearchAllNodeByField("patient", patient, chSearchPatient)
	}
	if dentistSelection > 0 {
		count++
		go (**appointmentTree).SearchAllNodeByField("dentist", dentist, chSearchDentist)
	}
	if appointmentSession > 0 {
		count++
		go (**appointmentTree).SearchAllNodeByField("session", appointmentSession, chSearchSession)
	}

	var result []*bst.BinaryNode
	for i := 0; i < count; i++ {
		select {
		case ret := <-chSearchDate:
			result = append(result, ret...)
		case ret2 := <-chSearchPatient:
			result = append(result, ret2...)
		case ret3 := <-chSearchDentist:
			result = append(result, ret3...)
		case ret4 := <-chSearchSession:
			result = append(result, ret4...)
		}
	}

	result = getDup(result, count)

	printAllAppointment(sessionList, result, appointmentDate)

	return nil
}

func getDup(list []*bst.BinaryNode, count int) []*bst.BinaryNode {

	duplicate_frequency := make(map[*bst.BinaryNode]int)
	var temp []*bst.BinaryNode

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := duplicate_frequency[item]

		if exist {
			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			duplicate_frequency[item] = 1 // else start counting from 1
		}
	}

	for v, n := range duplicate_frequency {
		if n == count {
			temp = append(temp, v)
		}
	}
	return temp
}

func printSimpleSession(sessionList **dll.DoublyLinkedlist) {
	fmt.Println("\nSessions")
	fmt.Println("-------------------------------------")
	list := (**sessionList).GetList()
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 5, '\t', 0)
	fmt.Fprintf(writer, "S/N\tStart Time\tEnd Time\n")
	for idx, data := range list {
		session := data.(Session)
		fmt.Fprintf(writer, "[%v]\t%s\t%s\n", idx+1, session.GetStartTime(), session.GetEndTime())
	}
	writer.Flush()
}

func printAllAppointment(sessionList **dll.DoublyLinkedlist, list []*bst.BinaryNode, appointmentDate time.Time) {

	fmt.Println("\nListing all appointments")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------------")

	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 5, '\t', 0)
	fmt.Fprintf(writer, "S/N\tDentist\tPatient\tMobile Number\tSession\tDate\tStart Time\tEnd Time\n")
	for idx, data := range list {
		patient := data.GetPatient().(*Patient)
		dentist := data.GetDentist().(*Dentist)
		r, _ := (**sessionList).Get(data.GetSession())
		session := r.(Session)
		fmt.Fprintf(writer, "%v\t%s\t%s\t%d\t%d\t%s\t%s\t%s\n", idx+1, dentist.GetName(), patient.GetName(), patient.GetMobileNum(), data.GetSession(), data.GetDate(), session.GetStartTime(), session.GetEndTime())
	}
	writer.Flush()
}
