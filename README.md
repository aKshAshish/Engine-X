# Nginx (Engine X)
**nginx** has one master process and several worker processes. The main purpose of the master process is to read and evaluate configuration, and maintain worker processes. Worker processes do actual processing of requests. nginx employs event-based model and OS-dependent mechanisms to efficiently distribute requests among worker processes. The number of worker processes is defined in the configuration file and may be fixed for a given configuration or automatically adjusted to the number of available CPU cores.

## Configuration File’s Structure
- **nginx** consists of modules which are controlled by directives specified in the configuration file. 
- **Directives** are divided into simple directives and block directives. 
    - A *simple directive* consists of the name and parameters separated by spaces and ends with a semicolon (;). 
    - A *block directive* has the same structure as a simple directive, but instead of the semicolon it ends with a set of additional instructions surrounded by braces ({ and }). If a block directive can have other directives inside braces, it is called a **context** (examples: events, http, server, and location).
- Directives placed in the configuration file outside of any contexts are considered to be in the `main` context. The `events` and `http` directives reside in the `main` context, `server` in `http`, and `location` in `server`.

- The rest of a line after the `#` sign is considered a *comment*.

## How Nginx processes a request?
nginx first decides which server should process the request. Suppose all the servers in the *configuration* listen on port ::80 
then nginx tests only the request’s *header field “Host”* to determine which server the request should be routed to. If its value does not match any server name, or the request does not contain this header field at all, then nginx will route the request to the ***default server*** for this port.
```
server {
    listen      80;
    server_name example.org www.example.org;
    ...
}
server {
    listen      80;
    server_name example.net www.example.net;
    ...
}
server {
    listen      80;
    server_name example.com www.example.com;
    ...
}
```
In the configuration above first server is considered as the default server and it is the default behaviour of nginx. The default server can be set by adding *default_server* parameter in the **listen** directive.

## How to prevent processing requests with undefined server names?
To drop the requests that have undefined "Host" in their headers we can simply add a server that just drops these requests.
```
server {
    listen      80;
    server_name "";
    return      444;
}
```

## Mixed name-based and IP-based virtual servers
```
server {
    listen      192.168.24.10:80 default_server;
    server_name example.org www.example.org;
    ...
}
server {
    listen      192.168.32.10:80;
    server_name example.net www.example.net;
    ...
}
server {
    listen      192.168.40.10:80;
    server_name example.com www.example.com;
    ...
}
```
- In the above configuration, nginx first tests IP and Port against the listen directives of the server blocks. 
- It then tests the “Host” header field of the request against the server_name entries of the server blocks that **matched** the IP address and port.
- If the server name is not found, the request will be processed by the default server. *For example*, a request for www.example.com received on the 192.168.1.1:80 port will be handled by the default server of the 192.168.1.1:80 port, i.e., by the first server, since there is no www.example.com defined for this port.

## A simple PHP sit configuration
Now let’s look at how nginx chooses a location to process a request for a typical, simple PHP site:
```
server {
    listen      80;
    server_name example.org www.example.org;
    root        /data/www;

    location / {
        index   index.html index.php;
    }

    location ~* \.(gif|jpg|png)$ {
        expires 30d;
    }

    location ~ \.php$ {
        fastcgi_pass    localhost:9000;
        fastcgi_param   SCRIPT_FILENAME
                        $document_root$fastcgi_script_name;
        include         $fascgi_params;
    }
}
```
- nginx first searches for the most specific prefix location given by ***literal strings*** regardless of the listed order and remebers it.
- Then nginx checks locations given by regular expression in the order listed in the configuration file. The first matching expression stops the search and nginx will use this location. 
- If no regular expression matches a request, then nginx uses the most specific prefix location found earlier. In the configuration above the only prefix location is “/” and since it matches any request it will be used as a last resort.

### Examples

- **Map request `/logo.png`**
    > The request will map to `/` prefix, and then nginx will try to match it with the regular expressions where it will match to the first listed regular expression `\.(gif|jpg|png)$` the nginx will stop here and return the file `/data/www/logo.png`

- **Map request `/index.php`**
    > Again this will start with matching with prefix `/`. Then nginx will start try to match it with regular expression in the listed order. The path will match with `\.php$`. Therefore, it will be handled by fastcgi_server hosted on `localhost:9000`.
- **Map request `/about.html`**
    > The request only matches to the prefix `/` hence using the root directive nginx will return the file `/data/www/about.html` to the client.
- **Map request `/`**
    > Handling a request “/” is more complex. It is matched by the prefix location “/” only, therefore, it is handled by this location. Then the index directive tests for the existence of index files according to its parameters and the “root /data/www” directive. If the file /data/www/index.html does not exist, and the file /data/www/index.php exists, then the directive does an internal redirect to “/index.php”, and nginx searches the locations again as if the request had been sent by a client. As we saw before, the redirected request will eventually be handled by the FastCGI server.