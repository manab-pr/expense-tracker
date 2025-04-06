# expanse-tracker 


# Expense Tracker

Expense Tracker is a robust and scalable application built to help users manage their budgets and track expenses efficiently. It is developed using **Go** with the **Gin-Gonic** framework for routing, **MongoDB** as the primary database, and integrates with **Redis** for caching, **AWS S3** for file storage, and **Nginx** as a reverse proxy. The application is containerized using **Docker** and orchestrated with **Kubernetes** (Minikube) for easy deployment and scaling.

---

## Features

- **Budget Management**: Create, update, and manage budgets for different categories.
- **Expense Tracking**: Track expenses for each category and calculate remaining budgets.
- **Monthly Reports**: Generate monthly reports for income and expenses.
- **Caching**: Use Redis to cache frequently accessed data for improved performance.
- **File Storage**: Store and retrieve files (e.g., receipts) using AWS S3.
- **Containerized**: Dockerized for easy deployment and scalability.
- **Kubernetes Support**: Deploy the application on Kubernetes using Minikube.
- **Reverse Proxy**: Nginx is used as a reverse proxy for load balancing and SSL termination.

---

## Technologies Used

- **Backend**: Go (Gin-Gonic)
- **Database**: MongoDB
- **Caching**: Redis
- **File Storage**: AWS S3
- **Containerization**: Docker
- **Orchestration**: Kubernetes (Minikube)
- **Reverse Proxy**: Nginx

---

## Prerequisites

Before running the project, ensure you have the following installed:

1. **Go** (v1.20 or higher)
2. **Docker** (v20.10 or higher)
3. **Kubernetes** (Minikube v1.30 or higher)
4. **MongoDB** (v6.0 or higher)
5. **Redis** (v7.0 or higher)
6. **AWS CLI** (for S3 integration)
7. **Nginx** (for reverse proxy)

---

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/expense-tracker.git
cd expense-tracker
```

### 2. Set Up Environment Variables

Create a `.env` file in the root directory and add the following variables:

```env
MONGODB_URI=mongodb://localhost:27017
REDIS_URL=localhost:6379
AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key
S3_BUCKET_NAME=your-s3-bucket-name
```

### 3. Build and Run with Docker

```bash
# Build the Docker image
docker build -t expense-tracker .

# Run the Docker container
docker run -d -p 8080:8080 --env-file .env expense-tracker
```

### 4. Deploy on Kubernetes (Minikube)

```bash
# Start Minikube
minikube start

# Apply Kubernetes configurations
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Access the application
minikube service expense-tracker-service
```

### 5. Set Up Nginx

Configure Nginx as a reverse proxy by adding the following to your Nginx configuration:

```nginx
server {
    listen 80;
    server_name expense-tracker.local;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

Restart Nginx to apply the changes:

```bash
sudo systemctl restart nginx
```

---

## API Endpoints

### Budget Management
- **GET /budget**: Fetch budget for a specific category and month.
- **POST /budget**: Create a new budget.
- **PUT /budget**: Update an existing budget.
- **DELETE /budget**: Delete a budget.

### Expense Tracking
- **GET /expenses**: Fetch expenses for a specific category and month.
- **POST /expenses**: Add a new expense.
- **PUT /expenses**: Update an existing expense.
- **DELETE /expenses**: Delete an expense.

### Reports
- **GET /reports/monthly**: Generate a monthly report for income and expenses.

---

## Project Structure

```
expense-tracker/
├── api/                  # API handlers
├── config/               # Configuration files
├── controllers/          # Business logic
├── models/               # Data models
├── repositories/         # Database operations
├── services/             # Service layer
├── k8s/                  # Kubernetes configurations
├── Dockerfile            # Docker configuration
├── go.mod                # Go dependencies
├── go.sum                # Go dependencies checksum
├── README.md             # Project documentation
└── .env                  # Environment variables
```

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Acknowledgments

- **Gin-Gonic** for the lightweight and fast HTTP framework.
- **MongoDB** for the flexible NoSQL database.
- **Redis** for efficient caching.
- **AWS S3** for reliable file storage.
- **Kubernetes** and **Docker** for container orchestration and deployment.
- **Nginx** for reverse proxy and load balancing.

---

## Contact

For any questions or feedback, feel free to reach out:

- **Email**: expensetracker@orion.com
- **GitHub**: https://github.com/manab-pr/expanse-tracker/tree/main

---

Happy Budgeting! 🚀
