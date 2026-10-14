# 🗂 MyPac

MyPac is a really simple cli application that helps you to organize your explcitly installed pacman packages in a yaml file.

![Preview of MyPac](https://github.com/user-attachments/assets/2467d07d-acbf-4107-8818-263d47f49064)

## Why tho?

I usually install too many packages when doing experiments, and after some time, I'd forget why do I have such a package installed and wasn't sure if removing it is a good idea or not. So I though of so many complex solutions and ended up with the simplest one.

It reads the `pacman -Qqe` output and compares with the YAML file I have written for myself. It not only helps me to have my packages kinda organized, but also lets me know if there's a package missing for whatever reason. And as the YAML file is completely handwritten, I can write comments all over the file for extra context of each package or group when needed.

The same goes for files: it lists the files no package owns (using [lostfiles](https://github.com/graysky2/lostfiles)) next to the reason you wrote for each, and the package files you've modified (mostly in `/etc`), so you know what to back up.

## Install

### With makepkg (recommended)

The repo ships a `PKGBUILD` for `mypac-git`, which builds the latest commit on `main`:

```sh
git clone https://github.com/haashemi/mypac.git
cd mypac
makepkg -si
```

`makepkg` clones the sources from GitHub itself; your checkout only provides the `PKGBUILD`. To update, run `git pull && makepkg -si` again. The version comes from `git describe`, so pacman treats each new commit as an upgrade. Remove it with `sudo pacman -R mypac-git`.

### With Go

```sh
go install github.com/haashemi/mypac@latest
```

This way you also need the `lostfiles` package (`sudo pacman -S lostfiles`); the PKGBUILD pulls it in for you.

## Usage

mypac must be run as root (`lostfiles` and reading every `/etc` file need it):

```
sudo mypac [command] [flags]

  packages    Categorize explicitly installed packages (default)
  lostfiles   List files not owned by any package, grouped by reason
  modified    List package files modified since installation (mostly /etc)
  help        Show this help
```

Every command accepts `-config <path>`. `lostfiles` and `modified` accept `-paths` to print only the paths, handy for backups:

```sh
{ sudo mypac modified -paths; sudo mypac lostfiles -paths; } | sudo tar -czf backup.tar.gz -T -
```

`lostfiles -strict` runs lostfiles in its strict mode.

## Config

Unless `-config` is given, the first file found is used:

1. `/etc/mypac.{yml,yaml}`
2. `/etc/mypac/mypac.{yml,yaml}`

```yaml
# Package groups, which can be nested.
packages:
  - name: Arch Linux
    packages:
      - base
      - linux
    groups:
      - name: Sound
        packages:
          - pipewire

  - name: Development
    packages:
      - docker
      - git
      - go

# Why a file not owned by any package exists. A directory covers everything
# under it, and the most specific path wins.
lostfiles:
  /etc/asusd: asusctl's runtime config
  /opt/some-app: installed manually from upstream tarball

# Modified package files to leave out of `mypac modified`.
modified:
  ignore_packages:
    - filesystem
  ignore_files:
    - /etc/locale.gen
```

The old format, where the whole file is just the list of package groups, still works.
