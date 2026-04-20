# Rate-Limited API Service  
### Backend Assignment – Source Asia

Hi Team,

Thank you for the opportunity to work on this assignment.  
This project is my implementation of the **Rate-Limited API Service** as part of the backend assignment for Source Asia.

---

##  Candidate Details

- **Name:** Keshav Raj  
- **LinkedIn:** https://www.linkedin.com/in/keshavraj18  

---

##  Project Overview

This project implements a simple backend service with:

- A rate-limited API (`POST /request`)
- A stats API (`GET /stats`)
- In-memory data storage (no database used)
- Proper handling of concurrent requests

---

##  Tech Stack

- **Language:** Go (Golang)
- **Framework:** Fiber
- **Storage:** In-memory (map + mutex)
- **Config:** Environment variables (.env)

---

##  API Endpoints

### 1. POST `/request`

Accepts:
```json
{
  "user_id": "string",
  "payload": "any data"
}
```

Behavior:
- Allows max **5 requests per user per minute**
- Returns error if limit exceeded

---

### 2. GET `/stats`

Returns:
- Total requests per user

Supports optional pagination:
```
/stats?page=1&limit=10
```

---

##  How to Run the Project

### 1. Clone the repository

```bash
git clone https://github.com/mrperfect2003/rate-limited.git
cd rate-limited
```

---

### 2. Install dependencies

```bash
go mod tidy
```

---

### 3. Create `.env` file

Create a file named `.env` in the root folder and add:

```env
PORT=5000
RATE_LIMIT_MAX_REQUESTS=5
RATE_LIMIT_WINDOW_SECONDS=60
```

---

### 4. Run the application

```bash
go run cmd/main.go
```

---

### 5. Server will start at

```
http://localhost:5000
```

---

##  How to Test APIs

### 🔹 POST `/request`

```
POST http://localhost:5000/request
```

Body:

```json
{
  "user_id": "user1",
  "payload": "test data"
}
```

- First 5 requests → success  
- 6th request → rate limit exceeded  

---

### 🔹 GET `/stats`

```
GET http://localhost:5000/stats
```

With pagination:

```
GET http://localhost:5000/stats?page=1&limit=10
```

---

##  Design Decisions

- Used **in-memory storage (map)** as per assignment requirement
- Implemented rate limiting using **timestamp-based sliding window**
- Used **sync.Mutex** to handle concurrent requests safely
- Clean separation of:
  - handler (API layer)
  - service (business logic)
  - storage (data layer)

---

##  Rate Limiting Logic

- Store timestamps of each request per user
- On each request:
  - Remove timestamps older than 60 seconds
  - Check if requests < 5
  - If yes → allow
  - Else → reject

---

##  Concurrency Handling

- Shared data protected using **mutex**
- Ensures:
  - No race conditions
  - Accurate counting under parallel requests

Tested using:

```bash
go run -race cmd/main.go
```

---

##  Sample Responses

### Success:
```json
{
  "message": "request accepted",
  "user_id": "user1"
}
```

### Rate Limit Exceeded:
```json
{
  "error": "rate limit exceeded"
}
```

---

##  Limitations

- Data is stored in memory (lost on restart)
- Not distributed (single instance)
- No persistent storage

---

##  Improvements (If given more time)

- Use Redis for distributed rate limiting
- Add retry mechanism / queue
- Add logging & monitoring
- Deploy to cloud (AWS / Azure)
- Add authentication (JWT)

---

##  Conclusion

This implementation focuses on:
- Correct rate limiting logic
- Concurrency safety
- Clean and simple backend design

Thank you again for the opportunity.  
Looking forward to your feedback!

---