## podcast-publisher

- A simple command-line application to self-host podcasts

I self-host a podcast feed to transfer lectures and audiobooks to my phone.

### Installation
```
go get -u github.com/marvelm/podcast-publisher
go install github.com/marvelm/podcast-publisher
```

### Usage
Assuming you have some files in `/my-media`:

```
my-media
└── book
    ├── 1
    │   ├── chapter1.mp3
    │   └── chapter2.mp3
    └── 2
        ├── chapter3.mp3
        └── chapter4.mp3
```


`podcast-publisher -url "https://example.com/podcasts" -folder /my-media/ > feed.xml`

This will create a podcast with episodes defined in `feed.xml`

#### Web server configuration
You'll need a webserver to serve your media and the podcast feed.
Here's an example of a NGINX configuration:

```nginx
server {
  location /podcasts {
    alias /my-media;
    autoindex on;
  }
  location /feed.xml {
    alias /feed.xml;
  }
}
```

Then, you can point your podcast app to `https://example.com/feed.xml`
