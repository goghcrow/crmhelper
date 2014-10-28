package main

/*
@from http://stackoverflow.com/questions/14661511/setting-up-proxy-for-http-client
you could set the HTTP_PROXY environment variable, if you do this Go will use it by default.

Bash:

export HTTP_PROXY="http://proxyIp:proxyPort"
Go:

os.Setenv("HTTP_PROXY", "http://proxyIp:proxyPort")
You could also construct your own http.Client that MUST use a proxy regardless of the environment's configuration:

proxyUrl, err := url.Parse("http://proxyIp:proxyPort")
myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
This is useful if you can not depend on the environment's configuration, or do not want to modify it.

You could also modify the default transport use by the "net/http" package. This would effect your entire program (including the default HTTP client).

proxyUrl, err := url.Parse("http://proxyIp:proxyPort")
http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
*/
