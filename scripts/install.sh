#!/bin/sh

freebsdH="dae9b5b61e96919e94f7cd9510fd06ce318290fdee9dcd459005baaf67471377"
openbsdH="49a8a899302267a103a7406032debef40d7cfd25cf856cf3a34cac969fba3490"
linuxH="711265cc460989a42906b0ff46946dd55810682587469623ce86776223f02347"

earlyCheck(){
    os=`uname`
    os=`echo $os | tr "[:upper:]" "[:lower:"]`

    case $os in
        *openbsd* | *linux* | *freebsd* ) ;;
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
    url="https://archive.org/download/cetus-v0.5.1/cetus-v0.5.1-$os-$cpu"
}

printURL(){
    echo "You can get the Pre-built binary here:"
    echo "$url"
    echo
    echo "Run these commands to install it on your device."
    echo "# curl -L -o /usr/local/bin/cetus $url"
    echo "# chmod +x /usr/local/bin/cetus"
    echo
    echo "You may want to verify the hash of the downloaded file."
    echo "This is sha256 hash for cetus built for: $os $cpu"
    case $os in
        *openbsd* )
            echo "$openbsdH"
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

echo "Cetus v0.5.1"
echo
earlyCheck
getURL
printURL
