Initialize my schedule downloader for stiny socks hockey

This code will pull games from the stinky socks website and add it to your google calendar.

Website: https://www.stinkysocks.net/

Prepare 
```
go get github.com/PuerkitoBio/goquery
go get -u google.golang.org/api/calendar/v3
go get -u golang.org/x/oauth2/google
```

Please go to [this URL](https://developers.google.com/workspace/calendar/api/quickstart/go) and follow instructions for:
- [Enable API](https://developers.google.com/workspace/calendar/api/quickstart/go#enable_the_api)
- [Configure the OAuth consent screen](https://developers.google.com/workspace/calendar/api/quickstart/go#configure_the_oauth_consent_screen)
- [Authorize credentials for a desktop application](https://developers.google.com/workspace/calendar/api/quickstart/go#authorize_credentials_for_a_desktop_application)

Fix for not being a test user:

Go back to Google Cloud Console → APIs & Services > OAuth consent screen.

Scroll down to Test users.

Click + Add Users.

Add your Gmail address (the one you’re logging in with).

Save.
  
How to run code:
```
go run cmd/main.go
```

How to run tests:
```
go test ./...
```