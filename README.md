# Sky

## Overview

The Sky database is used to store partially transient, temporal hashes.
This might not seem useful at first but it is a powerful construct when storing and analyzing behavioral data such as clickstreams, logs, and sensor data.
Hashes have a loose schema associated with them to determine whether each property is transient or permanent and the property's data type.

For example, let's say you want to analyze the behavior of users on your web site and combine it with other data in your database.
Each user would have a hash in Sky.
You want to track some properties on the hash over time such as the user's gender or current residence.
These are permanent properties and will be set to a value until changed.

Then there are properties that only exist such as the name of an action being taken or details about that action.
For instance, you might track an `action` property with values of `signup` or `checkout`.
You might also want to track the `purchase_amount` for the checkout action.
These are properties that exist only for the exact moment in which they occur.

Sky stores all this data efficiently for querying so you can aggregate these events at blazingly fast speeds.
For instance, on a typical commodity server you can run a funnel analysis at the rate of ~10MM events/core/second.

## API

### Overview

Sky uses a RESTful JSON over HTTP API.
Below you can find Table, Property, Event and Query endpoints.
The examples below use cURL but there are also client libraries available for different languages.

### Table API

```sh
# Creates an empty table named 'users'.
$ curl -X POST http://localhost:8585/tables -d '{"name":"users"}'
```

```sh
# Deletes the table named 'users'.
$ curl -X DELETE http://localhost:8585/tables/users
```

### Property API

```sh
# List all properties on the 'users' table.
$ curl http://localhost:8585/tables/users/properties
```

```sh
# Add the 'username' property to the 'users' table.
$ curl -X POST http://localhost:8585/tables/users/properties -d '{"name":"username","transient":false,"dataType":"string"}'
```

```sh
# Retrieve the 'username' property from the 'users' table.
$ curl http://localhost:8585/tables/users/properties/username
```

```sh
# Change the name of the 'username' property on the 'users' table to be 'username2'.
$ curl -X PATCH http://localhost:8585/tables/users/properties/username -d '{"name":"username2"}'
```

```sh
# Change the name of the 'username' property on the 'users' table to be 'username2'.
$ curl -X PATCH http://localhost:8585/tables/users/properties/username -d '{"name":"username2"}'
```

```sh
# Delete the 'username2' property on the 'users' table.
$ curl -X DELETE http://localhost:8585/tables/users/properties/username2
```
