#!/bin/sh

freebsdH="47913c8eebb2481bf1282535231cc82793c345ad4c0f382af50e1d7760f5ad32"
openbsdH="2abdad0520faff64cd3bdf92e194eaf9bafe1866dcde0b432f2a7e288a0806df"
linuxH="cba19402714d27ce64e57edecab9acc502542c1abb828ad9c5bd7db4241b0ff0"
netbsdH="d1346f5478931f21d6a9764287230d9362f312c91836773f6a795c7dcc77791c"
dragonflyH="88ba5cd45abf7a16e6ac3aee8d6e2403e700ed170c6f2191f4474e218cc31cb6"
darwinH="aea04ff4b027ca7d32d4cf56bdabd2cb3bc0f466ecc634cad0590b9ae774feb6"

earlyCheck(){
    os=`uname`
    os=`echo $os | tr "[:upper:]" "[:lower:"]`

    case $os in
	# Not sure about uname output on DragonFly BSD.
        *openbsd* | *linux* | *freebsd* | *netbsd* | *dragonfly* | *dragonflybsd* | *darwin* ) ;;
        *)
            echo "Pre-built binary not available for your os"
            exit 1
            ;;
    esac

    cpu=`uname -m`
    cpu=`echo $cpu | tr "[:upper:]" "[:lower:"]`

    case $cpu in
        *amd*64* | *x86*64* ) ;;
        *)
            echo "Pre-built binary not available for your cpu"
            exit 1
            ;;
    esac
}

getURL(){
    url="https://archive.org/download/cetus-v0.6.4/cetus-v0.6.4-$os-$cpu"
}

printURL(){
    echo "You can get the Pre-built binary here:"
    echo "$url"
    echo
    echo "Run these commands to install it on your device."
    echo "# curl -L -o /usr/local/bin/cetus $url"
    echo "# chmod +x /usr/local/bin/cetus"
    echo
    echo "This is sha256 hash for cetus built for: $os $cpu"
    case $os in
        *openbsd* )
            echo "$openbsdH"
            ;;
	*netbsd* )
            echo "$netbsdH"
            ;;
	*dragonflybsd* | *dragonfly* )
            echo "$dragonflyH"
            ;;
	*darwin* )
            echo "$darwinH"
            ;;
        *freebsd* )
            echo "$freebsdH"
            ;;
        *linux* )
            echo "$linuxH"
            ;;
    esac
    echo
    echo "Verify the hash by running sha256 on cetus binary."
    echo "$ sha256 /usr/local/bin/cetus"
}

echo "Cetus v0.6.4"
echo
earlyCheck
getURL
printURL
