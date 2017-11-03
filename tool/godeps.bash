set -e

pkg="$1"
main="$2"
suffix="$3"

local_dependencies() {
	go list -f '{{ join .Deps  "'${suffix}'\n"}}' $pkg | grep $main | grep -v vendor
}

local_dependencies
