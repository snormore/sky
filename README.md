# Sky -- Open Source Behavioral Database

### Introduction

Sky is an open source database used for high performance analysis of [behavioral data](http://skydb.io/blog/introduction-to-behavioral-databases.html).
This includes clickstream data, log data, sensor data, etc.

Sky's architecture is designed to address two use cases:

1. Fast analysis of an entire dataset.
1. Real-time analysis of a single object.

The database works in the context of objects.
An example of an object is a user of a web site or application.
Objects can be anything that has state and performs actions over time.
These actions and state changes are referred to as events within Sky.
For a given object, Sky can analyze the actions the object performs or how the object's state changes over time.

Operating in the context of a single object at a time has advantages and trade offs.
It's advantages are speed and simplicity of code: Sky can easily aggregate tens of million of events per second per core.
The tradeoff is that Sky is not general purpose.
It's not relational, it's not good for graph analysis, it isn't mean to store documents.
That's okay.
It's not meant to replace those technologies.


### Getting Started

The best way to get up and running with Sky is to read through the [Documentation][] and [Getting Started][] pages.
Sky only requires ZeroMQ as a dependency so it's easy to get up and running.


### More Information

If you have questions about Sky, feel free to visit the [Github repository][], [e-mail the mailing list](mailto:sky@librelist.com) or [find me on Twitter](https://twitter.com/benbjohnson).

  [Github repository]: https://github.com/skydb/sky
  [Documentation]: http://skydb.io/docs
  [Getting Started]: http://skydb.io/docs/getting-started.html