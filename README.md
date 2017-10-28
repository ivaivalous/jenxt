# Jenxt
Jenkins Extender - Extend the Jenkins API with custom Groovy scripts

# What is Jenxt

Jenxt acts as a web server, written in go, that lets you easily extend the Jenkins API through Groovy scripts. In a nutshell, it executes Groovy against Jenkins' script console, which gives you full control over Jenkins to make your automation tasks easier.

Normally, it is dangerous to let users run arbitrary scripts as they might take actions that are malicious. With Jenxt, the content of the script is abstracted away from the user and they can't view or modify it. They can only execute it, passing in parameters, and retrieving a response. User requests to Jenxt can be validated.

# How to setup

Installation is meant to be as easy as possible. You either build Jenxt from source with a simple `go build` or download a released executable. The choice of Go as a development language is exactly to make running Jenxt as easy as possible. No frameworks needed, no advanced setup.

```
jenxt
jenxt.json
scripts
```

`jenxt` is the main executable. It starts the server and listens for incoming requests.
`jenxt.json` is the file where you configure the system. `scripts` contains all your user scripts.

To use Jenxt, you simply create Groovy scripts and place them in the scripts directory. Scripts are normal Groovy Jenkins would understand, with some special annotations to register them with the Jenxt server. You'll learn about this later.

To start using Jenxt, first edit `jenxt.json`:

```
{
    "server": {
        "host": "127.0.0.1",
        "port": "8899"
    }
}
```

Then just run the executable.

