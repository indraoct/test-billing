# Loan Billing Engine
For some purpose

## How to Use this Backend Apps

    1. Make sure you already install docker in your local computer
    2. To run this just run this script : 
```cmd
    docker-compose up
```
    It will be run the go apps and dummy docker database (just for test). The Apps will run in port 8080 
    and the database on port 5432 (postgres)
    3. Before you run the workflow of system, you should read the "System Design for Loan Billing Engine" to know the big picture of this app.
    4. There is postman collection on this application you can test it using postman.
    5. For dummy data and authentication i just put 2 model auth (dummy) there are : admin-token and user-token. In this application i just prepare 1 customer : customer_id : John Doe (because i dont have may times to do this test)
    6. The notification process is still dummy (in the real case should using kafka/rabbitmq).
    7. The application maybe not perfect but I already do the best that I could due to limited time and my schedule on my current weekday working time is tight. Hope you can understand well about my work. Thanks :)

## System Design for Loan Billing Engine

### 1. High-Level Architecture

The Loan Billing Engine is a service that handles loan management, including billing schedules, payments, outstanding balances, and delinquency tracking. The system can be split into several components:

 #### a. API Layer
    Handles incoming requests from clients and routes them to appropriate services.

 #### b. Service Layer
    Encapsulates the business logic for loan management, such as:
    - Generating loan schedules.
    - Tracking payments and calculating outstanding balances.
    - Checking delinquency status.

 #### c. Data Layer
    Manages data persistence in a relational database like PostgreSQL or MySQL.

 #### d. External Services

    Notification Service: Notifies borrowers of upcoming or missed payments.
    Analytics Service: Tracks loan performance metrics for internal use.

### 2. Components

 #### a. Loan Management Service
    Core service responsible for:
    - Creating loans.
    - Generating repayment schedules.
    - Handling payments.
    - Delinquency checks.

 #### b. Notification Service
    Stores all loan-related data:
    - Loan details (ID, principal, interest rate, etc.).
    - Repayment schedules and payment statuses.
    - Customer details.

 #### c. Database
    Stores all loan-related data:
    - Loan details (ID, principal, interest rate, etc.).
    - Repayment schedules and payment statuses.
    - Customer details.

 #### d. Frontend (Optional)
    Web or mobile UI for borrowers and admins to interact with the system.

### 3. Database Design
Here’s a relational schema for the database:

 #### Tables:

    Customers

    - id (Primary Key)
    - name
    - email
    - phone
    
    Loans
    
    - id (Primary Key)
    - customer_id (Foreign Key -> Customers)
    - principal_amount
    - interest_rate
    - weeks
    - weekly_payment
    - outstanding_balance
    - status (e.g., Active, Paid, Delinquent)
    
    RepaymentSchedules

    - id (Primary Key)
    - loan_id (Foreign Key -> Loans)
    - week_number
    - due_date
    - status (e.g., Paid, Unpaid)


    Payments
    
    - id (Primary Key)
    - loan_id (Foreign Key -> Loans)
    - week_number
    - amount
    - payment_date

### 4. API Endpoints
Here are the key endpoints:


|Endpoint|	Method|	Description
|:---:| :---: | :---: |
|/loans|	POST|	Create a new loan.
|/loans/{id} |	GET	| Get loan details, including schedule.
|/loans/{id}/outstanding | 	GET |	Get outstanding balance.
|/loans/{id}/delinquent	| GET |	Check if the borrower is delinquent.
|/loans/{id}/payments |	POST |	Make a payment for a loan.
| /customers/{id}	| GET	| Get customer details.


### 5. System Workflow
 #### a. Loan Creation

    - Admin calls /loans API with loan details.
    - Loan schedule is generated and stored in the RepaymentSchedules table.
    - Initial status is set to Active.

 #### b. Payment Processing

    - Borrower makes a payment via /loans/{id}/payments.
    - Payment is recorded in the Payments table.
    - Outstanding balance is recalculated.
    - Delinquency status is updated if applicable.

 #### c. Delinquency Tracking

    - A background job checks for loans with missed payments and marks them as delinquent if necessary.

 #### d. Notification

    - Notification service sends reminders for upcoming payments and delinquency warnings.

#### 6. Technology Stack (for this test i just code in backend)

 #### a. Backend

    - Language: Go (for performance and concurrency).
    - Framework: Echo for REST APIs.

 #### b. Database

    PostgreSQL: For relational data with transactional guarantees.

 #### c. Queueing System

    RabbitMQ/Kafka: For asynchronous communication between services (e.g., sending notifications).

 #### d. Frontend

    - React or Vue.js for web.
    - Flutter for mobile apps.

 #### e. Deployment

    - Docker & Kubernetes: For containerization and orchestration.
    - CI/CD: GitHub Actions, Jenkins, or GitLab CI for automated deployments.

 #### f. Monitoring

    - Prometheus & Grafana: For metrics and alerting.
    - ELK Stack: For centralized logging.

### 7. High Level Diagram
               +-------------------+       +-------------------+
               |   Notification    |       |    Analytics      |
               |      Service      |       |      Service      |
               +-------------------+       +-------------------+
                         ^                         ^
                         |                         |
                         |                         |
    +----------+   +------------+          +-------------------+
    |  Client  |<->| API Layer  |<-------->| Loan Management   |
    +----------+   +------------+          |    Service        |
                        |                         |
                        |                         |
                    +-------------+         +-------------------+
                    |  Database   |<------->| Repayment Tracker |
                    +-------------+         +-------------------+
