package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var messages []openai.ChatCompletionMessageParamUnion

// Doctor struct with available slots mapped by date
type Doctor struct {
	FirstName      string              `json:"first_name"`
	LastName       string              `json:"last_name"`
	Specialty      string              `json:"specialty"`
	AvailableSlots map[string][]string `json:"available_slots"` // Key: Date, Value: List of slots
}

func main() {
	// Initialize OpenAI client
	client := openai.NewClient(option.WithAPIKey("api-key"))

	// Terminal input reader
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Super Clinic Chatbot!")
	fmt.Println("Type 'exit' to quit.")

	systemRole := `You are a virtual assistant for Super Clinic, specializing in scheduling doctor appointments. You have access to the clinic's doctor roster and their available time slots. Your primary objectives include:

    ### Appointment Scheduling:
    - Assist patients in booking appointments by cross-referencing their requested times with doctors' availability.
    - Ensure no double bookings occur for the same doctor during the same time slot.

    ### Alternative Suggestions:
    - If the requested time is unavailable, suggest the nearest available slots with the same doctor.
    - If the requested doctor has no availability, recommend alternative doctors with similar expertise.

    ### Slot Management:
    - Keep track of booked appointments, patient names, and remaining slots for each doctor.
    - Respond to specific requests like "give slot info" by listing both booked and available slots, along with patient names for booked slots.

    ### Function Calling for Checking Doctor Availability:
    - When a user asks about a doctor’s availability for a specific time (e.g., "Is Dr. Jamil Ahmad available at 10:00 AM?"), use the function **is_doctor_available**.
    - The function **is_doctor_available** will check the database and return true if the doctor is available or false if not.
    - If the requested slot is unavailable, suggest alternative available slots.

    ### Confirmation Summary:
    - Upon successful booking, provide a summary that includes:
        - Doctor's name
        - Appointment date and time
        - Instructions to arrive 10 minutes early and bring any necessary medical documents.

    ### Unavailable Slots Notification:
    - If no slots are available for the requested time, notify the patient politely and suggest alternatives, either with the same doctor or with other doctors.

    ### Operating Hours:
    - The clinic operates Monday through Friday.
    - Doctors are available during the following sessions:
        - Morning Session: 9:00 AM to 1:00 PM
        - Evening Session: 5:00 PM to 9:00 PM
    - Each appointment lasts between 30 and 50 minutes, depending on the doctor.
    - The clinic is closed on weekends.

    ### Sample Appointment Confirmation:
    "Your appointment has been successfully scheduled with Dr. [First Name] [Last Name] on [Date] at [Time]. Please arrive 10 minutes early and bring any necessary medical documents."

    ### Sample Alternative Time Suggestion:
    "I'm sorry, Dr. [Last Name] is not available at your requested time. However, the closest available slots are:

    - [Date] at [Alternative Time 1]
    - [Date] at [Alternative Time 2]"
    Would you like to book one of these slots, or would you prefer to see another doctor with similar expertise?"`

	assistantRole := "You are a helpful assistant for a doctor's clinic, assisting with appointment scheduling."
	messages = append(messages, openai.SystemMessage(systemRole))
	messages = append(messages, openai.AssistantMessage(assistantRole))

	for {
		// Get user input
		fmt.Print("\nYou: ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}

		// Exit condition
		if userInput == "exit\n" {
			fmt.Println("Goodbye!")
			break
		}

		// Append user message to history
		messages = append(messages, openai.UserMessage(userInput))

		response := simpleRequestChat(client)

		// Append assistant response to history
		messages = append(messages, openai.AssistantMessage(response))
	}
}

// Function calling tool definition
func getDoctorAvailabilityTool() openai.ChatCompletionToolParam {
	return openai.ChatCompletionToolParam{
		Type: openai.F(openai.ChatCompletionToolTypeFunction),
		Function: openai.F(openai.FunctionDefinitionParam{
			Name:        openai.String("is_doctor_available"),
			Description: openai.String("Check if a doctor is available on a given date and time."),
			Parameters: openai.F(openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"doctor_name": map[string]interface{}{
						"type":        "string",
						"description": "Full name of the doctor to check availability for. Example: 'Jamil Ahmad'",
					},
					"requested_date": map[string]interface{}{
						"type":        "string",
						"description": "Date of the requested appointment in YYYY-MM-DD format. Example: '2025-01-21'",
					},
					"requested_time": map[string]interface{}{
						"type":        "string",
						"description": "Time slot for the requested appointment in HH:MM format (24-hour format). Example: '09:00'",
					},
				},
				"required": []string{"doctor_name", "requested_date", "requested_time"},
			}),
		}),
	}
}

// OpenAI Chat Request Function
func simpleRequestChat(client *openai.Client) (response string) {
	//ex : I want to book appointment for  Jamil Ahmad at 10 AM on tommorow

	// Define function tool
	doctorAvailabilityTool := openai.F([]openai.ChatCompletionToolParam{getDoctorAvailabilityTool()})

	params := openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Tools:    doctorAvailabilityTool,
		Seed:     openai.Int(0),
		Model:    openai.F(openai.ChatModelGPT4o),
	}

	ctx := context.Background()
	chatCompletion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err.Error())
	}

	// Handle function calls from OpenAI
	toolCalls := chatCompletion.Choices[0].Message.ToolCalls
	if len(toolCalls) == 0 {
		fmt.Println("No function call triggered.")
	} else {
		// Process tool calls
		params.Messages.Value = append(params.Messages.Value, chatCompletion.Choices[0].Message)
		for _, toolCall := range toolCalls {
			if toolCall.Function.Name == "is_doctor_available" {
				var args map[string]interface{}
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
					panic(err)
				}
				doctorName := args["doctor_name"].(string)
				requestedDate := args["requested_date"].(string)
				requestedTime := args["requested_time"].(string)

				// Check doctor availability
				available := isDoctorAvailable(doctorName, requestedDate, requestedTime)

				// Send function result back to OpenAI
				params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, fmt.Sprintf("%v", available)))
			}
		}
	}

	chatCompletion, err = client.Chat.Completions.New(ctx, params)
	if err != nil {
		log.Fatalf("Error creating chat completion: %v", err)
	}

	//Print chatbot response
	response = chatCompletion.Choices[0].Message.Content
	fmt.Printf("Chatbot: %s\n", response)
	return response
}

// check if doctor is available
func isDoctorAvailable(doctorName string, requestedDate string, requestedTime string) bool {
	fmt.Printf("Checking availability for Doctor: %s, Date: %s, Time: %s\n", doctorName, requestedDate, requestedTime)
	time.Sleep(1 * time.Second)
	doctors, err := parseDoctorSchedule("doctor_schedule.csv")
	if err != nil {
		log.Fatalf("Failed to load doctor schedule: %v", err)
	}

	// Check if the doctor exists
	doctor, exists := doctors[doctorName]
	if !exists {
		fmt.Printf("Doctor %s not found.\n", doctorName)
		return false
	}

	// Check if the requested date exists
	availableSlots, dateExists := doctor.AvailableSlots[requestedDate]
	if !dateExists {
		fmt.Printf("Doctor %s is not available on %s.\n", doctorName, requestedDate)
		return false
	}

	// Check if the requested time exists in available slots
	for _, slot := range availableSlots {
		if strings.Contains(slot, requestedTime) {
			return true
		}
	}
	fmt.Printf("Doctor %s is not available at %s on %s.\n", doctorName, requestedTime, requestedDate)
	return false
}

// parseDoctorSchedule reads the doctor_schedule.csv file and returns a map of Doctor structs
func parseDoctorSchedule(filename string) (map[string]Doctor, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %v", err)
	}

	doctors := make(map[string]Doctor)

	// Skip the first row (header)
	for i, row := range records {
		if i == 0 {
			continue
		}

		firstName := row[0]
		lastName := row[1]
		specialty := row[2]
		date := row[3]
		availableSlots := strings.Split(row[4], ",")

		// Create doctor key
		fullName := firstName + " " + lastName

		// If doctor exists, append slots to the existing date key
		if doc, exists := doctors[fullName]; exists {
			doc.AvailableSlots[date] = availableSlots
			doctors[fullName] = doc
		} else {
			doctors[fullName] = Doctor{
				FirstName:      firstName,
				LastName:       lastName,
				Specialty:      specialty,
				AvailableSlots: map[string][]string{date: availableSlots},
			}
		}
	}

	return doctors, nil
}
