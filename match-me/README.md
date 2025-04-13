# 👀 Match-Me Web

**Match-Me Web** is a full-stack recommendation platform that connects users based on their interests, preferences, and profiles. Whether you're seeking friendship, professional networking, or a hobby partner, Match-Me Web helps you find the perfect match.

## 🚀 Features

- **User Authentication** – Secure login, JWT-based sessions.
- **User Profiles** – Customizable bios, profile pictures.
- **Recommendations** – Smart user recommendations.
- **Connections** – Send, accept, and manage connections.
- **Real-Time Chat** – WebSocket-based messaging.
- **GraphQL API** – Optimized queries for mobile efficiency.

---

## 🛠️ Tech Stack

- **Backend:** Go, PostgreSQL
- **Frontend:** React
- **API:** GraphQL + REST
- **Authentication:** JWT
- **Real-Time Communication:** WebSockets

---

## 🛠️ Setup Instructions

### **1️⃣ Clone the Repository**
```bash
git clone https://gitea.kood.tech/katriinsartakov/match-me.git
cd match-me
```

### **2️⃣ Set Up PostgreSQL**
```bash
sudo systemctl start postgresql
psql -U postgres -h localhost
```
```sql
CREATE DATABASE match_me;
ALTER USER postgres PASSWORD 'yourpassword';
\q
```

### **3️⃣ Configure Environment Variables**
Create a `.env` file in the root of the project:
```env
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=match_me
DB_HOST=localhost
DB_PORT=5432
DB_SSLMODE=disable
```

### **4️⃣ Install Dependencies**
#### Backend:
```bash
go mod tidy
```
#### Frontend:
```bash
cd frontend
npm install
```

---

## 🚀 Running the Application

### **Run Backend**
```bash
go run main.go
```
- Runs on **http://localhost:8080**
- **GraphQL Playground** is only available in **developer mode** (see below).

### **Run Frontend**
```bash
cd frontend
npm start
```
- Runs on **http://localhost:3000**

---

## 🎯 Using the GraphQL API

### **Enable Developer Mode**
To use **GraphQL Playground**, start the backend with:
```bash
go run main.go -d
```
Then visit:  
👉 **http://localhost:8080/graphql**

### **Test Queries**
You can test the GraphQL API via **Postman**, **GraphQL Playground**, or a **GraphQL client**.

### **Authorization for User-Specific Queries**
For queries that fetch specific user data (such as `me`, `myBio`, or `connections`), you must include an **Authorization** header with a valid JWT token. Example:
```json
{
  "Authorization": "Bearer your_jwt_token_here"
}
```
Ensure you're logged in and using a valid token to access protected endpoints.

#### **1️⃣ Get Current Logged-in User**
```graphql
{
  me {
    id
    name
    email
    profilePicture
    bio {
      interests
      age
      location
      info
    }
  }
}
```

#### **2️⃣ Get a Specific User by ID**
```graphql
{
  user(id: "104") {
    id
    name
    email
    profilePicture
    bio {
      interests
      age
      gender
      location
      info
    }
  }
}
```

#### **3️⃣ Get Recommended Users**
```graphql
{
  recommendations {
    id
    name
    profilePicture
    bio {
      interests
      age
    }
  }
}
```

#### **4️⃣ Get Accepted Connections**
```graphql
{
  connections {
    id
    name
    email
    bio {
      interests
      age
    }
  }
}
```

#### **5️⃣ Get a Specific Bio by ID**
```graphql
{
  bio(id: "104") {
    id
    interests
    age
    location
    info
    user {
      id
      name
    }
  }
}
```

### **Check Available Queries**
Run this to inspect available queries:
```graphql
{
  __schema {
    queryType {
      fields {
        name
      }
    }
  }
}
```

---

## 📆 Generating 100 Random Users
Run this before starting the server:
```bash
go run scripts/main.go scripts/seed_users.go
```

---

