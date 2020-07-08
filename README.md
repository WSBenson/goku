# goku-http-server

This is my Goku http server written in go and containerized using Docker.

The purpose of this server is to return a response or message that ranks each character by their power level. You have complete control over which characters  you choose to create, and eventually battle.

There are two different ways to use this program: simply sending a JSON POST request to the HTTP server at the address that you can specify, or by using elastic search with  some simple ready-to-go Cobra commands from your shell.

While there are two different ways to get the results of the program, there are also two different ways to build it: using docker-compose, or with Minikube and Helm 3.

I am first going to talk about building the program to get the HTTP server up and running.

- Before we go any further, please clone this repository somewhere and make sure you have the latest version of Golang and Docker installed
- If you plan to use Minikube and Helm 3, please install them as well

## Build Option 1: Simple

- Go to where ever you have cloned the repository and cd (change directory) into the goku folder (if you are not already there)
- Run the command:
  ```shell
  docker-compose up --build
  ```
  If you run into this beautiful error for some reason:
  ```shell
  docker.credentials.errors.InitializationError: docker-credential-gcloud not installed or not available in PATH [69711] Failed to execute script docker-compose
  ```
  Then all you have to do is run this command:
  ```shell
  rm ~/.docker/config.json
  ```
  And run ```docker-compose up --build``` again

  You now should have multiple containers running: one for the goku HTTP server, one for Elasticsearch, one for Kafka, and a final one for Zookeeper

  You may verify this by running this command:
  ```shell
  docker ps
  ```
  
  Now the Goku HTTP server will be running at http://localhost:3000/goku by default.
  
  Elasticsearch will also be running at http://localhost:9200/ by default.

## Build Option 2: A Little More Confusing

Unfortunately this method does not currently work with Elasticsearch since I had difficulty finding an Elasticsearch chart that would fit my needs and work with Helm 3.

- Go to where ever you have cloned the repository and cd (change directory) into the goku folder (if you are not already there)
- Run this command to start up Minikube:
  ```shell
  minikube start
  ```
- Run this to get the Kafka Helm chart:
  ```Shell
  helm dep update ./goku-chart
  ```
- Then run this to build the containers for the goku server, Kafka, and Zookeeper:
  ```shell
  helm install anynameyouwant ./goku-chart --set service.type=NodePort
  ```
- Copy the command it prints out under "1. Get the application URL by running these commands:"
- It should look something like this:
  ```shell
    export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=goku-chart,app.kubernetes.io/instance=anynameyouwant" -o jsonpath="{.items[0].metadata.name}")
  echo "Visit http://127.0.0.1:8080/goku to use your application"
  kubectl --namespace default port-forward $POD_NAME 8080:3000
  ```
- Paste that command into your Shell to run it

This will port forward  port 8080 to port 3000 so you can access the server at: 
http://127.0.0.1:8080/goku

- If you would like to view the containers being built you can open another Shell tab/window and run ```kubectl get pods``` for a listing inside Shell or ```minikube dashboard``` for a more useful graphical interface.

The minikube dashboard will show you your successful deployment and all the pods that have been created.

WARNING: Again, this will not work with Elasticsearch at the moment so you can only see the results of the fighter evaluation with the POST Request method.



## Seeing The Results: POST Request

Upon a POST request to the Goku HTTP server, a unique message evaluating the different power levels of the characters you entered will be returned.

For Example:

The JSON file for the POST request should look something like this:

```json
{
	"fighters": [
			{
				"name": "Goku",
				"power": 150000000
			}
		]
}
```

Which will return a response of:

```json
{"message":"The scouter says Goku's power level is over 9000! You better start running."}
```

You may add as many fighters as you want, just be careful with the formatting of your JSON. But don't worry it will return a descriptive logger error if you end up having some typos.

A simple program like Postman makes it super easy to send post requests to the server. Just make sure you can reach it first at http://localhost:3000/goku if you built it with docker-compose or http://127.0.0.1:8080/goku if you built it with Minikube and Helm.

## Seeing the Results: Cobra Commands, Elasticsearch, and Kafka

If you don't feel like sending a JSON POST message you can evaluate your fighters through shell with the Cobra commands I have added.
 - Ensure you have successfully built the containers with one of the two options above (at this point in time Elasticsearch will only work with the docker-compose method)
 - Open a new Shell tab or window (leave your previous tab/window open)
 - In the goku folder run the command ```goku``` to see all available Cobra commands
 - From there, create the Elasticsearch index by running the command:
   ```shell
   goku es
   ```
- Then run the goku addf command with the --help flag to see available flags:
  ```shell
  goku addf --help
  ```
- You will see that the flag --name allows you to specify your fighters name and --level their power
- If I wanted to create Goten with a power level of 1000 I would run:
  ```shell
  goku addf --name Goten --level 1000
  ```
- This will index Goten with a power level of 1000 in Elasticsearch and also create a new Kafka producer that will produce a message to be consumed by a Kafka Consumer that can do something with that information.
- After you are done adding all of the fighters you want to be evaluated (in battle)...
- You can run:
  ```shell
  goku battle
  ```
Which will return a message with the result of the evaluation by querying all fighters in your Elasticsearch index.

If you want to see each fighter that was indexed to Elasticsearch you can visit: http://localhost:9200/fighters/_search


## Credits

Shout out to Decipher Technology Studios for being a mint company.

Shout out to Joseph Knoebel (knoebelja) for being patient with me on this project and reviewing too many PRs.

                        And shout out to Alec, he sucks.