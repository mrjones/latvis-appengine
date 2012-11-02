if [[ $1 == "run" ]]; then
    dev_appserver.py --address linode.mrjon.es .
fi

if [[ $1 == "upload" ]]; then
    appcfg.py update .
fi

