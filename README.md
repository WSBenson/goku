# goku-http-server

This is my Goku http server written in go and containerized using docker.

The purpose of this server is to take a JSON post request and return a response ranking each character by their power level.

At the moment it will simply add "Hello, " + the name in the JSON body

For Example:

The JSON file for the POST request should look something like this:

```json
{
	"name": "Goku"
}
```

Which will return a response of:

```json
{
    "Response": "Hello, Goku"
}
```
This will be changed in the future once the program is ready to rank characters by power levels.

## Gettting Started

### Building

To Build this project follow these steps:
- Make sure you have docker installed
- Clone my Git repository where ever you want
- Then change directory to wsbenson/goku-http-server
- Run these two commands back to back in bash:
    ```bash
       docker build -t goku .

       docker run --env PORT=3000 -p 3000:3000 goku
                            OR
                docker run -p 3000:3000 goku
    ```
- By default, even if you leave out the --env PORT=3000 part of the command, the server will run on port 3000. So if port 3000 is occupied for you, change the PORT env variable in the second command, EX: --env PORT=4005
- Also make sure to update the first port number in that second command with your new port. If you changed the port to 4005 it would become: 4005:3000
- Don't actually use port 4005, it won't work. This was just used for a formatting example.

- The http server should now be running, you can check by going to:
    `http://localhost:3000/test`   <- Change the port number to match yours



### Testing

- Send a JSON POST request to: `http://localhost:3000/test` or whatever port you used
- You should get a response similar to the one above the "Getting Started" section.