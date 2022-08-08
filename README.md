# mac-sync

Sync specified programs and config files between macs using Git.

If you want to backup your whole mac os system, use [Time Machine](https://support.apple.com/en-us/HT201250) instead.

## Prerequisite

- Git account setting - This program assumes there is your email information in your `~/.gitconfig`.

## How to set up

1. Create private repository named `mac-sync-configs` in Github.

2. Add `programs.yaml`, `configs.yaml` to the `main` branch of the repository.

3. Run `mac-sync upload-configs` to upload configuration files to the repository.

4. In another mac, run `mac-sync download-configs` to download configuration files from the repository.

## Configuration files

### `programs.yaml`

Example:

```
homebrew:
  install: brew install {program}
  uninstall: brew uninstall {program}
  programs:
    - python3

npm:
  install: npm i -g {program}
  uninstall: npm rm -g {program}
  programs:
    - ts-node
```

### `configs.yaml`

Example:

```
sync:
  # macos configs
  - ~/Library/Preferences/com.apple.dock.plist
```

## Usage

```
NAME:
   mac-sync - Sync specified programs and config files between macs using Git.

USAGE:
   mac-sync [global options] command [command options] [arguments...]

COMMANDS:
   upload-configs, u    Upload local config files to remote
   download-configs, d  Download configs from remote
   sync-programs, s     Sync programs with remote
   clear-cache, c       Clear cache
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

## Example

- [mac-sync-configs](https://github.com/jopemachine/mac-sync-configs)
