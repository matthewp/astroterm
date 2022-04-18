#/bin/bash

mkdir tmp

for filename in build/*; do
    echo "Packaging $filename"

    # Make executable
    chmod +x $filename
    echo "\tGave executable permissions"

    # Figure out our filenames
    ext=$(basename $filename | sed 's/.*\.//')

    bn=$(basename $filename)
    if [ $ext == "exe" ]; then
        bn=$(basename $filename | sed 's/\.[^.]*$//')
    fi
    echo "\tBasename is $bn"
    
    aext="tar.gz"
    if [ $ext == "exe" ]; then
        aext="zip"
    fi
    
    archivename="$bn.$aext"
    echo "\tArchive name $archivename"

    cp LICENSE README.md tmp

    if [ $ext == "exe" ]; then
        mv $filename "tmp/astroterm.exe"
        echo "\tMoved executable to tmp/astroterm.exe"
        cd tmp
        zip $archivename *
        echo "\tZipped to $archivename"
    else
        mv $filename "tmp/astroterm"
        echo "\tMoved executable to tmp/astroterm"
        cd tmp
        tar zcvf $archivename *
        echo "\tTarred and gzipped $archivename"
    fi
    mv $archivename ../build

    cd ..
    rm tmp/*
done

rmdir tmp