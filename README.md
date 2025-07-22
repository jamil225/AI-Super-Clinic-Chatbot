# 🏥 AI Super Clinic Chatbot

A conversational AI assistant for doctor appointment scheduling built with Go and OpenAI's Function Calling API.

## 📖 Overview

This chatbot helps patients book appointments by checking doctor availability in real-time. It understands natural language queries and manages appointment scheduling through an interactive chat interface.

**Key Features:**
- 🤖 Natural language conversation for appointment booking
- 📅 Real-time doctor availability checking  
- 🔍 Alternative time slot suggestions
- 👨‍⚕️ Multi-doctor, multi-specialty support
- 📊 CSV-based schedule management

## 🚀 How to Run

### Prerequisites
- Go 1.22 or higher
- OpenAI API key

### Quick Start

1. **Clone and setup:**
   ```bash
   git clone https://github.com/jamil225/AI-Super-Clinic-Chatbot.git
   cd AI-Super-Clinic-Chatbot
   go mod tidy
   ```

2. **Configure API key:**
   ```bash
   export OPENAI_API_KEY="your_openai_api_key_here"
   ```
   
   *Alternative: Replace the API key directly in `main.go` line 28*

3. **Run the application:**
   ```bash
   go run main.go
   ```

4. **Start chatting:**
   ```
   Welcome to the Super Clinic Chatbot!
   Type 'exit' to quit.
   You: Is Dr. Ahmad available tomorrow at 10 AM?
   ```

## 💬 Usage Examples

**Check availability:**
```
You: Is Dr. Jamil Ahmad available on 2025-01-21 at 10:00?
```

**Book appointment:**
```
You: I want to book an appointment with a cardiologist for tomorrow morning
```

**Get schedule info:**
```
You: give slot info
```

**Exit:**
```
You: exit
```

## 📁 Project Structure

```
├── main.go                 # Main application with AI logic
├── doctor_schedule.csv     # Doctor availability data
├── go.mod                  # Go dependencies
└── README.md              # Documentation
```

## 🏥 Doctor Schedule Format

The `doctor_schedule.csv` contains doctor availability:

```csv
First Name,Last Name,Specialty,Date,Available Slots
Jamil,Ahmad,Cardiology,2025-01-21,"09:00-09:30,10:00-10:30,11:00-11:30"
```

## 🛠️ How It Works

1. **Chat Interface** - Continuous conversation loop with user input
2. **AI Processing** - OpenAI processes queries and determines if function calling needed  
3. **Function Calling** - When appointment-related, calls `is_doctor_available` function
4. **Schedule Check** - Reads CSV file to verify doctor availability
5. **Response** - AI provides natural language response with results

## 🔗 Resources

- [Medium Article: Building a Doctor’s Appointment Chat Application with Golang and OpenAI](https://medium.com/@jamil.ahmad7720/building-a-doctors-appointment-chat-application-with-golang-and-openai-a-step-by-step-guide-3cbb4357ea2a)
- [OpenAI Go SDK](https://github.com/openai/openai-go)

## 📞 Contact

**Jamil Ahmad**  
📧 jamil.ahmad7720@gmail.com  
💼 [LinkedIn](https://www.linkedin.com/in/jamil-ahmad-7720/)

---

*Demo project showcasing OpenAI Function Calling integration with Go for intelligent appointment scheduling.*
