# Go Application Deployment and MySQL Setup

## 1. Explanation of Source Code

The Go application is a simple web service that interacts with a MySQL database to log user requests. The application exposes two primary endpoints:

- **`/current-time`**: Logs the current time (in the Toronto time zone), along with the username and IP address, into a MySQL database.
- **`/logs`**: Retrieves and returns all the logs from the MySQL database.

The application uses the `github.com/go-sql-driver/mysql` driver to connect to the MySQL database and perform SQL queries.

Key parts of the application:
- **Database Connection**: The Go application reads environment variables for database credentials and uses these to form a DSN (Data Source Name) for the MySQL connection.
- **Time Handling**: The application uses the Go `time` package to handle time zones, specifically converting the current time to Toronto time (`America/Toronto`).
- **Logging**: The `/current-time` endpoint logs the username, IP address, and current time to the MySQL database in a table called `time_log`.

---

## 2. Instructions for Setting Up and Running the Application

### **Prerequisites**
- Docker and Docker Compose installed on your system.

### **Step-by-Step Setup:**

1. **Clone the Repository** (if not already done):
   ```bash
   git clone https://github.com/your-repository/go-app.git
   cd go-app
   ```

2. **Update Docker Compose File**:
   Update the `docker-compose.yml` file in the project with your choice of Username and Password:

   ```yaml
   services:
     mysql:
       image: mysql:latest
       container_name: mysql-container
       restart: always
       environment:
         MYSQL_ROOT_PASSWORD: admin@4321
         MYSQL_DATABASE: UserData
         MYSQL_USER: admin
         MYSQL_PASSWORD: admin@4321
       ports:
         - "3306:3306"
       volumes:
         - ./sql-scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
       
   
     backend:
       build:
         context: .                     # The directory where Dockerfile is located
       container_name: goApp_container
       restart: always
       ports:
         - "8080:8080"                  # Expose the app on port 8080
       environment:
         DB_HOST: mysql                 # MySQL service name (DNS resolution via Docker network)
         DB_PORT: 3306                  # MySQL port
         DB_USER: root                  # MySQL user
         DB_PASSWORD: admin@4321        # MySQL password
         DB_NAME: UserData              # MySQL database name
       depends_on:
         - mysql                        # Ensure the MySQL container starts first
   ```

3. **Build and Start the Containers**:
   Build and start the containers using the following Docker Compose command:
   ```bash
   docker-compose up --build
   ```

4. **Access the Application**:
   - Once the containers are up, you can access the Go application by visiting `http://localhost:8080` in your browser.
   - Use `localhost:8080/logs` to view the logs stored in the database.

5. **Stop the Containers**:
   To stop the containers, use:
   ```bash
   docker-compose down
   ```

---

## 3. Setting Up the MySQL Database and Table

### **MySQL Setup Using Docker**

The MySQL container with the Database and Table is automatically set up when you run the `docker-compose up` command (using init.sql), but if you may need to manually set up the database:
#### **Step 1: Access MySQL Container**

To connect to the MySQL container, run the following command to open the MySQL shell:
```bash
docker exec -it mysql_container mysql -u root -p
```
Enter the password for the `root` user (`admin@4321` in this case).

#### **Step 2: Create the Database**

If the database `UserData` is not already created by the `docker-compose.yml` file, you can create it manually by running:
```sql
CREATE DATABASE UserData;
```

#### **Step 3: Create the `time_log` Table**

Run the following SQL query to create the `time_log` table that the Go application uses to store logs:
```sql
CREATE TABLE time_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    logged_time DATETIME,
    username VARCHAR(255),
    ip_address VARCHAR(255)
);
```

#### **Step 4: Verify the Table**

You can verify that the table has been created by running:
```sql
SHOW TABLES;
```

This should display the `time_log` table in the list.

---

## 4. Docker Commands Used

### **Docker Compose Commands**
- **Build and Start Containers**:
  ```bash
  docker-compose up --build
  ```

- **Stop Containers**:
  ```bash
  docker-compose down
  ```

- **Start Containers (without rebuilding)**:
  ```bash
  docker-compose up
  ```

- **View Logs**:
  To view logs from the containers:
  ```bash
  docker-compose logs
  ```

- **Run the Go Application  Manually** (if you want to run the Go app container without Compose):
  ```bash
  go run main.go //this is your localhost server
  docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=admin@4321 -e MYSQL_DATABASE=UserData -e MYSQL_USER=admin -e MYSQL_PASSWORD=admin@4321 -p 3306:3306 -d mysql:lates
  ```

---

## Conclusion

With this setup, you can easily run your Go application with MySQL using Docker Compose. The instructions cover:
- How to set up and run the Go application and MySQL database in containers.
- How to create the MySQL database and table.
- The necessary Docker commands for building, running, and managing containers.

