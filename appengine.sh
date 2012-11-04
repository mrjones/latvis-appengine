if [[ $1 == "run" ]]; then
    dev_appserver.py --address linode.mrjon.es .
fi

if [[ $1 == "upload" ]]; then
    appcfg.py --oauth2 --noauth_local_webserver update .
fi

