
## Setup Instructions
1. **Clone the repository**:
   ```
   git clone <repository-url>
   cd delivery-management-system
   ```

2. **Install dependencies**:
   Ensure you have Go installed, then run:
   ```
   go mod tidy
   ```  

3. **Set up the MySQL database**:
   Create a MySQL database and update the database connection parameters in `config/config.go`.

4. **Run the application**:
   ```
   go run cmd/student-api/main.go --config=config/local.yaml
   ```

5. **Access the API**:
   The application will be running on `http://localhost:8080`. You can use tools like Postman or curl to interact with the API.

6. **Set Git Initilization**: git init
