- create in @internal/app/start.go startServer function
- create http server 
- see which config values (whole address, stuff like that) can be externalized, put them in config.json and then load them in ConfigData etc. 
- create one Root endpoint 
- all the endpoint declaration is in @internal/endpoint/router.go
- each endpoit stays in file called get_root.go so get is the http method, root is info about entry point

