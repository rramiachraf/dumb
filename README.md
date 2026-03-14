# dumb
With the massive daily increase of useless scripts on Genius's web frontend, and having to download megabytes of clutter, [dumb](https://github.com/rramiachraf/dumb) tries to make reading lyrics from Genius a pleasant experience, and as lightweight as possible.

<a href="https://codeberg.org/rramiachraf/dumb"><img src="https://img.shields.io/badge/Codeberg-%232185d0" /></a>

![Screenshot](https://raw.githubusercontent.com/rramiachraf/dumb/main/screenshot.png)

## Installation & Usage
### Docker
```bash
docker run -p 8080:5555 --name dumb ghcr.io/rramiachraf/dumb:latest
```

### Build from source
[Go 1.24+](https://go.dev/dl) is required.
```bash
git clone https://github.com/rramiachraf/dumb
cd dumb
make build
./dumb
```

#### Notes:
- The default port is 5555, you can use other ports by setting the `PORT` environment variable.
- Genius servers are behind a Cloudflare reverse proxy, which means certain IPs won't be able to send requests, to partially mitigate this, you can specify a proxy by setting the `PROXY` variable (must be a valid URI).

## Public Instances
| URL | Tor | I2P | Region | CDN? | Operator |
| --- | :----: | :----: | :----: | :----: | --- |
| <https://dumb.ducks.party> | No | No | NL | No | https://ducks.party |
| <https://dumb.privacydev.net> | [Yes](http://dumb.g4c3eya4clenolymqbpgwz3q3tawoxw56yhzk4vugqrl6dtu3ejvhjid.onion) | No | FR | No | https://privacydev.net |
| <https://dumb.hyperreal.coffee> | No | No | US | No | https://hyperreal.coffee |
| <https://dm.vern.cc> | [Yes](http://dm.vernccvbvyi5qhfzyqengccj7lkove6bjot2xhh5kajhwvidqafczrad.onion) | [Yes](http://vernxpcpqi2y4uhu7to4rnjmyjjgzh3x3qxyzpmkhykefchkmleq.b32.i2p) | US | No | https://vern.cc |
| <https://dumb.lunar.icu> | No | No | DE | Yes | @MaximilianGT500 |
| <https://dumb.privacyfucking.rocks> | No | No | DE | - | https://privacyfucking.rocks |
| <https://sing.whatever.social> | No | No | US/DE | Yes | Whatever Social |
| <https://dumb.bloat.cat> | No | No | DE | No | https://bloat.cat |
| <https://dumb.jeikobu.net> | No | No | DE | Yes | https://jeikobu.net |
| <https://dumb.canine.tools> | No | No | US | No | https://canine.tools |
| <https://lyr.dc09.ru> | No | No | RU | No | @DarkCat09 |
| <https://db.kuuro.net> | [Yes](http://db.kuuro5abqix6tfku77wj32srkicgqh3f7ro77ctpda7ub7a7mlv4jsid.onion/) | No | US | No | https://kuuro.net/ |
| <https://dumb.artemislena.eu> | [Yes](http://dumb.lpoaj7z2zkajuhgnlltpeqh3zyq7wk2iyeggqaduhgxhyajtdt2j7wad.onion) | No | DE | No | https://artemislena.eu |

[Status Page](https://github.com/rramiachraf/dumb-instances)

#### Notes:
- Instances list in JSON format can be found in [instances.json](instances.json) file.
- For people who might be capable and interested in hosting a public instance feel free to do so, and don't forget to open a pull request, so your instance can be included here.

## Contributing
Contributions are welcome.

## License
[MIT](https://github.com/rramiachraf/dumb/blob/main/LICENCE)

