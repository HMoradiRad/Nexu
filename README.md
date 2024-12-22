# Nexu challange
### Develop and deploy a basic Golang application with PostgreSQL, set up Prometheus metrics,and deploy the system on Kubernetes(minikube).


This project demonstrates the development and deployment of a Golang RESTful API application with PostgreSQL, NGINX reverse proxy, and Prometheus monitoring on Kubernetes.

### Features:
- A Golang RESTful API to manage meeting data.
- PostgreSQL database deployment in Kubernetes.
- Prometheus metrics exposed for API performance monitoring.
- NGINX as a reverse proxy, serving static content for Google bots, routing /blog to an external IP, and routing all other traffic to the Golang application.

### Prerequisites
To run this project, you need:

- Minikube or any Kubernetes cluster setup
- kubectl configured to interact with your Kubernetes cluster
- Docker to build and push images
- Go 1.22.2 or later for building the application
- PostgreSQL database configured in Kubernetes
- Prometheus set up for monitoring


### Project Setup

1) **Clone the Repository**
```
https://github.com/HMoradiRad/golang-meeting-app.git

cd golang-meeting-app
```

### Create the Go application.

- The application is built as a RESTful API that handles meetings:
-  POST /meetings: Create or update a meeting.
-  GET /meetings: List all meetings.
-  GET /metrics: Expose Prometheus metrics.



2) **Build And Push Docker Image**

```
docker build -t golang-meeting-app:1.0 .
```
**Push to a Docker registry (if required):**

```
docker push <your-docker-registry>/golang-meeting-app
```

### Deploy to Kubernetes
1. Set Up PostgreSQL
**Apply the PostgreSQL StatefulSet:**

```
kubectl apply -f postgres-statefulset.yaml
```

2) Deploy the Golang Application
**Deploy the application:**

```
kubectl apply -f golang-app-deployment.yaml
```

3) Deploy NGINX as Reverse Proxy
```
kubectl apply -f nginx-deployment.yaml
```

4) Deploy Prometheus
```
kubectl apply -f prometheus-deployment.yaml
```
**Now the application should be accessible at:**
```
http://<minikube-ip>:<port>
```

#  Test the API Endpoints

**Add a New Meeting**

```
curl -X POST http://<minikube-ip>:<port-nginx>/meetings \
-H "Content-Type: application/json" \
-d '{
  "title": "nexu",
  "invitee_email": "nexu@gmail.com",
  "host_email": "nexu@gmail.com",
  "start_at": "2024-12-22T09:00:00Z",
  "end_at": "2024-12-22T10:00:00Z",
  "status": "scheduled"
}'

```
**Update a Meeting**
```
curl -X POST http://<minikube-ip>:<port-nginx>/meetings  \
-H "Content-Type: application/json" \
-d '{
  "id": 1,
  "status": "completed"
}'
```
