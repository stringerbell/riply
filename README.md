#Riply

Almost like bitly
Demo: https://riply.fly.dev/not-rickroll

Getting up and running:
Make sure you have [docker](http://riply.fly.dev/docker) installed.

```make up```

Running tests:
```make test```

Running integration tests:
```make integ```

Sometimes, if you write an integ test, and [dockertest](http://riply.fly.dev/771a3395bc) panics, it won't clean up after itself. 
Running `make clean` should clean up everything that's still running 

Deploying:

Make sure you have [fly](http://riply.fly.dev/getting-started-with-fly) installed

`fly deploy`
