:
path="`pwd`/mypkg"
cnt=`echo $GOPATH | grep -c $path`
if [ $cnt -eq 0 ] ; then
	export GOPATH="$GOPATH:`pwd`/mypkg:$GOROOT/pkg"
fi
echo "(new) GOPATH = $GOPATH"
