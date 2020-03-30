# **Init:**

`stevedore init -f example.yaml`

###### **What it does:**

Stevedore eventually calls helm. Hence the tiller in the appropriate namespace has to be installed.
The init command initialises the tiller in the namespace(s) given in the manifest(s).

|  Flags 	| Description  	| Notes  	|
|---	|---	|---	|
| -f  	|  manifest file or directory with manifest files 	| Atleast one should be given |
| -n  	| namespace(s) in which to do the init  	| Optional |
| --kube-config  	| path to kube config file  	| Optional. Defaults to `~/.kube/config` |
| --recursive  	| respects `depends_on` and does init based on dependent manifests also  	| Optional |
| --force-upgrade  	| forces upgrade/downgrade of tiller to the helm version supported by stevedore	| Optional. May be compulsory depending on the already installed tiller version installed in the namespace |

This commands precedes all other commands.
 
Todo:

Gif