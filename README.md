Restless
========

Restless is a Framework for building small to medium (up to 1000 concurrent users) RESTfull APIs in golang and Ecmascript 6.
This can be used to setup whole projects in a few hours without to care about the webstack or complicated configurations.

** Principals

Restless is model-based. To add an endpoint you create a Model. This is basically just a name and up to six backend hooks in JS.
The basic operations on an endpoint are: create, update, get and delete. Creating and updating is called "commiting" an instance of a model.

** Example

For example you have the endpoint 'article'. If you POST a request with any JSON in its body to 'http://myserver/api/article/' the backend creates an instance of 'article' with the given data and calls the 'onPreCommit' hook. In that hook you can add data to your instance (like a lastModified timestamp or something). After that, the instance is written to DB and the 'onPostCommit' hook is called. Here you can log that a instance was created/updated or something. In case of a POST request additionally an unique id is generated returned to the caller of the HTTP request. The PUT method is only a little different since it updates an existing instance. The request would go to 'http://myserver/api/article/<unique_id>'. Especially the 'onPreCommit' and 'onPostCommit' hooks are called. If you send a GET request to 'http://myserver/api/article/<unique_id>' the instance is requested. At first the 'onPreInit' hook of the model is called, after that the instance-data is retrieved from the DB. After getting the data from DB, the 'onPostInit' hook is called. Here you can do serverside post processing of the data, like resolving foreign keys in the data. The last possible action is DELETE. At first the 'onPreRemove' hook is called, then the instance is removed from the DB and after that, the 'onPostRemove' hook is called.

Lets have a look how these hooks could look like:

	app.Model('Article', class {
		onPreInit(){
			console.debug(`onPreInit of instance ${this.__uid} of model ${this.__name}`);
		}
		onPostInit(){
			this.set('example.dynamic.requestTime',(new Date()).toString());
		}
		onPreRemove(){
			console.debug(`onPreRemove of instance ${this.__uid} of model ${this.__name}`);
		}
		onPostRemove(){
			console.debug(`deleted instance ${this.__uid} of model ${this.__name}`);
		}
		onPreCommit(){
			this.set('lastModified',Date.now());
		}
		onPostCommit(){
			console.debug(`commited instance ${this.__uid} of model ${this.__name}`);
		}
	});

See how its ES6? It gets precompiled by babel and then interpreted by otto. There are some limitations, but this is work-in-progress. You can write ES5 too if you want to.

All you need to do to bring this up and running is placing the code in the correct place (/js/models/article/server.js) and to write your configfile like this:

	{
	  "endpoints" : [
	    {
	      "url": "/article",
	      "model": "Article"
	    }
	  ]
	  "db": "/tmp/restless.db",
	  "address": ":8080",
	  "assets": "./bin/assets",
	  "jsBase": "./bin/server.js",
	}

As you see, you specify only a db path (its a leveldb, so no need to run a DB server) and the endpoints. You can also specify a directory where your static assets are taken from. They are available via GET'ing 'http://myserver/assets/filename.ext'. This way you can provide your frontend which uses the API for data storage and retrieval.