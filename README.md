# Randomware

Requirements:
    Walk (GUI)
    Setup:
        -install https://lh5.googleusercontent.com/ddRBeYejX9Qp5Vdf5tNI319D5XUkt6BHXrkSFGbxVzjgYhBM8uYOyB87KVhh08j-QKJJKgIxUDYrIqilIBpCLsjcTpg9BUiqT6WogEtSSoEg9TCj2f8wEY-SLlJq7g=w740
        -check that C:\TDM-GCC-64\bin is in PATH
        -execute : go get github.com/lxn/walk
        -execute: go get github.com/akavel/rsrc
    Compile:
        go build -ldflags="-H windowsgui"