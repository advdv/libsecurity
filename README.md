![Timeglass Screenshot](/overview.png?raw=true "Timeglass Screenshot")


# libsecurity
When we're talking about security, the internet is the true wild wild west of today and the problems doesn't seem to be disappearing anytime soon. Until we find a way to create bulletproof server software we'll need a system in place to react as quick as possible to the release of new vulnerabilities. 

Traditionally it can be very difficult for a software vendor to inform their consumers. Especially for open source projects, the product is often acquired (anonymously) through a variety of channels and contacting the user can be downright impossible. Docker can provide a solution, as an universal runtime it can also act as a real time channel for receiving and acting on new vulnerabilities.

## How does it work
fundamental to the workings of libsecurity is the involvement two seperate parties: the vendor and the consumer. The vendor releases new versions of software and is responsible for fixing vulnerability issues and broadcasting the existence of such vulnerability to the consumers. 

### The Vendor
For this party, libsecurity provides a container that allows vendors to test certain images for the existence of a vulnerability on any of the images that are managed by Docker:

```
docker run -v /var/run/docker.sock:/var/run/docker jerbi/cve-check
```

The container will output an example message that the vendor can broadcast across any distributed messaging platform, in this case it tells us that Docker image "9e1ed860cc088ae4b68ce28fb8888739652729e1107054f58dff90979f7dc935" is vulnerable to "CVE-2014-6271": the imfamous Shellshock vulnerability:

```
CVE-2014-6271 in 9e1ed860cc088ae4b68ce28fb8888739652729e1107054f58dff90979f7dc935
```

The vendor is reponsible for broadcasting it across a messaging platform, this demo uses twitter but one could imagine more secure channels such as irc bots or the Docker hub.

![Timeglass Screenshot](/screenshot.png?raw=true "Timeglass Screenshot")

### The Consumer
The consumer, on his part, requires to run a container that monitors certain twitter feeds for new vulnerabilities. The consumer can decide what Twitter user to follow. For example, inside a corporate network one might want to watch to the twitter feed of the security office (advanderveer).

```
docker run -it -v /var/run/docker.sock:/var/run/docker.sock advanderveer/docksec --twitter_user=advanderveer
```

Whenever a message such as _"CVE-2014-6271 in 9e1ed860cc088ae4b68ce28fb8888739652729e1107054f58dff90979f7dc935"_ is broadcasted across the network and picked up by the container above it will act in several ways:

1. It will validated if any of the images managed by the daemon include the broadcasted vulnerability image id (e.g 9e1ed...)
2. If any images or containers are vulnerable it will reply on Twitter stating the fact it is vulnerable
3. The owner of the twitter account can then reply to it with "use latest" to ask the Daemon to fix the problem by itself
2. When asked, it will pull the latest version, if any of the running containers is vulnerable it will restart each container with the latest images. 

### Contributors

- [Amir Jerbi](https://github.com/jerbia)
- [Daniel Sachse](https://github.com/w0mbat)
- [Peter Rossbach](https://github.com/rossbachp)
- [Meir Wahnon](https://github.com/meirwah)
- [Ad van der Veer](https://github.com/advanderveer)
- [Greg Deed](https://github.com/tegbiz)
