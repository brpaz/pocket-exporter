# Pocket Exporter

[![Build status](https://img.shields.io/travis/com/brpaz/pocket-exporter.svg)](https://travis-ci.org/brpaz/pocket-exporter)
[![License](https://img.shields.io/github/license/brpaz/pocket-exporter.svg)](https://github.com/brpaz/pocket-exporter/LICENSE)
[![Latest Release](https://img.shields.io/github/release/brpaz/pocket-exporter.svg)](https://github.com/brpaz/pocket-exporter/releases/latest)

> Command line tool that allows to export your [Pocket](https://getpocket.com) articles in a json file.

## Pre-Requisites

- You must have a [Pocket](https://getpocket.com/register) account and a valid access token to authenticate on the Pocket API. You can get one by creating an application [here](https://getpocket.com/developer/). Make sure you grant "retrieve" permissions and select "Desktop (other)" as platform.

## Install

To install, the best way is to use the compiled binary.

You can get the latest one from the [Releases](https://github.com/brpaz/pocket-exporter/releases) page.

Then extract its contents and move to the appropriate folder:

```
sudo mv pocket-exporter /usr/local/bin/pocket-exporter
sudo chmod +x /usr/local/bin/pocket-exporter
```

## Usage

Open a terminal window, run:

```
pocket-exporter --consumerKey="some-consumer-key" -o /tmp/teste.json-o="/path/to/some/file.json"
```

## Contributing

All contributions are welcome.

Please see [Contributing.md](CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for details.

## Authors

- [Bruno Paz](https://github.com/user/brpaz)

## License

This project is Licensed under [MIT](LICENSE) License
