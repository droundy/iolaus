# How to install iolaus

Get a copy of iolaus with

    git clone github.com/droundy/iolaus.git

Then you can build and install it with

    cd iolaus
    make install

This last step will download and install my `goopt` package (using
`goinstall`) if you haven't got it.  If you already have an
out-of-date copy of `goopt` installed, you may need to run

    goinstall -u github.com/droundy/goopt

before compiling iolaus.
