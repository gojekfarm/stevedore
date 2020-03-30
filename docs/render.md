# **Render:**

###### **With Specific Files:**

`stevedore render -f example.yaml -o example_override_example.yaml -e dev_environment.yaml`

###### **With directories:**
`stevedore render -f services/ -o overrides/ -e envs/`

###### **What it does:**

A stevedore manifest may contain variables that will be filled during plan and apply. This command is used to render the manifest with all these variables.
It gives the user the chance to have a look at the various variables that will be used and sent to the helm charts eventually. 

These variables may come from different sources. It can be from an override or pulled from a different store. 

**Flags accepted:**

|  Flags 	| Description  	| Notes  	|
|---	|---	|---	|
| -f  	|  manifest file or directory with manifest files 	| Atleast one should be given |
| -o  	| override file or directory with override files  	| Optional |
| -e  	| env file or directory with env files  	| Optional |
| -a  	| output of artifacts path 	| Optional |
| -r  	| helm repo name to which the charts will be pushed/pulled  	| Optional. Defaults to `chartmuseum` |
| -t  	| helm timeout  	| Optional. Defaults to `600 seconds` |
| --recursive  	| respects `depends_on` in manifest files and recursively renders/plans/applies for those manifests too  	| Optional|
| --config  	|  config file that stevedore looks at for environment information 	| Optional. Defaults to `/usr/local/etc/stevedore-config.yaml` |
| --log-level  	|  sets the logger level 	| Optional. Defaults to `error` |

_**Atleast one manifest file should be given**_. Overrides/Envs file(s) are optional. If provided they will be used.

Stevedore goes through all the manifests, collates all matching overrides and envs for the appropriate manifests.
It shows the various kubernetes entities created/modified/deleted. 

A very useful detail to be noted is the `depends_on` portion. Using this, one manifest can depend on another manifest.
If the recursive flag is provided, now stevedore will recursively also do a render/plan/apply for the manifests dependencies mentioned in this part of the manifest. 

The render command will show the manifest with the final values that will be sent to helm and also displays information to the user about the source of the variable values.

ToDo: 

Gif to be added.


