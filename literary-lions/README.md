# 🦁 Literary Lions Forum 

## 👀 Overview 

Welcome to the Literary Lions Forum, a web-based book club platform designed to enhance and streamline book discussions. Our platform transforms the chaotic paper-based discussions into a structured digital forum, enabling book enthusiasts to engage in vibrant discussions, categorize posts, like/dislike comments, and filter posts effectively. This project is built using Go for the backend, SQLite for data management, and Docker for seamless deployment.

## 👇 Prerequisites

Before you get started, ensure you have the following installed:
- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)


## Setup and Installation

1. **Clone the Repository**:
    ```bash
    git clone https://gitea.kood.tech/katriinsartakov/literary-lions.git
    cd literary-lions
    ```

2. **Install Dependencies**:
Ensure all Go dependencies are installed.
    ```bash
    go mod tidy
    ```

## Docker Instructions

1. **Build the Docker Image**:
    ```bash
    docker build -t literary-lions-app .
    ```

2. **Run Docker Container**:
    ```bash
    docker run -p 8080:8080 literary-lions-app
    ```

3. **Access the Application**:
Open your web browser and navigate to [http://localhost:8080](http://localhost:8081).

## Database Setup
The SQLite database will be automatically initialized with the necessary tables using the SQL script located in `dbscripts/0001_create_users_table.sql`.


## 🫶 Project Features 


- **User Authentication**: Register and login functionalities with encrypted passwords.


- **Post Filtering**: Filter posts by category, created posts, and liked posts.


- **Profile Management**: User profiles with personal information and posts.


- **Communication Features**: Create posts and comments, view posts and comments without registration.

- **Like/Dislike Functionality**

- **Profile Pictures**


## 🤝 Support

If you encounter any errors or issues, please contact us via Discord. We will be happy to help!
