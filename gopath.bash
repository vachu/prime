:
path="`pwd`/mypkg"
cnt=`echo $GOPATH | grep -c $path`
if [ $cnt -eq 0 ] ; then
	export GOPATH="`pwd`/mypkg"
else
	export GOPATH="$GOPATH:`pwd`/mypkg"
fi
echo "(new) GOPATH = $GOPATH"
