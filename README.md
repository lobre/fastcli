# fastcli

Wrapper around httpie to easily call Fastly API

Can be called with something like:

    fastcli -e myenv version/36/response_object

Taking to accound that `$HOME/.fastcli` exists with:

    [
      { "name": "myenv", "id": "myfastlyid", "token": "myfastlytoken }
    ]
