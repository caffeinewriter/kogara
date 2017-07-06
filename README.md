# Kogara

## A simple, ugly, Go/Redis-backed URL shortener.

### Naming

Kogara (小柄) is the result of extensive Google Translating variations of the word "short" into different languages, because that's how I like to name my projects. Plus, Kogara sounds pretty damn cool.

### Development Status

This is a hobby project designed to increase my own familiarity with Go, Redis, and Gin. I'm striving to have it become a stable, full-fledged project, but keep that in mind if you decide to use it. 

### Routes

`/` is a simple homepage to shorten links.

`/+/:id` displays a simple counter of the times the link has loaded from the server. (It uses a 301 redirect, so browsers will typically aggressively cache the result)

`/r/:id` is the redirect itself. 

`/check/:id` will return a JSON blob revealing the existence of a particular ID. Mainly implemented so if I decide to add custom URLs.

### Current caveats

* At present solely generates sequential base62 link IDs.
* At present has no way to administer links beyond directly managing them in Redis.
* ~~Does not confirm link existence, or even structure, so you may end up with some weird results.~~
* ~~May be vulnerable to XSS, but most browsers reject execution of JavaScript unless it's directly entered by the user into the address bar.~~
