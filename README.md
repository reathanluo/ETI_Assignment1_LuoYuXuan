# ETI Assignment 1 - Luo YuXuan

The assignment implements a ride-sharing platform using microservices. It has two user groups: passengers and drivers. Passenger can initiate a ride request and the driver will be responsible for the ride.

<img width="620" alt="image" src="https://user-images.githubusercontent.com/72963379/208301661-ff1a4aa3-b971-4369-aa01-2d9959a5dde2.png">

## Design consideration of the microservices
The design contains three microservices, which are authentication, user action and trip. They will be used to perform different tasks while staying cohesive

### Authentication
- The microservice is in charge of user sign up and login (check input credentials)

### User action
- The microservice is in charge of editing user profiles

### Trip
- The microservice is in charge of creating a new trip (by the passenger) and start and end a trip (by the driver). It could also help the passenger retrieve the past trip history

## Architecture diagram

![design](https://user-images.githubusercontent.com/72963379/208302239-88a21858-fa9d-4d9d-940b-2939b156cc00.jpg)

## Instructions for setting up and running the microservices
1. Download the repository and the sql script
2. Set up the database with the script in MySQL Workbench
3. Install the browser extension Moesif Orgin & CORS Changer on your browser and turn it on
4. Run the three microservices with command **`go run file.go`** 
5. Open the browser and type in the link http://127.0.0.1:5500/frontend/login.html (you can start live server in vscode as well)


