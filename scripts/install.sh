#!/bin/sh

freebsdH="b14b8e41892ec7b7f6a8faa0c868927f68e300b3cda36e1af79683c3313daed7"
openbsdH="5c6105233b886327699a571e1ed48c434124e94dfa4e08e1b8348a46da7ad1b2"
linuxH="79b15576fb37ef78ffb424e414eec483c40d237ea40e8234235fa5b1a322a3b2"
netbsdH="b66ce80fca51cce56b2c48723a5bfb137ef378836062f5054b5e447fc183967f"
dragonflyH="8eeebfebc71cba05eff71f3b38ab5b60e1c1334730083b3705ba463845034fb8"
darwinH="bd46b41a86a19c74d791240ae241a0784b7fbd3920c85ac12d1f51f0d741981f"

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
    url="https://archive.org/download/cetus-v0.6.0/cetus-v0.6.0-$os-$cpu"
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

echo "Cetus v0.6.0"
echo
earlyCheck
getURL
printURL
