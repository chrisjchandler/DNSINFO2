DNS Query API Version 2 only default go libraries
Because sometimes, you just need to query DNS records and you can't do it with a dig or nslookup.

What is this?
It's an API. It queries DNS records. You send it a domain and a record type, it gives you an answer. It's not rocket science.

How to Use
Start the Server: Run go run . or go build followed by ./your-binary-name 

Make Requests: Use curl, Postman, or whatever floats your boat. 


Example specifying a recursive resolver:


curl "http://localhost:8080/dns-query?domain=msitproject.site&nameserver=1.1.1.1"
Replace msitproject.site with your desired domain, and 1.1.1.1 with your preferred nameserver, if you have a preference.

This is sample output from that curl

{"a":["129.146.40.194"],"mx":["eforward5.registrar-servers.com.","eforward4.registrar-servers.com.","eforward1.registrar-servers.com.","eforward2.registrar-servers.com.","eforward3.registrar-servers.com."],"ns":["dns1.registrar-servers.com.","dns2.registrar-servers.com."],"txt":["v=spf1 include:spf.efwd.registrar-servers.com ~all"]}


you can also not specify a recursive and use the servers default resolver 

curl "http://localhost:8080/dns-query?domain=cnn.com&nameserver="

AND you can also hard code the recursive in the api code itself by modifying this block:

const defaultResolver = "8.8.8.8" 

Change 8.8.8.8 to whatever you want. 

Get Results: The API spits back JSON with DNS record info. If there's nothing, it means there's nothing. 



Features
Queries A, AAAA, CNAME, MX, NS, TXT records – the whole gang.
Uses github.com/miekg/dns too cool for the standard net package.
Lets you specify a nameserver


Limitations
Doesn't make coffee.
Won't do your laundry.
Might be overkill if you just want an IP address.
Why Use This?
Maybe you're a sysadmin, or maybe you're just nosy. Whatever your reasons, I don't judge. you crave DNS info, this gets it for you. End of story.

Now with TLS!

Configuration Options
UseTLS: Set this to true to enable TLS (HTTPS), or false to disable it (HTTP).
CertFile: Path to the TLS certificate file.
KeyFile: Path to the TLS private key file.
ServerPort: Port number on which the server will listen.
Enabling TLS
To enable TLS, you need to provide a valid certificate and private key file. Once you have these files:

Set the UseTLS flag to "true" in the configuration.
Provide paths for CertFile and KeyFile with your certificate and key.
Ensure the ServerPort is set to an appropriate value for HTTPS (e.g., 443 or any other if using a custom configuration).

use curl -k to query the api and validate ssl functionality 

e.g "curl -k https://localhost:{port}/dns-query?domain=cnn.com"
Contributing
Found a bug? Open an issue. 

License
Do what you want with it I don't care

