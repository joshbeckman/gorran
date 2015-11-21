# gorran
So if I'm going to learn, I must do it by listening.

Powering [on.narro.co](http://on.narro.co).

## Deploy notes
```
go get github.com/kr/godep
godep save
git add -A
git commit -m "dependencies"
# maybe
heroku config:set BUILDPACK_URL=https://github.com/kr/heroku-buildpack-go.git
```

## ENV Config Vars
`MONGO_URI`
