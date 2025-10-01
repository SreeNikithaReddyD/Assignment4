# CS6650 - Simple Product API Service (Go & Terraform)

This repository contains a simple Go web server that implements a Product API, containerized with Docker. The entire cloud infrastructure required to deploy this service on AWS Elastic Container Service (ECS) is managed declaratively using Terraform.

## Repository Structure

-   **/src**: Contains the Go server code (`main.go`), its dependencies (`go.mod`, `go.sum`), and the `Dockerfile` for containerization.
-   **/terraform**: Contains all Terraform configuration files (`.tf`) for managing the AWS infrastructure, including modules for network, ECR, ECS, and logging.
-   **/testing**: Contains the Locust load testing script (`locustfile.py`).
-   **.gitignore**: Configured to exclude sensitive files like `.tfvars` (if any present)and Terraform state files.
-   I didnt create any .tfvars file as I had default values set and did not really need any .tfvars file.

---

## Part I: Deployment Instructions

Follow these steps to deploy the infrastructure and service to your AWS account.

### Prerequisites

1.  An AWS Account with programmatic access credentials.
2.  Terraform installed locally.
3.  Docker installed and running locally.
4.  AWS CLI configured on your machine (`aws configure`) and aws configure set aws_session_token {temporary_session_token}.

### Steps

1.  **Clone the Repository**

    git clone <your-repo-url>
    cd <your-repo-name>/terraform
 

2.  **Initialize and Apply Terraform**
    These commands will set up all the necessary AWS resources.

    # Initialize Terraform providers and modules
    terraform init

    # Apply the configuration to create the infrastructure
    terraform apply -auto-approve
 
    This process automatically builds the Docker image from the `/src` directory, creates an ECR repository, pushes the image, and deploys it as a service on ECS.

3.  **Get the Service Public IP Address**
    Once the deployment is complete, run the command below from the `terraform` directory to get the public IP of your running container.

    aws ec2 describe-network-interfaces --network-interface-ids $(aws ecs describe-tasks --cluster $(terraform output -raw ecs_cluster_name) --tasks $(aws ecs list-tasks --cluster $(terraform output -raw ecs_cluster_name) --service-name $(terraform output -raw ecs_service_name) --query 'taskArns[0]' --output text) --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" --output text) --query 'NetworkInterfaces[0].Association.PublicIp' --output text


---

## Part II: API Endpoints and Examples

All API requests should be made to `http://<YOUR_PUBLIC_IP>:8080`.

### 1. Create a Product

-   **Endpoint**: `POST /products`
-   **Description**: Adds a new product to the store.
-   **Request Body**: JSON object with `name`, `price`, and `quantity`.

#### Example Responses:

-   **201 Created (Success)**
    curl -X POST http://<YOUR_PUBLIC_IP>:8080/products \
    -H "Content-Type: application/json" \
    -d '{"name": "Gaming Mouse", "price": 75.99, "quantity": 200}'

    # Response Body:
    # {"id":"1","name":"Gaming Mouse","price":75.99,"quantity":200}

-   **400 Bad Request (Invalid Input)**
    *This happens if required fields are missing or if `price` is not positive.*

    curl -X POST http://<YOUR_PUBLIC_IP>:8080/products \
    -H "Content-Type: application/json" \
    -d '{"name": "Defective Keyboard", "price": -50, "quantity": 10}'

    # Response Body:
    # {"error":"Price must be positive, and Quantity cannot be negative."}
    

### 2. Get a Product by ID

-   **Endpoint**: `GET /products/{id}`
-   **Description**: Retrieves a single product by its ID.

#### Example Responses:

-   **200 OK (Success)**
    curl http://<YOUR_PUBLIC_IP>:8080/products/1

    # Response Body:
    # {"id":"1","name":"Gaming Mouse","price":75.99,"quantity":200}
 

-   **404 Not Found (Product does not exist)**
    curl http://<YOUR_PUBLIC_IP>:8080/products/999

    # Response Body:
    # {"error":"Product not found"}

---

## Part III: Clean Up

To tear down all the infrastructure created by Terraform, run the following command from the `/terraform` directory:
```bash
terraform destroy -auto-approve