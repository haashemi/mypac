# 🗂 MyPac

MyPac is a really simple cli application that helps you to organize your explcitly installed pacman packages in a yaml file.

![Preview of MyPac](https://github.com/user-attachments/assets/2467d07d-acbf-4107-8818-263d47f49064)

## Why tho?

I usually install too many packages when doing experiments, and after some time, I'd forget why do I have such a package installed and wasn't sure if removing it is a good idea or not. So I though of so many complex solutions and ended up with the simplest one.

It reads the `pacman -Qqe` output and compares with the YAML file I have written for myself. It not only helps me to have my packages kinda organized, but also lets me know if there's a package missing for whatever reason. And as the YAML file is completely handwritten, I can write comments all over the file for extra context of each package or group when needed.

## Usage:

0. Download the binary or build it yourself.
1. Create a `package.yml` file
2. Run the application using `./mypac`

*Note:* You can optionally pass a `-path` flag to use different `pacakge.yml` file anywhere with different names.

## Package.yml structure:

This file has a really simple structure; It's just an array of package groups which can be nested togehter too.

Here is its type:
```go
type PackageGroup struct {
	Name     string          `yaml:"name"`
	Packages []string        `yaml:"packages"`
	Groups   []*PackageGroup `yaml:"groups"`
}
```

And here's an example:
```yaml
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
```
