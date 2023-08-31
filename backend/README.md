# :robot: Chatbot

This solution is based on the JobSity code challenge. It was developed using the following technologies:

- Golang.
- Docker.
- MySQL.
- SQL.
- RabbitMQ.
- Bash script.
- HTML.
- JavaScript.
- Git.
- Makefile.
- Dockerize.
- JWT.

## Prerequisites
Please, ensure you have previously installed:
1. Go 1.19+
2. Make tool.
3. Docker

## Progress

Within the proposed mandatory features, the status of the solution is detailed below:

- [X] Allow registered users to log in and talk with other users in a chatroom.
- [X] Allow users to post messages as commands into the chatroom with the following format __/stock=stock_code__
- [X] Create a decoupled bot that will call an API using the stock_code as a parameter (https://stooq.com/q/l/?s=aapl.us&f=sd2t2ohlcv&h&e=csv, here __aapl.us__ is the __stock_code__)
- [ ] The bot should parse the received CSV file and then it should send a message back into the chatroom using a message broker like RabbitMQ.
- [X] The message will be a stock quote using the following format: "__APPL.US quote is $93.42 per share__". The post owner will be the bot.
__Note:__ Due to lack of time, the __RabbitMQ__ functionality could not be covered.
- [X] Have the chat messages ordered by their timestamps and show only the last __50__ messages.
- [ ] Unit test the functionality you prefer.
__Note:__ Due to lack of time, proper unit and integration tests couldn't be adequately covered in the code.

## Bonus (Optional)
None of the Bonus tasks could be completed.
- [ ] Have more than one chatroom.
- [ ] Handle messages that are not understood or any exceptions raised within the bot.

## Considerations
- [X] We will open 2 browser windows and log in with 2 different users to test the functionalities.
- [X] The stock command won’t be saved on the database as a post.
- [X] The project is totally focused on the backend; please have the frontend as simple as you can.
- [X] Keep confidential information secure.
- [X] Pay attention if your chat is consuming too many resources.
- [X] Keep your code versioned with Git locally.
- [X] Feel free to use small helper libraries.

## How to run

The solution uses `Docker` and `Makefile` to dockerize and deploy the solution transparently, as well as to automate the execution of various tasks useful for running the solution.

#### Steps

The main steps to execute the solution in a local environment are detailed below:

1. ```$ make install```
Installs the necessary dependencies for the solution.

2. ```$ make docker-build ```
Downloads the images and sets up the containers that are part of the `docker-compose.yml` file, then runs the application and services on the ports specified in the `.env` variables file.

3. __Web Browser Execution__
Once the Docker services are built, the application can be accessed through a web browser using the configured port. For example:

```HTML
    http//localhost:8086/
```

4. __Seeded Data__
The application comes with 4 pre-seeded user records, which are automatically populated into the database during the provisioning of the services defined in the ```docker-componse.yml``` file.
The users with which tests can be performed are as follows:
- __User:__ alice | __Password:__ alice
- __User:__ bob | __Password:__ bob
- __User:__ louis | __Password:__ louis

5. ```$ make test```
Runs tests (which are very limited :disappointed:) within the solution.