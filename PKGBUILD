# Maintainer: Ali Hashemi <pvwvwx@gmail.com>
pkgname=mypac-git
_pkgname=mypac
pkgver=1.0.0.r2.g64babbb
pkgrel=1
pkgdesc='Organize your explicitly installed pacman packages, lost files and modified files with a YAML file'
arch=('x86_64' 'aarch64')
url='https://github.com/haashemi/mypac'
license=('MIT')
depends=('glibc' 'pacman' 'lostfiles')
makedepends=('git' 'go')
provides=("$_pkgname")
conflicts=("$_pkgname")
source=("git+$url.git")
sha256sums=('SKIP')

pkgver() {
  cd "$_pkgname"
  git describe --long --tags --abbrev=7 | sed 's/^v//;s/\([^-]*-g\)/r\1/;s/-/./g'
}

prepare() {
  cd "$_pkgname"
  go mod download
}

build() {
  cd "$_pkgname"
  export CGO_CPPFLAGS="$CPPFLAGS"
  export CGO_CFLAGS="$CFLAGS"
  export CGO_CXXFLAGS="$CXXFLAGS"
  export CGO_LDFLAGS="$LDFLAGS"
  export GOFLAGS='-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw'
  go build -o "$_pkgname" .
}

check() {
  cd "$_pkgname"
  go test ./...
}

package() {
  cd "$_pkgname"
  install -Dm755 "$_pkgname" -t "$pkgdir/usr/bin"
  install -Dm644 README.md -t "$pkgdir/usr/share/doc/$_pkgname"
  install -Dm644 LICENSE -t "$pkgdir/usr/share/licenses/$_pkgname"
}
