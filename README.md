# Task: Develop a RESTful API with JWT Authentication, Rate Limiting, Data Validation, Role-Based Access Control, and Caching

## Objective:

Create a RESTful API that includes user authentication with JSON Web Tokens (JWT), rate limiting, data validation, role-based access control (RBAC), and caching. The API will have endpoints for user registration, login, user profile management, and data retrieval.

## Requirements:

- Language: Use any backend language/framework of your choice (e.g., Python with Flask/Django, Node.js with Express, Java with Spring Boot, etc.).
- Endpoints:
    POST /register: Register a new user with roles.
    POST /login: Authenticate a user and return a JWT.
    GET /profile: Retrieve the user profile, protected by JWT authentication.
    PUT /profile: Update the user profile, protected by JWT authentication.
    GET /admin/data: Retrieve data for admin users, protected by JWT authentication and RBAC.
- Authentication: Use JWT for securing endpoints.
- Rate Limiting: Implement rate limiting to restrict each user to a maximum of 5 requests per minute to the /profile and /admin/data endpoints.
- Data Validation: Validate input data for registration, login, and profile management endpoints.
- Role-Based Access Control: Implement RBAC to allow only admin users to access the /admin/data endpoint.
- Caching: Implement caching for the /admin/data endpoint to improve performance.

## Example Requests and Responses:

POST /register:

json
```
{
    "username": "user1",
    "password": "password123",
    "role": "user"
}
```
POST /login:

json
```
{
    "username": "user1",
    "password": "password123"
}
```
Response:

json
```
{
    "token": "your-jwt-token"
}
```
GET /profile (with JWT in Authorization header):

json
```
{
    "username": "user1",
    "role": "user",
    "email": "user1@example.com"
}
```
PUT /profile (with JWT in Authorization header):

json
```
{
    "email": "newemail@example.com"
}
```
GET /admin/data (with JWT in Authorization header):

json
```
{
    "data": "Admin-specific data"
}
```
## Instructions:

### Project Setup:

Set up a basic project structure for your chosen backend language/framework.
Configure a database (e.g., PostgreSQL, MongoDB, SQLite) to store user credentials and roles.
### User Registration:

- Create the POST /register endpoint to register a new user with a specified role.
- Validate the input data (e.g., ensure username, password, and role are provided and meet certain criteria).
- Hash the user's password before storing it in the database.
- User Login and JWT Authentication:


- Create the POST /login endpoint to authenticate a user.
- Validate the input data.
- Generate a JWT upon successful authentication and return it in the response.
### User Profile Management:

- Create the GET /profile endpoint to retrieve the authenticated user's profile.
- Create the PUT /profile endpoint to update the authenticated user's profile.
- Protect these endpoints with JWT authentication.
- Implement data validation for the profile update.

### Admin Data Retrieval:

Create the GET /admin/data endpoint to retrieve admin-specific data.
Protect this endpoint with JWT authentication and RBAC.
Ensure that only users with the admin role can access this endpoint.

### Rate Limiting:

- Implement middleware or a function to track the number of requests made by each user based on their JWT.
- Limit each user to a maximum of 5 requests per minute for the /profile and /admin/data endpoints.
- Return a 429 status code and a message "Too Many Requests" if the limit is exceeded.

### Caching:

- Implement caching for the /admin/data endpoint to store the response for a short duration (e.g., 1 minute) using an in-memory store like Redis or a simple in-memory cache in your application.

### Testing:

- Test the API using a tool like Postman or curl.
- Verify that user registration, login, JWT authentication, role-based access control, rate limiting, and caching work as expected.

### Recording:

## After completing the task, please document the following:

- Steps taken to set up the project, including user registration, login, JWT authentication, profile management, role-based access control, rate limiting, and caching.
- Code snippets for the key parts of the implementation.
- Testing method and results, including logs or screenshots showing successful and unsuccessful requests.
- Explanation of any challenges faced and how they were overcome.

This task will challenge his ability to implement advanced backend features such as authentication, RBAC, rate limiting, data validation, and caching, demonstrating his extensive technical skills.