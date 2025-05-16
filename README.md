# 🏥 AI-Super-Clinic Chatbot

An AI-powered Doctor's Appointment Chat Application built with Golang and OpenAI's Function Calling capabilities.

## 📖 Overview

This application simulates a virtual assistant for a doctor's clinic, enabling users to check doctor availability and book appointments through natural language interactions. It leverages OpenAI's Chat API with Function Calling to process user queries and interact with a predefined doctor schedule.

## 🚀 Features

- Interactive chat interface for appointment scheduling
- Integration with OpenAI's Chat API and Function Calling
- Checks doctor availability against a CSV-based schedule
- Handles natural language inputs for seamless user experience

## 🛠️ Prerequisites

- Go 1.18 or higher
- OpenAI API Key

## 📦 Installation

1. **Clone the Repository**
   ```bash
   git clone https://github.com/jamil225/AI-Super-Clinic-Chatbot.git
   cd AI-Super-Clinic-Chatbot
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Set Up Environment Variables**

   Set your OpenAI API key as an environment variable:
   ```bash
   export OPENAI_API_KEY=your_openai_api_key
   ```

   *Alternatively*, you can directly replace the placeholder in `main.go` with your API key:
   ```go
   client := openai.NewClient(option.WithAPIKey("your_openai_api_key"))
   ```

## 📄 Usage

1. **Run the Application**
   ```bash
   go run main.go
   ```

2. **Interact with the Chatbot**

   After running the application, you'll be prompted to enter your queries. For example:
   ```
   You: Is Dr. Smith available on 2025-05-20 at 10:00?
   ```

   The chatbot will respond based on the availability data provided in `doctor_schedule.csv`.

3. **Exit the Application**

   To terminate the chat, type:
   ```
   exit
   ```

## 📂 Project Structure

```
├── doctor_schedule.csv   // Contains the doctor's availability schedule
├── go.mod                // Go module file
├── go.sum                // Go dependencies checksum file
├── main.go               // Main application file
└── README.md             // Project documentation
```

## 🧠 How It Works

1. **Initialization**

   The application initializes the OpenAI client using the provided API key and sets up the system and assistant roles to guide the chatbot's behavior.

2. **User Interaction Loop**

   The application enters a loop where it continuously prompts the user for input, processes the input, and responds accordingly.

3. **Function Calling**

   When a user query pertains to checking a doctor's availability, the application utilizes OpenAI's Function Calling feature to invoke the `is_doctor_available` function. This function checks the `doctor_schedule.csv` file to determine availability and returns the result to the chatbot, which then communicates it to the user.

## 📊 Doctor Schedule Format

The `doctor_schedule.csv` file should be structured as follows:

```csv
doctor_name,date,time
Dr. Smith,2025-05-20,10:00
Dr. Johnson,2025-05-21,14:30
```

Each row represents a time slot when a doctor is available. The chatbot references this file to determine availability.

## 🔗 Related Resources

- [Medium Article: Building a Doctor’s Appointment Chat Application with Golang and OpenAI](https://medium.com/@jamil.ahmad7720/building-a-doctors-appointment-chat-application-with-golang-and-openai-a-step-by-step-guide-3cbb4357ea2a)
- [OpenAI Go SDK](https://github.com/openai/openai-go)

## 📬 Contact

For questions or feedback, feel free to reach out:

- **Author**: Jamil Ahmad
- **Email**: jamil.ahmad7720@gmailcom
- **LinkedIn**: [LinkedIn](https://www.linkedin.com/in/jamil-ahmad-7720/)

---

*This project is a demonstration of integrating OpenAI's Function Calling with a Golang application to create a functional and interactive chatbot for scheduling doctor appointments.*
