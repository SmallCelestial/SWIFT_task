# SWIFT_task

## 🛠 Tech Stack

### **Programming Language & Frameworks**
- 🐹 **Go (Golang)** – Main backend programming language
- **Gin** – A fast web framework for building REST APIs
- **Gorm** – An ORM for database management

### **Database**
- **PostgreSQL** – Relational database management system

### **Other Technologies**

- **Docker** - Containerization platform used to build and run services in isolated environments.
- **Docker Compose** - Tool for defining and running multi-container Docker applications.

---

## 📋 Requirements

To run this application locally, you need to have the following installed:

- **Docker**: [Install Docker](https://www.docker.com/get-started)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)

---

## ⚙️ Setup Instructions

### 1. Clone the Repository

Start by cloning the repository to your local machine:

```bash
git clone https://github.com/SmallCelestial/SWIFT_task.git
cd SWIFT_task
```

### 2. Build and Start the Application

Ensure Docker and Docker Compose are installed and  **running**. Then, navigate to the project directory and run the following command to
build and start all services:

```bash
docker-compose up --build
```

If the application is already built, and you just want to start the containers, use this command:

```bash
docker-compose up
```

### 3.Access the Application

Once the containers are up and running, you can access the application locally:

`http://localhost:8080/`

---

## 🧪 Testing
You can easily test application using prepared postman collection.

[Click here to view the Postman documentation](https://documenter.getpostman.com/view/37233656/2sAYX3sPvo)