Assuming you guys already have GO and Docker installed, here are the commands I used to get the server running:

1. docker build -t go-docker-server .
2. docker run -p 8080:8080 go-docker-server

Steps 1 and 2 step is for building out the server and running it.

Here are the url's for the endpoints I used to call my own back-end server:

POST endpoint: http://localhost:8080/receipts/process
GET endpoint: http://localhost:8080/receipts/{id}/points

Replace {id} with the id that gets returned from the POST endpoint

IMPORTANT NOTE: In the instructions, it says if the time is AFTER 2pm and BEFORE 4pm, then to add the points to the score. I wasn't entirely sure if this was inclusive for the 2pm and 4pm times, but since the wording said strictly AFTER and BEFORE, I assumed it wasn't inclusive.

Proof it works:

![image](https://github.com/user-attachments/assets/637599d7-f596-444f-9a6c-598329b870a7)

![image](https://github.com/user-attachments/assets/4788a90b-04bc-41b2-9557-29c3bc542f52)
