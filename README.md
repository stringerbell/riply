#Riply

Almost like bitly
Demo: https://riply.fly.dev/not-rickroll

Getting up and running:
Make sure you have docker installed.

```make up```

Running tests:
```make test```

Running integration tests:
```make integ```

Sometimes, if you write an integ test, and dockertest panics, it won't clean up after itself. 
Running `make clean` should clean up everything that's still running 

Deploying:
Make sure you have https://fly.io/docs/getting-started/ installed

`fly deploy`
