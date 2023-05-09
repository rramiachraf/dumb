# dumb
With the massive daily increase of useless scripts on Genius's web frontend and having to download megabytes of clutter, [dumb](https://github.com/rramiachraf/dumb) tries to make reading lyrics from Genius a pleasant experience and as lightweight as possible.

<a href="https://codeberg.org/rramiachraf/dumb"><img src="https://img.shields.io/badge/Codeberg-%232185d0" /></a>

![Screenshot](https://raw.githubusercontent.com/rramiachraf/dumb/main/screenshot.png)

## Installation & Usage
[Go 1.18+](https://go.dev/dl) is required.
```bash
git clone https://github.com/rramiachraf/dumb
cd dumb
go build
./dumb
```

The default port is 5555, you can use other ports by setting the `PORT` environment variable.

## Public Instances

| URL                                           | Region | CDN? | Operator         |
| ---                                           | ---    | ---  | ---              |
| <https://dm.vern.cc>                          | US     | No   | https://vern.cc  |
| <https://dumb.nunosempere.com> (experimental) | DE     | No   | @NunoSempere     |
| <https://sing.whatever.social>                | US/DE  | Yes  | Whatever Social  |
| <https://dumb.lunar.icu>                      | DE     | Yes  | @MaximilianGT500 |
| <https://dumb.esmailelbob.xyz>                | CA     | No   | https://esmailelbob.xyz |
| <https://dumb.privacydev.net>                 | DE     | No   | https://privacydev.net  |
| <https://sing.whateveritworks.org>                 | DE     | Yes  | https://whateveritworks.org  |

### Tor
| URL                                                                        | Operator        |
| ---                                                                        | ---             |
| <http://dm.vernccvbvyi5qhfzyqengccj7lkove6bjot2xhh5kajhwvidqafczrad.onion> | https://vern.cc |
| <http://dumb.esmail5pdn24shtvieloeedh7ehz3nrwcdivnfhfcedl7gf4kwddhkqd.onion> | https://esmailelbob.xyz |
| <http://dumb.g4c3eya4clenolymqbpgwz3q3tawoxw56yhzk4vugqrl6dtu3ejvhjid.onion> | https://privacydev.net  |

### I2P
| URL                                                                   | Operator        |
| ---                                                                   | ---             |
| <http://vernxpcpqi2y4uhu7to4rnjmyjjgzh3x3qxyzpmkhykefchkmleq.b32.i2p> | https://vern.cc |

For people who might be capable and interested in hosting a public instance feel free to do so and don't forget to open a pull request so your instance can be included here.

## Contributing
Contributions are welcome.

## License
[MIT](https://github.com/rramiachraf/dumb/blob/main/LICENCE)

