# Jenxt
Jenkins API Extender - Extend the Jenkins API with custom Groovy scripts you can run en masse on your Jenkins servers.

# What is Jenxt

Jenxt acts as a web server, written in go, that lets you easily extend the Jenkins API through Groovy scripts. In a nutshell, it executes Groovy against Jenkins' script console, which gives you full control over Jenkins to make your automation tasks easier.

Normally, it is dangerous to let users run arbitrary scripts as they might take actions that are malicious. With Jenxt, the content of the script is abstracted away from the user and they can't view or modify it. They can only execute it, passing in parameters, and retrieving a response. User requests to Jenxt can be validated.

# How to setup

Installation is meant to be as easy as possible. You either build Jenxt from source with a simple `go build` or download a released executable. The choice of Go as a development language is exactly to make running Jenxt as easy as possible. No frameworks needed, or advanced setup.

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
        "port": 8899
    },
    "remotes": [{
        "name": "Jenkins Production Server",
        "url": "https://ci.example.com",
        "username": "user123",
        "password": "My Password",
        "labels": ["default"]
    }, {
        "name": "Jenkins Server in Paris",
        "url": "https://ci2.example.com",
        "username": "user",
        "password": "PASSWORD",
        "labels": ["default", "Europe"]
    }]
}
```

Under server > port, you can change the port where the server is started. Use `remotes` to specify a list of remote Jenkins servers you would like to manage. The `remotes` array consists of server descriptions.

`name` is a descriptive name for the server, like "Jenkins Production Server". This name will be returned in responses so you know how each server responded to your request. In `url` you specify the full address of your Jenkins server. `username` and `password` are the credentials of a user who has the permission to run Groovy scripts.

To avert the risk of committing your passwords to version control by mistake, you can separate passwords from the rest of the configuration. Simply remove the passwords from the configuration file so that it looks like this:

```
...
{
    "name": "Jenkins Server in Paris",
    "url": "https://ci2.example.com",
    "username": "user",
    "labels": ["default", "Europe"]
}
...
```

And then create a new JSON file with a name of your choice, with the following structure:

```
{
    "remotes": [{
        "name": "Jenkins Server in Paris",
        "password": "YouR Pa$$w0RD"
    }]
}
```

The configuration and secrets file can be stored anywhere on the file system your program can access, and you can give it the location by running it with two arguments:

```
./jenxt <path to jenxt.json> <path to secrets.json>
```

The last component of the config file is the `labels` list that allows you to group your servers. You can use this to enable executing scripts only for the groups you'd like. If you want to store a server but never execute any scripts on it, you can just not set any labels for it.

# Creating scripts

Groovy scripts you create go in the `scripts` directory. The repository contains some scripts for example purposes that you can safely delete. To enable a script, you just have to place it in the folder. Please only put Groovy scripts in there as Jentx will try to register any file it finds in the directory.

You can add, change or remove scripts without restarting Jenxt. Checks for changes are run every ten seconds. This should allow you to just run the server somewhere and sync it though your version control system.

A script looks like this:

```
/*
<jenxt>
{
    "expose": "epoch-time"
}
</jenxt>
*/

return new Date().getTime()
```

Wrapped in `<jenxt>...</jenxt>` is the so called meta of the script. It instructs Jenxt how to run the script. Inside is a simple JSON object that may currently contain the following settings:

 - `expose` - this is the endpoint one needs to access so they can run the script. In the above example, an HTTP request to `http://127.0.0.1:8899/epoch-time` will run the given script against the configured Jenkins instances that have the "default" label. To execute against instances labelled "XYZ", add a `label` parameter to the request, like in `http://127.0.0.1:8899/epoch-time?label=XYZ`.
  - `jsonResponse` - Set this to false or omit to get Jenkins' response as a string. Set it to true to return a JSON. *Note*: if jsonResponse is `true` and the response can't be converted to a JSON (for example, when it is a normal string), the response will be returned as `null`.

**Note**: In next releases, additional configuration will be added for parameters, request validation, etc.

# Responses

Running the script from the above example against two machines yields a response similar to this one.

```
{
   "results":[
      {
         "server":"Jenkins Production Server",
         "response":"1509269708370"
      },
      {
         "server":"Jenkins Server in Paris",
         "response":"1509269708721"
      }
   ]
}
```

If there was an error in executing the script, it will be restured in the `response` field, too. Additionally, there will be an `error` fiels with the value of `true`. Like:

```
{
  "results": [{
    "server": "Jenkins Server in Pune",
    "response": "Authentication failed",
    "error": true
  }]
}
```

And if we have the below, more complex script, giving us the SCM configuration for all jobs that have one,

```
/*
<jenxt>
{
    "expose": "scm-triggers",
    "jsonResponse": true
}
</jenxt>
*/

import groovy.json.*

def result = [:]

Jenkins.instance.getAllItems().each { it ->
  if (!(it instanceof jenkins.triggers.SCMTriggerItem)) {
    return
  }

  def itTrigger = (jenkins.triggers.SCMTriggerItem)it
  def triggers = itTrigger.getSCMTrigger()

  triggers.each { t->
    def builder = new JsonBuilder()
    result[it.name] = builder {
      spec "${t.getSpec()}"
      ignorePostCommitHooks "${t.isIgnorePostCommitHooks()}"
    }
  }
}

return new JsonBuilder(result).toPrettyString()
```

The result would resamble this:

```
{
   "results":[
      {
         "server":"Jenkins Production Server",
         "response":{
            "Build Project A":{
               "ignorePostCommitHooks":"false",
               "spec":"@hourly"
            },
            "ivodb.deploy":{
               "ignorePostCommitHooks":"false",
               "spec":"H 15 * * *"
            }
         }
      }
   ]
}
```

# Roadmap

 - Unit tests :) :)
 - Parameters in requests
 - Request parameters and body validation
 - DONE - Automatic pick up of configuration and script changes
 - Authentication of requests to Jenxt
 - Password encryption
 - CI
